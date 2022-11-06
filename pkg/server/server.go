package server

import (
	"net/http"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	// For Swagger docs
	_ "github.com/adh-partnership/api/docs"
	"github.com/adh-partnership/api/internal/v1/router"
	v1storage "github.com/adh-partnership/api/internal/v1/storage"
	"github.com/adh-partnership/api/pkg/config"
	"github.com/adh-partnership/api/pkg/database"
	"github.com/adh-partnership/api/pkg/database/models"
	"github.com/adh-partnership/api/pkg/discord"
	"github.com/adh-partnership/api/pkg/gin/middleware/auth"
	ginLogger "github.com/adh-partnership/api/pkg/gin/middleware/logger"
	"github.com/adh-partnership/api/pkg/gin/response"
	"github.com/adh-partnership/api/pkg/logger"
	"github.com/adh-partnership/api/pkg/oauth"
	"github.com/adh-partnership/api/pkg/storage"
)

type ServerStruct struct {
	Engine          *gin.Engine
	Config          *config.Config
	TrackedPrefixes map[string]bool
}

var Server *ServerStruct

type ServerOpts struct {
	ConfigFile string
}

var log = logger.Logger.WithField("component", "server")

func NewServer(o *ServerOpts) (*ServerStruct, error) {
	s := &ServerStruct{}

	log.Infof("Loading config file: %s", o.ConfigFile)
	cfg, err := config.ParseConfig(o.ConfigFile)
	if err != nil {
		return nil, err
	}
	s.Config = cfg
	config.Cfg = cfg

	s.TrackedPrefixes = make(map[string]bool)

	for _, prefix := range cfg.Facility.Stats.Prefixes {
		s.TrackedPrefixes[prefix] = true
	}

	log.Info("Connecting to database")
	err = database.Connect(database.DBOptions{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		Database: cfg.Database.Database,
		Driver:   "mysql",
		Logger:   logger.Logger,
	})
	if err != nil {
		return nil, err
	}

	log.Info("Running migrations...")
	err = database.DB.AutoMigrate(
		&models.APIKeys{},
		&models.ControllerStat{},
		&models.DelayedJob{},
		&models.Document{},
		&models.EmailTemplate{},
		&models.EventPosition{},
		&models.Event{},
		&models.EventSignup{},
		&models.Feedback{},
		&models.Flights{},
		&models.OAuthClient{},
		&models.OAuthLogin{},
		&models.OAuthRefresh{},
		&models.OnlineController{},
		&models.Rating{},
		&models.Role{},
		&models.TrainingNote{},
		&models.User{},
		&models.VisitorApplication{},
	)
	if err != nil {
		log.Errorf("Failed to run migrations: %v", err)
		return nil, err
	}

	log.Info("Configuring Discord package")
	discord.SetupWebhooks(cfg.Discord.Webhooks)

	log.Info("Building OAuth2 Clients")
	oauth.BuildWithConfig(cfg)

	log.Info("Building storage objects")
	log.Info(" - Uploads")
	log.Debugf("Config: %+v", cfg)
	_, err = storage.Configure(cfg.Storage, "uploads")
	if err != nil {
		return nil, err
	}
	if cfg.Storage.BaseURL != "" {
		log.Infof(" - Setting BaseURL to %s", cfg.Storage.BaseURL)
		v1storage.SetBase(cfg.Storage.BaseURL)
	}

	log.Info("Building gin engine")
	gin.SetMode(gin.ReleaseMode)
	s.Engine = gin.New()
	s.Engine.Use(gin.Recovery())
	s.Engine.Use(ginLogger.Logger)
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowMethods = []string{"GET", "PATCH", "POST", "PUT", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Authorization", "X-Requested-With", "Accept"}
	corsConfig.AllowCredentials = true
	corsConfig.AllowWildcard = true
	// Use this instead of AllowAllOrigins so that we return the origin and not '*'
	corsConfig.AllowOriginFunc = func(origin string) bool {
		return true
	}
	s.Engine.Use(cors.New(corsConfig))

	cookieOpts := sessions.Options{
		Domain:   cfg.Session.Cookie.Domain,
		Path:     cfg.Session.Cookie.Path,
		MaxAge:   cfg.Session.Cookie.MaxAge,
		HttpOnly: true,
		Secure:   cfg.Session.Cookie.Secure,
	}
	switch strings.ToLower(cfg.Session.Cookie.SameSite) {
	case "none": // Useful for local development against the staging API
		cookieOpts.SameSite = http.SameSiteNoneMode
		cookieOpts.Secure = true
	case "lax":
		cookieOpts.SameSite = http.SameSiteLaxMode
	case "strict":
		cookieOpts.SameSite = http.SameSiteStrictMode
	default:
		cookieOpts.SameSite = http.SameSiteDefaultMode
	}

	store := cookie.NewStore([]byte(cfg.Session.Cookie.Secret))
	store.Options(cookieOpts)
	s.Engine.Use(sessions.Sessions(cfg.Session.Cookie.Name, store))
	s.Engine.Use(auth.UpdateCookie)
	s.Engine.Use(auth.Auth)

	log.Info("Registering static routes and templates")
	s.Engine.LoadHTMLGlob("static/*.html")
	s.Engine.Static("/static", "static")
	s.Engine.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{})
	})
	s.Engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	s.Engine.GET("/ping", func(c *gin.Context) {
		response.RespondMessage(c, http.StatusOK, "PONG")
	})

	log.Info("Registering routes")
	router.SetupRoutes(s.Engine)
	s.Engine.HandleMethodNotAllowed = true

	Server = s

	return s, nil
}

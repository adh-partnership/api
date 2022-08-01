package server

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	// For Swagger docs
	_ "github.com/kzdv/api/docs"
	"github.com/kzdv/api/internal/v1/router"
	"github.com/kzdv/api/pkg/config"
	"github.com/kzdv/api/pkg/database"
	"github.com/kzdv/api/pkg/discord"
	"github.com/kzdv/api/pkg/gin/middleware/auth"
	ginLogger "github.com/kzdv/api/pkg/gin/middleware/logger"
	"github.com/kzdv/api/pkg/gin/response"
	"github.com/kzdv/api/pkg/logger"
	"github.com/kzdv/api/pkg/oauth"
)

type ServerStruct struct {
	Engine *gin.Engine
	Config *config.Config
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

	log.Info("Configuring Discord package")
	discord.SetupWebhooks(cfg.Discord.Webhooks)

	/*
		Commented out while we plan for Redis and possible RabbitMQ

			log.Info("Connecting to redis")
			database.ConnectRedis(database.RedisOptions{
				Password:      cfg.Redis.Password,
				DB:            cfg.Redis.Database,
				Sentinel:      cfg.Redis.Sentinel,
				MasterName:    cfg.Redis.MasterName,
				SentinelAddrs: cfg.Redis.SentinelAddrs,
				Addr:          cfg.Redis.Address,
			})
	*/

	log.Info("Building OAuth2 Client")
	oauth.Build(
		cfg.OAuth.ClientID,
		cfg.OAuth.ClientSecret,
		fmt.Sprintf("%s%s", cfg.OAuth.MyBaseURL, "/v1/login/callback"),
		fmt.Sprintf("%s%s", cfg.OAuth.BaseURL, cfg.OAuth.Endpoints.Authorize),
		fmt.Sprintf("%s%s", cfg.OAuth.BaseURL, cfg.OAuth.Endpoints.Token),
	)

	log.Info("Building gin engine")
	gin.SetMode(gin.ReleaseMode)
	s.Engine = gin.New()
	s.Engine.Use(gin.Recovery())
	s.Engine.Use(ginLogger.Logger)
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowMethods = []string{"GET", "PATCH", "POST", "PUT", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Authorization", "X-Requested-With", "Accept"}
	corsConfig.AllowCredentials = true
	corsConfig.AllowWildcard = true
	s.Engine.Use(cors.New(corsConfig))

	store := cookie.NewStore([]byte(cfg.Session.Cookie.Secret))
	store.Options(sessions.Options{
		Domain:   cfg.Session.Cookie.Domain,
		Path:     cfg.Session.Cookie.Path,
		MaxAge:   cfg.Session.Cookie.MaxAge,
		HttpOnly: true,
	})
	s.Engine.Use(sessions.Sessions(cfg.Session.Cookie.Name, store))
	s.Engine.Use(auth.UpdateCookie)

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

	Server = s

	return s, nil
}

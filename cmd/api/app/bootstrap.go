package app

import (
	"fmt"
	"time"

	"github.com/urfave/cli/v2"

	"github.com/adh-partnership/api/pkg/config"
	"github.com/adh-partnership/api/pkg/database"
	"github.com/adh-partnership/api/pkg/database/models"
	"github.com/adh-partnership/api/pkg/jobs/roster"
	"github.com/adh-partnership/api/pkg/logger"
	"github.com/adh-partnership/api/pkg/network/vatusa"
)

func newBootstrapCommand() *cli.Command {
	return &cli.Command{
		Name:  "bootstrap",
		Usage: "Bootstrap the application",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "config",
				Value: "config.yaml",
				Usage: "Path to the configuration file",
			},
			&cli.BoolFlag{
				Name: "skip-seed",
			},
			&cli.BoolFlag{
				Name: "skip-roster",
			},
			&cli.BoolFlag{
				Name: "skip-migration",
			},
		},
		Action: func(c *cli.Context) error {
			log := logger.Logger.WithField("component", "bootstrap")
			configfile := c.String("config")
			log.Infof("Loading config file: %s", configfile)
			cfg, err := config.ParseConfig(configfile)
			if err != nil {
				return err
			}
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
				return err
			}

			if !c.Bool("skip-migration") {
				log.Info("Running database migrations")
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
					&models.Rating{},
					&models.Role{},
					&models.TrainingNote{},
					&models.User{},
					&models.VisitorApplication{},
				)
				if err != nil {
					return err
				}
			}

			if !c.Bool("skip-seed") {
				log.Info("Seeding database")
				seeds := []interface{}{
					&models.Rating{
						ID:    -1,
						Long:  "Inactive",
						Short: "INA",
					},
					&models.Rating{
						ID:    1,
						Long:  "Observer",
						Short: "OBS",
					},
					&models.Rating{
						ID:    2,
						Long:  "Student 1",
						Short: "S1",
					},
					&models.Rating{
						ID:    3,
						Long:  "Student 2",
						Short: "S2",
					},
					&models.Rating{
						ID:    4,
						Long:  "Student 3",
						Short: "S3",
					},
					&models.Rating{
						ID:    5,
						Long:  "Controller",
						Short: "C1",
					},
					&models.Rating{
						ID:    6,
						Long:  "Controller 2",
						Short: "C2",
					},
					&models.Rating{
						ID:    7,
						Long:  "Senior Controller",
						Short: "C3",
					},
					&models.Rating{
						ID:    8,
						Long:  "Instructor",
						Short: "I1",
					},
					&models.Rating{
						ID:    9,
						Long:  "Instructor 2",
						Short: "I2",
					},
					&models.Rating{
						ID:    10,
						Long:  "Senior Instructor",
						Short: "I3",
					},
					&models.Rating{
						ID:    11,
						Long:  "Supervisor",
						Short: "SUP",
					},
					&models.Rating{
						ID:    12,
						Long:  "Administrator",
						Short: "ADM",
					},
					&models.Role{
						Name: "atm",
					},
					&models.Role{
						Name: "datm",
					},
					&models.Role{
						Name: "ta",
					},
					&models.Role{
						Name: "ec",
					},
					&models.Role{
						Name: "fe",
					},
					&models.Role{
						Name: "wm",
					},
					&models.Role{
						Name: "events",
					},
					&models.Role{
						Name: "mtr",
					},
					&models.EmailTemplate{
						Name:    "activity_warning",
						Subject: "Inactivity Warning",
						Body: `<p>Hello {{.FirstName}} {{.LastName}},</p>

<p>This is a warning that you have not met the activity requirements as set forth under the facility policy as of today. 
If you do not meet the requirements by the 1st of the following month, you may be removed from the facility due to inactivity.</p>

<p>Obviously, we understand this is a hobby and that sometimes real life gets in the way. If this is the case, please reach out 
to the senior staff and let us know.</p>

<p>Thank you for your time and we hope to see you on the network soon!</p>

<p>Best Regards,<br>
{{range .findRole "atm"}}
{{.}}<br>
{{end}}
{{range .findRole "datm"}}
{{.}}<br>
{{end}}</p>`,
						EditGroup: "admin",
						CC:        "atm@denartcc.org, datm@denartcc.org",
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					},
					&models.EmailTemplate{
						Name:    "visitor_accepted",
						Subject: "Visitor Application Accepted",
						Body: `<p>Hello {{.FirstName}} {{.LastName}},</p>

<p>Your visitor application has been accepted. Shortly, the staff will be adding you to the roster.</p>

<p>Please ensure you read and adhere to our facility SOPs and join our Discord server to stay up to date with the latest 
information. The invite for this can be found on our website.</p>

<p>Thank you for your time and we hope to see you on the network soon!</p>

<p>Best Regards,<br>
{{range .findRole "atm"}}
{{.}}<br>
{{end}}
{{range .findRole "datm"}}
{{.}}<br>
{{end}}</p>`,
						EditGroup: "admin",
						CC:        "atm@denartcc.org, datm@denartcc.org",
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					},
					&models.EmailTemplate{
						Name:    "visitor_denied",
						Subject: "Visitor Application Denied",
						Body: `<p>Hello {{.FirstName}} {{.LastName}},</p>

<p>We regret to inform you that your visiting application has been denied. This usually means that you did not meet the requirements to be a visitor.</p>

<p>If you have any questions about the reason for this denial, please do not hesitate to contact the senior staff.</p>

<p>Best Regards,<br>
{{range .findRole "atm"}}
{{.}}<br>
{{end}}
{{range .findRole "datm"}}
{{.}}<br>
{{end}}</p>`,
						EditGroup: "admin",
						CC:        "atm@denartcc.org, datm@denartcc.org",
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					},
				}

				for _, seed := range seeds {
					log.Infof("Seed: %+v", seed)
					// This will warn if a record exists, in case the API has already been deployed
					// prior to bootstrapping... some records (ie Ratings) will already exist.
					if err := database.DB.FirstOrCreate(seed, seed).Error; err != nil {
						log.Warnf("Error seeding %+v: %s", seed, err)
					}
				}
			}

			if !c.Bool("skip-roster") {
				log.Info("Populating users table")
				err = roster.UpdateRoster()
				if err != nil {
					return err
				}
				roster.UpdateForeignRoster()
			}

			log.Info("Populating staff roles")
			fac, err := vatusa.GetFacility(cfg.VATUSA.Facility)
			if err != nil {
				return err
			}

			atm, err := database.FindUserByCID(fmt.Sprint(fac.Info.ATM))
			if err != nil {
				return err
			}
			if atm != nil {
				err := database.AddRoleStringToUser(atm, "atm")
				if err != nil {
					return err
				}
			}

			datm, err := database.FindUserByCID(fmt.Sprint(fac.Info.DATM))
			if err != nil {
				return err
			}
			if datm != nil {
				err := database.AddRoleStringToUser(datm, "datm")
				if err != nil {
					return err
				}
			}

			ta, err := database.FindUserByCID(fmt.Sprint(fac.Info.TA))
			if err != nil {
				return err
			}
			if ta != nil {
				err := database.AddRoleStringToUser(ta, "ta")
				if err != nil {
					return err
				}
			}

			ec, err := database.FindUserByCID(fmt.Sprint(fac.Info.EC))
			if err != nil {
				return err
			}
			if ec != nil {
				err := database.AddRoleStringToUser(ec, "ec")
				if err != nil {
					return err
				}
			}

			fe, err := database.FindUserByCID(fmt.Sprint(fac.Info.FE))
			if err != nil {
				return err
			}
			if fe != nil {
				err := database.AddRoleStringToUser(fe, "fe")
				if err != nil {
					return err
				}
			}

			wm, err := database.FindUserByCID(fmt.Sprint(fac.Info.WM))
			if err != nil {
				return err
			}
			if wm != nil {
				err := database.AddRoleStringToUser(wm, "wm")
				if err != nil {
					return err
				}
			}

			return nil
		},
	}
}

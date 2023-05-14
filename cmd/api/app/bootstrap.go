package app

import (
	"fmt"

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

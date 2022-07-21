package app

import (
	"fmt"

	"github.com/kzdv/api/pkg/config"
	"github.com/kzdv/api/pkg/database"
	"github.com/kzdv/api/pkg/facility"
	"github.com/kzdv/api/pkg/logger"
	"github.com/kzdv/api/pkg/network/global"
	"github.com/kzdv/api/pkg/network/vatusa"
	dbTypes "github.com/kzdv/types/database"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/urfave/cli/v2"
)

func newJobCommand() *cli.Command {
	return &cli.Command{
		Name:  "sync",
		Usage: "Sync Job",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "config",
				Value: "config.yaml",
				Usage: "Path to the configuration file",
			},
		},
		Action: func(c *cli.Context) error {
			log := logger.Logger.WithField("component", "job")
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

			controllers, err := vatusa.GetFacilityRoster("both")
			if err != nil {
				return err
			}

			updateid, _ := gonanoid.New(24)
			err = facility.UpdateControllerRoster(controllers, updateid)
			if err != nil {
				return err
			}

			// Update foreign visitors
			var users []dbTypes.User
			if err := database.DB.
				Where(dbTypes.User{ControllerType: dbTypes.ControllerTypeOptions["visit"]}).
				Not(dbTypes.User{Region: "AMAS", Division: "USA"}).Find(&users).Error; err != nil {
				log.Errorf("Error getting foreign visitors: %s", err)
			}
			for _, user := range users {
				location, err := global.GetLocation(fmt.Sprint(user.CID))
				if err != nil {
					log.Errorf("Error getting location for user %d: %s", user.CID, err)
					continue
				}
				user.Region = location.Region
				user.Division = location.Division
				user.Subdivision = location.Subdivision
				if err := database.DB.Save(&user).Error; err != nil {
					log.Errorf("Error saving user %d: %s", user.CID, err)
				}
			}

			// Users not part of the VATUSA roster will be removed from our roster
			if err := database.DB.Model(&dbTypes.User{}).Updates(dbTypes.User{
				ControllerType: dbTypes.ControllerTypeOptions["none"],
				UpdateID:       updateid,
			}).Not(dbTypes.User{UpdateID: updateid}).Error; err != nil {
				return err
			}

			return nil
		},
	}
}

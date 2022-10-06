package app

import (
	"github.com/urfave/cli/v2"

	"github.com/adh-partnership/api/pkg/config"
	"github.com/adh-partnership/api/pkg/database"
	"github.com/adh-partnership/api/pkg/database/models"
	"github.com/adh-partnership/api/pkg/discord"
	"github.com/adh-partnership/api/pkg/jobs/roster"
	"github.com/adh-partnership/api/pkg/logger"
)

func newUpdateRosterCommand() *cli.Command {
	return &cli.Command{
		Name:  "update-roster",
		Usage: "Update roster",
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

			log.Infof("Configuring Discord webhooks")
			discord.SetupWebhooks(cfg.Discord.Webhooks)

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

			log.Info("Running database migrations")
			err = database.DB.AutoMigrate(
				&models.DelayedJob{},
				&models.EmailTemplate{},
				&models.User{},
			)
			if err != nil {
				return err
			}

			log.Info("Running update roster job")
			err = roster.UpdateRoster()
			if err != nil {
				log.Errorf("Error updating roster: %s", err)
				return err
			}

			return nil
		},
	}
}

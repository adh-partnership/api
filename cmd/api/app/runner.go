package app

import (
	"time"

	"github.com/go-co-op/gocron"
	"github.com/urfave/cli/v2"

	"github.com/kzdv/api/pkg/config"
	"github.com/kzdv/api/pkg/database"
	dbTypes "github.com/kzdv/api/pkg/database/types"
	"github.com/kzdv/api/pkg/jobs/delayedjobs"
	"github.com/kzdv/api/pkg/jobs/emails"
	"github.com/kzdv/api/pkg/jobs/roster"
	"github.com/kzdv/api/pkg/logger"
	"github.com/kzdv/api/pkg/messaging"
)

func newRunnerCommand() *cli.Command {
	return &cli.Command{
		Name:  "runner",
		Usage: "Job Runner",
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

			log.Info("Running database migrations")
			err = database.DB.AutoMigrate(
				&dbTypes.DelayedJob{},
				&dbTypes.EmailTemplate{},
				&dbTypes.User{},
			)
			if err != nil {
				return err
			}

			log.Info("Configuring messaging")
			messaging.Setup(
				cfg.RabbitMQ.Host,
				cfg.RabbitMQ.Port,
				cfg.RabbitMQ.User,
				cfg.RabbitMQ.Password,
			)

			log.Info("Building email consumer")
			err = messaging.BuildConsumer("emails", emails.Handler)
			if err != nil {
				return err
			}

			log.Info("Building scheduled jobs")
			s := gocron.NewScheduler(time.UTC)
			log.Info(" - Roster")
			err = roster.ScheduleJobs(s)
			if err != nil {
				return err
			}
			log.Info(" - Delayed Jobs")
			err = delayedjobs.ScheduleJobs(s)
			if err != nil {
				return err
			}

			log.Info("Starting scheduled jobs on main thread")
			s.StartBlocking()

			return nil
		},
	}
}

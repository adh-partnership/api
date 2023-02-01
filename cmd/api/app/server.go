package app

import (
	"context"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/urfave/cli/v2"

	"github.com/adh-partnership/api/pkg/jobs/activity"
	"github.com/adh-partnership/api/pkg/jobs/dataparser"
	"github.com/adh-partnership/api/pkg/jobs/roster"
	"github.com/adh-partnership/api/pkg/logger"
	"github.com/adh-partnership/api/pkg/metrics"
	"github.com/adh-partnership/api/pkg/server"
)

var log = logger.Logger.WithField("component", "server")

func newServerCommand() *cli.Command {
	return &cli.Command{
		Name:  "server",
		Usage: "Start the server",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "config",
				Value: "config.yaml",
				Usage: "Path to the configuration file",
			},
		},
		Action: func(c *cli.Context) error {
			log.Info("Starting ADH-PARTNERSHIP API")
			log.Debugf("config=%s", c.String("config"))
			srvr, err := server.NewServer(&server.ServerOpts{
				ConfigFile: c.String("config"),
			})
			if err != nil {
				return err
			}

			ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
			defer stop()

			log.Info("Building scheduled jobs")
			s := gocron.NewScheduler(time.UTC)
			log.Info(" - Activity")
			err = activity.ScheduleJobs(s)
			if err != nil {
				return err
			}
			log.Info(" - Roster")
			err = roster.ScheduleJobs(s)
			if err != nil {
				return err
			}
			log.Info(" - VATSIM Data Parser")
			err = dataparser.Initialize(s)
			if err != nil {
				return err
			}

			if srvr.Config.Metrics.Enabled {
				log.Infof("Building Metrics")
				m := metrics.GetMonitor()
				m.SetMetricPath(srvr.Config.Metrics.Path)
				m.SetMetricPort(srvr.Config.Metrics.Port)
				log.Info("Registering Metrics middleware")
				m.Use(srvr.Engine)
				log.Info("Starting Metrics server on :%s%s", srvr.Config.Metrics.Port, srvr.Config.Metrics.Path)
				m.Start()
			}

			log.Info("Starting scheduled jobs")
			s.StartAsync()

			log.Infof("Starting server on :%s", srvr.Config.Server.Port)
			srv := &http.Server{
				Addr:    ":" + srvr.Config.Server.Port,
				Handler: srvr.Engine,
			}

			// We will run the server in a separate goroutine so that we can
			// stop it gracefully when we receive a signal to do so.
			go func() {
				if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					log.Errorf("Error starting server: %s", err.Error())
				}
			}()

			<-ctx.Done() // Block main thread until we receive a signal to terminate or interrupt
			log.Infof("Shutting down server")
			stop()
			ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
			defer cancel()
			if err := srv.Shutdown(ctx); err != nil {
				log.Errorf("Error shutting down server: %s. Server will be forced shutdown.", err.Error())
			}
			log.Infof("Server shut down complete.")

			return nil
		},
	}
}

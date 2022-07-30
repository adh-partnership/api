package app

import (
	"context"
	"errors"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/urfave/cli/v2"

	"github.com/kzdv/api/pkg/logger"
	"github.com/kzdv/api/pkg/server"
)

func NewRootCommand() *cli.App {
	return &cli.App{
		Name:  "app",
		Usage: "KZDV Monolithic API",
		Commands: []*cli.Command{
			newServerCommand(),
			newJobCommand(),
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "log-level",
				Value:   "info",
				Aliases: []string{"l"},
				Usage:   "Log level (accepted values: trace, debug, info, warn, error, fatal, panic)",
			},
			&cli.StringFlag{
				Name:  "log-format",
				Value: "text",
				Usage: "Log format (accepted values: text, json)",
			},
		},
		Before: func(c *cli.Context) error {
			format := c.String("log-format")
			if !logger.IsValidFormat(format) {
				return errors.New("invalid log format")
			}
			logger.NewLogger(format)

			if !logger.IsValidLogLevel(c.String("log-level")) {
				return errors.New("invalid log level")
			}

			l, _ := logger.ParseLogLevel(c.String("log-level"))
			logger.Logger.SetLevel(l)

			logger.Logger.Info("Starting KZDV API")

			s, err := server.NewServer(&server.ServerOpts{
				ConfigFile: c.String("config"),
			})
			if err != nil {
				return err
			}

			log.Info("Building Notify Context")
			ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
			defer stop()

			srv := &http.Server{
				Addr:    s.Config.Server.Host + ":" + s.Config.Server.Port,
				Handler: s.Engine,
			}
			go func() {
				// We will utilize Kubernetes' readiness and liveness probes to determine if the server is ready to serve traffic.
				// kubelet will kill the pod if it takes too long to start.
				if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					logger.Logger.WithField("component", "server").Errorf("Error starting server: %+v", err)
				}
			}()

			<-ctx.Done()
			logger.Logger.Info("Shutting down server")
			stop()

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			if err := srv.Shutdown(ctx); err != nil {
				logger.Logger.WithField("component", "server").Errorf("Error shutting down server cleanly: %+v", err)
			}
			logger.Logger.Info("Server stopped")

			return nil
		},
	}
}

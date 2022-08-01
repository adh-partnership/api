package app

import (
	"context"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/urfave/cli/v2"

	"github.com/kzdv/api/pkg/logger"
	"github.com/kzdv/api/pkg/server"
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
			log.Info("Starting KZDV API")
			log.Debugf("config=%s", c.String("config"))
			srvr, err := server.NewServer(&server.ServerOpts{
				ConfigFile: c.String("config"),
			})
			if err != nil {
				return err
			}

			ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
			defer stop()

			log.Infof("Starting server on %s:%s", srvr.Config.Server.Host, srvr.Config.Server.Port)
			srv := &http.Server{
				Addr:    srvr.Config.Server.Host + ":" + srvr.Config.Server.Port,
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

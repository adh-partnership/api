package app

import (
	"errors"

	"github.com/urfave/cli/v2"

	"github.com/adh-partnership/api/pkg/logger"
)

func NewRootCommand() *cli.App {
	return &cli.App{
		Name:  "app",
		Usage: "ADH-PARTNERSHIP Monolithic API",
		Commands: []*cli.Command{
			newServerCommand(),
			newRunnerCommand(),
			newUpdateRosterCommand(),
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

			return nil
		},
	}
}

/*
 * Copyright ADH Partnership
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */

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

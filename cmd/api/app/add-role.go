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
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/adh-partnership/api/pkg/config"
	"github.com/adh-partnership/api/pkg/database"
	"github.com/adh-partnership/api/pkg/logger"
)

func newAddRoleCommand() *cli.Command {
	return &cli.Command{
		Name:  "add-role",
		Usage: "Add a role to a user",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "cid",
				Required: true,
				Usage:    "CID of the user",
			},
			&cli.StringFlag{
				Name:     "role",
				Required: true,
				Usage:    "Role to add",
			},
			&cli.StringFlag{
				Name:  "config",
				Usage: "Path to the config file",
				Value: "config.yaml",
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

			log.Info("Attempting to add role to user")
			user, err := database.FindUserByCID(c.String("cid"))
			if err != nil {
				return err
			}
			if user == nil {
				return fmt.Errorf("user not found")
			}

			log.Info("Adding role to user")
			err = database.AddRoleStringToUser(user, c.String("role"))
			if err != nil {
				return err
			}

			log.Info("Role added successfully")

			return nil
		},
	}
}

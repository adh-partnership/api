package main

import (
	"os"

	"github.com/adh-partnership/api/cmd/api/app"
	"github.com/adh-partnership/api/pkg/logger"
)

func main() {
	a := app.NewRootCommand()
	err := a.Run(os.Args)
	if err != nil {
		logger.Logger.Errorf("Error starting application: %v", err)
		os.Exit(1)
	}
}

// @title ADH API
// @version 1.0
// @description ADH API

// @contact.name Daniel Hawton
// @contact.email daniel@hawton.org

// @license.name Apache
// @license.URL https://github.com/adh-partnership/api/blob/main/LICENSE

// @host network.denartcc.org
// @BasePath /v1
// @query.collection.format multi

// @securityDefinitions.basic BasicAuth
// @description Session Cookie

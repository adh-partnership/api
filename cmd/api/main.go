package main

import (
	"os"

	"github.com/kzdv/api/cmd/api/app"
	"github.com/kzdv/api/pkg/logger"
)

func main() {
	a := app.NewRootCommand()
	err := a.Run(os.Args)
	if err != nil {
		logger.Logger.Errorf("Error starting application: %v", err)
		os.Exit(1)
	}
}

// @title KZDV API
// @version 1.0
// @description KZDV API

// @contact.name Daniel Hawton
// @contact.email wm@denartcc.org

// @license.name Apache
// @license.URL https://github.com/kzdv/api2/blob/main/LICENSE

// @host network.denartcc.org
// @BasePath /v1
// @query.collection.format multi

// @securityDefinitions.basic BasicAuth
// @description Session Cookie

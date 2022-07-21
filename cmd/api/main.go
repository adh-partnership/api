package main

import (
	"os"

	"github.com/kzdv/api/cmd/api/app"
)

func main() {
	a := app.NewRootCommand()
	_ = a.Run(os.Args)
}

// @title KZDV API
// @version 1.0
// @description KZDV API

// @contact.name Daniel Hawton
// @contact.email wm@denartcc.org

// @license.name Apache
// @license.url https://github.com/kzdv/api2/blob/main/LICENSE

// @host network.denartcc.org
// @BasePath /v1
// @query.collection.format multi

// @securityDefinitions.basic BasicAuth
// @description Session Cookie

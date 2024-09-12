package main

import (
	"github.com/fvbock/endless"
	"lmich.com/tononkira/config"
	"lmich.com/tononkira/routes"
)

func main() {
	setupOption := new(config.SetupOptions)
	setupOption.Run = func() {
		router := new(routes.Router)
		router.InitRoutes()
		// router.Run()
		endless.ListenAndServe(":"+config.Env.ApiPort, router.Router)

	}
	config.Setup(setupOption)
}

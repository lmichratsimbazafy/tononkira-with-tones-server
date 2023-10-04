package main

import (
	"lmich.com/tononkira/config"
	"lmich.com/tononkira/controllers"
)

func main() {
	setupOption := new(config.SetupOptions)
	setupOption.Run = func() {
		router := new(controllers.Router)
		router.Init()
		router.Run()
	}

	config.Setup(setupOption)
}

package config

import "fmt"

type SetupOptions struct {
	Run         func()
	OnClose     func()
	NeedLocalDb bool
}

func Setup(options *SetupOptions) {
	Getenv()
	dbConnexion := new(Db)
	dbConnexion.Client = dbConnexion.Connect(options.NeedLocalDb)
	DatabaseInstance = dbConnexion.GetDbInstance()
	fmt.Println("connected to mongoDb")
	options.Run()
	defer close(dbConnexion, options.OnClose)
}

func close(dbConnection *Db, onClosing func()) {
	if onClosing != nil {
		onClosing()
	}
	dbConnection.Disconnect()
}

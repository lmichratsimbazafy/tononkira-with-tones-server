package main

import (
	"lmich.com/tononkira/config"
	"lmich.com/tononkira/scripts/data"
)

func seed() {
	setupOption := new(config.SetupOptions)
	setupOption.Run = func() {
		data.UpsertSuperRoles()
		data.UpsertSuperAdminAdmin()
		authorData := data.GetAuthorsFromFile()
		data.UpsertAuthors(authorData)

		lyricsData := data.GetLyricsFromFile()
		data.UpsertLyrics(lyricsData)
	}
	setupOption.NeedLocalDb = true

	config.Setup(setupOption)

}

func main() {
	seed()
}

package helpers

import (
	"fmt"
	"path/filepath"
	"runtime"
)

func GetCurrentDirPath() string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		fmt.Println("Impossible de récupérer le chemin du fichier source.")
	}
	return filepath.Dir(filename)
}

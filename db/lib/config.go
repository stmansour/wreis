package db

import (
	"extres"
	"log"

	"github.com/kardianos/osext"
)

// DBConfig is the shared struct of configuration values
var DBConfig extres.ExternalResources

// ReadConfig will read the configuration file "config.json" if
// it exists in the current directory
func ReadConfig() error {
	folderPath, err := osext.ExecutableFolder()
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Printf("Executable folder = %s\n", folderPath)
	fname := folderPath + "/config.json"
	err = extres.ReadConfig(fname, &DBConfig)
	return err
}

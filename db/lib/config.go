package db

import (
	"database/sql"
	"extres"
	"fmt"
	"log"
	"math/rand"
	"time"
	util "wreis/util/lib"

	"github.com/kardianos/osext"
)

// Wdb is a struct with all variables needed by the db infrastructure
var Wdb struct {
	Prepstmt PrepSQL
	Config   extres.ExternalResources
	DB       *sql.DB
	DBFields map[string]string // map of db table fields DBFields[tablename] = field list
	Zone     *time.Location    // what timezone should the server use?
	Key      []byte            // crypto key
	Rand     *rand.Rand        // for generating Reference Numbers or other UniqueIDs
	noAuth   bool              // is authrization needed to access the db?
}

// ReadConfig will read the configuration file "config.json" if
// it exists in the current directory
func ReadConfig() error {
	folderPath, err := osext.ExecutableFolder()
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Printf("Executable folder = %s\n", folderPath)
	fname := folderPath + "/config.json"
	if err = extres.ReadConfig(fname, &Wdb.Config); err != nil {
		return err
	}

	Wdb.Zone, err = time.LoadLocation(Wdb.Config.Timezone)
	if err != nil {
		fmt.Printf("Error loading timezone %s : %s\n", Wdb.Config.Timezone, err.Error())
		util.Ulog("Error loading timezone %s : %s", Wdb.Config.Timezone, err.Error())
		return err
	}
	return err
}

// Init initializes the db subsystem
func Init(db *sql.DB) error {
	Wdb.DB = db
	Wdb.DBFields = map[string]string{}
	BuildPreparedStatements()

	return nil
}

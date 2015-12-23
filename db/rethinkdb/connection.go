package rethinkdb

import (
	"fmt"
	"log"
	"os"

	r "github.com/dancannon/gorethink"
)

// S is a pointer to a RethinkDB Session. Don't set it to something else!
var S *r.Session

// GetSession retyrbs a pointer to a `gorethink` Session struct which can be
// used across manu goroutines
func GetSession() *r.Session {
	return S
}

// ConnectDatabase establishes database session connection to RethinkDB using the
// the default localhost:28015 or from the environment variables
// DATABASE_PORT_28015_TCP_ADDR and DATABASE_PORT_28015_TCP_PORT for address and
// port respectively
func ConnectDatabase(name string) {
	var err error
	var dbhost, dbport string

	dbhost = os.Getenv("DATABASE_PORT_28015_TCP_ADDR")
	dbport = os.Getenv("DATABASE_PORT_28015_TCP_PORT")
	if dbhost != "" && dbport != "" {
		log.Println("Using environment variables DATABASE_PORT_28015_TCP ...")
	} else {
		dbhost = "localhost"
		dbport = "28015"
	}

	dbpeer := fmt.Sprintf("%s:%s", dbhost, dbport)
	log.Printf("Connecting to %s ...", dbpeer)

	S, err = r.Connect(r.ConnectOpts{
		Address:       dbpeer,
		Database:      name,
		DiscoverHosts: true,
	})
	if err != nil {
		log.Fatalf("Could not establish database connection: %s", err)
	}

	log.Println("Database connected")
}

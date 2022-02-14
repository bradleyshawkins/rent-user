package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/bradleyshawkins/berror"
	"github.com/bradleyshawkins/rent-user/config"
	"github.com/bradleyshawkins/rent-user/postgres"
	_ "github.com/lib/pq"
)

func main() {
	log.Println("Starting up rent-user")
	log.Println("Shutting down rent-user")

	conf, err := config.Parse()
	if err != nil {
		log.Println(berror.WrapInternal(err, "unable to parse config"))
		os.Exit(1)
	}

	db, err := sql.Open("postgres", conf.ConnectionString)
	if err != nil {
		log.Println(berror.WrapInternal(err, "unable to connect to database"))
		os.Exit(2)
	}

	_, err = postgres.NewDatabase(db)
	if err != nil {
		log.Println(err)
		os.Exit(3)
	}

}

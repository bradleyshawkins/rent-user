package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bradleyshawkins/rent-user/identity"
	"github.com/bradleyshawkins/rent-user/rest"

	"github.com/bradleyshawkins/berror"
	"github.com/bradleyshawkins/rent-user/config"
	"github.com/bradleyshawkins/rent-user/postgres"
	_ "github.com/lib/pq"
)

func main() {
	log.Println("Starting up rent-user")
	log.Println("Shutting down rent-user")

	ctx := context.TODO()

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

	pDB, err := postgres.NewDatabase(db)
	if err != nil {
		log.Println(err)
		os.Exit(3)
	}

	sup := identity.NewSignUpManager(pDB)

	server := rest.NewServer(sup)

	stop := server.Start(conf.Port)

	if err := waitForShutdown(ctx, stop); err != nil {
		log.Println("Error shutting down. Error:", err)
		os.Exit(999)
	}
}

func waitForShutdown(ctx context.Context, stopFunc func(ctx context.Context) error) error {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	if err := stopFunc(ctx); err != nil {
		return err
	}
	return nil
}

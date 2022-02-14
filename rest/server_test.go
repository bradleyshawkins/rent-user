package rest_test

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/bradleyshawkins/rent-user/config"

	"github.com/bradleyshawkins/rent-user/identity"

	"github.com/bradleyshawkins/rent-user/postgres"
	"github.com/bradleyshawkins/rent-user/rest"
)

var (
	server     *rest.Server
	serverAddr string
	httpClient *http.Client
)

func TestMain(m *testing.M) {
	// flag.Parse() must be called before testing.Short() or else it will panic
	flag.Parse()
	// Check to see if -short argument was used on go test to signify not to run integration tests
	if testing.Short() {
		log.Println("Skipping Integration Tests")
		os.Exit(0)
	}

	conf, err := config.Parse()
	if err != nil {
		log.Println("Unable to parse config")
		os.Exit(1)
	}

	// Create Database
	db, err := sql.Open("postgres", conf.ConnectionString)
	if err != nil {
		log.Println(err)
		os.Exit(2)
	}

	log.Println("Beginning integration tests")
	pg, err := postgres.NewDatabase(db)
	if err != nil {
		log.Println("Unable to create database connection. Error:", err)
		os.Exit(999)
	}

	// Create registrar
	sup := identity.NewSignUpManager(pg)

	// Create Router
	server = rest.NewServer(sup)

	svr := httptest.NewServer(server.Mux)
	serverAddr = svr.URL
	httpClient = svr.Client()

	code := m.Run()

	log.Println("Completed integration tests")
	svr.Close()
	os.Exit(code)
}

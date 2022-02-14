package config

import (
	"fmt"
	"log"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	ConnectionString string `envconfig:"DATABASE_URL" default:"postgresql://postgres:password@localhost:5432/rent-user?sslmode=disable"`
	Port             string `envconfig:"PORT" default:"8080"`
}

func Parse() (*Config, error) {
	log.Println("Parsing config")
	var c Config
	err := envconfig.Process("", &c)
	log.Printf("%+v\n", c)
	if err != nil {
		return nil, fmt.Errorf("unable to process config. Error:%v", err)
	}
	return &c, nil
}

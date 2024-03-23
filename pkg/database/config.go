package database

import (
	"fmt"
	"github.com/echovisionlab/aws-weather-updater/pkg/env"
	"os"
)

type config struct {
	Name string
	User string
	Pass string
	Host string
	Port string
}

func (d *config) ConnStr() string {
	var host, port = d.Host, d.Port
	if len(port) == 0 {
		port = "5432"
	}
	if len(host) == 0 {
		host = "localhost"
	}
	return fmt.Sprintf("user=%v dbname=%v sslmode=disable password=%v host=%v port=%v",
		d.User,
		d.Name,
		d.Pass,
		host,
		port)
}

func readDatabaseConfig() (*config, error) {
	if err := env.CheckRequiredEnv(); err != nil {
		return nil, err
	}
	return &config{
		Name: os.Getenv(env.DatabaseName),
		User: os.Getenv(env.DatabaseUser),
		Pass: os.Getenv(env.DatabasePass),
		Host: os.Getenv(env.DatabaseHost),
		Port: os.Getenv(env.DatabasePort),
	}, nil
}

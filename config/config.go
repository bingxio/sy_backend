package config

import (
	"os"

	"github.com/bingxio/dotenv"
)

type Config struct {
	ApiPort      string
	User         string
	Password     string
	Host         string
	Port         int
	DbPath       string `env:"db_path"`
	ResourcePath string `env:"resource_path"`
}

var Conf = new(Config)

func LoadConf() error {
	buffer, err := os.ReadFile(".env")
	if err != nil {
		return err
	}
	return dotenv.Unmarshal(buffer, Conf)
}

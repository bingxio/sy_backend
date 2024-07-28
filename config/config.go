package config

import (
	"os"

	"github.com/bingxio/dotenv"
)

type Config struct {
	ApiPort  string
	User     string
	Password string
	Host     string
	Port     int
	DbPath   string `env:"db_path"`
}

var C = new(Config)

func LoadConf() error {
	buffer, err := os.ReadFile(".env")
	if err != nil {
		return err
	}
	return dotenv.Unmarshal(buffer, C)
}

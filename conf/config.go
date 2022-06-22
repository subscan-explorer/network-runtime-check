package conf

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

var Conf Config

type Config struct {
	Network []string `yaml:"network"`
	APIKey  string   `yaml:"apikey"`
}

func InitConf(path string) {
	fd, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("failed to open configuration file. err: %s\n", err.Error())
	}
	if err = yaml.Unmarshal(fd, &Conf); err != nil {
		log.Fatalf("failed to parse configuration file. err: %s\n", err.Error())
	}
}

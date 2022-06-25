package conf

import (
	"context"
	"log"
	"os"

	"github.com/subscan-explorer/network-runtime-check/internal/api/github/self"
	"gopkg.in/yaml.v3"
)

var Conf Config

type Config struct {
	Network []string `yaml:"network"`
	APIKey  string   `yaml:"apikey"`
}

func InitConf(ctx context.Context, path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		if data, err = self.GetConfigData(ctx); err != nil {
			log.Fatalln("failed to open configuration file.")
			return
		}
	}
	if err = yaml.Unmarshal(data, &Conf); err != nil {
		log.Fatalf("failed to parse configuration file. err: %s\n", err.Error())
	}
}

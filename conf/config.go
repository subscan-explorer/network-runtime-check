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
	Network   []string    `yaml:"network"`
	APIKey    string      `yaml:"apikey"`
	ParamRule []ParamRule `yaml:"param_rule"`
}

type ParamRule struct {
	Domain      string `yaml:"domain"`
	WsAddr      string `yaml:"wsAddr"`
	RuleInherit string `yaml:"rule_inherit"`
	Pallet      []Rule `yaml:"pallet"`
}

type Rule struct {
	Name  string `yaml:"name"`
	Event []struct {
		Name  string   `yaml:"name"`
		Param []string `yaml:"param"`
	} `yaml:"event"`
	Extrinsic []struct {
		Name  string   `yaml:"name"`
		Param []string `yaml:"param"`
	}
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

	// rule inherit
	ruleMap := make(map[string][]Rule)
	for _, rule := range Conf.ParamRule {
		if len(rule.RuleInherit) != 0 {
			ruleMap[rule.RuleInherit] = nil
		}
	}
	for _, rule := range Conf.ParamRule {
		if _, ok := ruleMap[rule.Domain]; ok {
			ruleMap[rule.Domain] = rule.Pallet
		}
	}
	for i := 0; i < len(Conf.ParamRule); i++ {
		if d := Conf.ParamRule[i].RuleInherit; len(d) != 0 {
			Conf.ParamRule[i].Pallet = ruleMap[d]
		}
	}
}

package conf

import (
	"log"
	"net/url"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type Rule struct {
	Network []NetworkRule `yaml:"network"`
	Rule    []struct {
		Name   string   `yaml:"name"`
		Pallet []Pallet `yaml:"pallet"`
	} `yaml:"rule"`
}

type NetworkRule struct {
	Name        string   `yaml:"name"`
	Domain      string   `yaml:"domain"`
	WsAddr      string   `yaml:"wsAddr"`
	RuleInherit []string `yaml:"rule_inherit"`
	Pallet      []Pallet `yaml:"pallet"`
}

type Pallet struct {
	Name      string      `yaml:"name"`
	Event     []Event     `yaml:"event"`
	Extrinsic []Extrinsic `yaml:"extrinsic"`
}

type Event struct {
	Name  string   `yaml:"name"`
	Param []string `yaml:"param"`
}
type Extrinsic struct {
	Name  string   `yaml:"name"`
	Param []string `yaml:"param"`
}

func LoadRule(path string) Rule {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalln("failed to loading rule file.")
	}
	var rule Rule
	if err = yaml.Unmarshal(data, &rule); err != nil {
		log.Fatalf("failed to parse rule file. err: %s\n", err.Error())
	}
	// rule inherit
	ruleMap := make(map[string][]Pallet)

	for _, r := range rule.Rule {
		ruleMap[r.Name] = r.Pallet
	}

	for i := 0; i < len(rule.Network); i++ {
		if len(rule.Network[i].Name) == 0 {
			if len(rule.Network[i].Domain) != 0 {
				rule.Network[i].Name = rule.Network[i].Domain
			} else if len(rule.Network[i].WsAddr) != 0 {
				if r, err := url.Parse(rule.Network[i].WsAddr); err == nil {
					rule.Network[i].Name = r.Host
				}
			}
		}
		if len(rule.Network[i].Name) == 0 {
			log.Fatalln("The network name in the rule file is required")
		}
		rule.Network[i].Name = strings.ToLower(rule.Network[i].Name)
		if len(rule.Network[i].RuleInherit) > 0 {
			rule.Network[i].Pallet = mergeRule(rule.Network[i].Pallet, ruleMap, rule.Network[i].RuleInherit)
		}
	}
	return rule
}

func mergeRule(baseRule []Pallet, ruleMap map[string][]Pallet, ruleName []string) []Pallet {
	eventMap := make(map[string]map[string][]string)
	extMap := make(map[string]map[string][]string)
	for _, pallet := range baseRule {
		em := make(map[string][]string)
		for _, e := range pallet.Event {
			em[e.Name] = e.Param
		}
		eventMap[pallet.Name] = em
		etm := make(map[string][]string)
		for _, e := range pallet.Extrinsic {
			etm[e.Name] = e.Param
		}
		extMap[pallet.Name] = etm
	}

	for _, rn := range ruleName {
		if m, ok := ruleMap[rn]; ok {
			for _, r := range m {
				em := eventMap[r.Name]
				if em == nil {
					em = make(map[string][]string)
				}
				for _, e := range r.Event {
					if _, ok := em[e.Name]; !ok {
						em[e.Name] = e.Param
					}
				}
				eventMap[r.Name] = em
				etm := extMap[r.Name]
				if etm == nil {
					etm = make(map[string][]string)
				}
				for _, e := range r.Extrinsic {
					if _, ok := etm[e.Name]; !ok {
						etm[e.Name] = e.Param
					}
				}
				extMap[r.Name] = etm
			}
		}
	}
	var pallet []Pallet
	for pl, m := range eventMap {
		p := Pallet{Name: pl}
		for e, param := range m {
			p.Event = append(p.Event, Event{Name: e, Param: param})
		}
		for e, param := range extMap[pl] {
			p.Extrinsic = append(p.Extrinsic, Extrinsic{Name: e, Param: param})
		}
		pallet = append(pallet, p)
	}
	return pallet
}

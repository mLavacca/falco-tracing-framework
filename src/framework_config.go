package main

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type TracerConfigurations struct {
	Record           RecordConfiguration            `yaml:"record,omitempty"`
	Report           ReportConfiguration            `yaml:"report,omitempty"`
	RulesBreakers    []RulesBreakersConfiguration   `yaml:"rules_breakers"`
	BreakingProfiles []BreakingProfileConfiguration `yaml:"breaking_profiles"`
}

type RecordConfiguration struct {
	ProgConfig      ProgConfiguration `yaml:"prog_config,omitempty"`
	BreakingProfile string            `yaml:"breaking_profile,omitempty"`
}

type ReportConfiguration struct {
	ProgConfig ProgConfiguration `yaml:"prog_config,omitempty"`
}

type ProgConfiguration struct {
	ProgBin  string   `yaml:"prog_bin,omitempty"`
	ProgArgs []string `yaml:"prog_args,omitempty"`
}

type RulesBreakersConfiguration struct {
	Rule     string `yaml:"rule,omitempty"`
	RuleId   int    `yaml:"rule_id,omitempty"`
	Trigger  string `yaml:"trigger,omitempty"`
	Rollback string `yaml:"rollback,omitempty"`
}

type BreakingProfileConfiguration struct {
	Name     string `yaml:"name,omitempty"`
	Sequence []int  `yaml:"sequence,omitempty"`
	Ratio    int    `yaml:"ratio,omitempty"`
	Limit    int    `yaml:"limit,omitempty"`
}

func (tc *TracerConfigurations) UnmarshalYAML(configFile string) error {

	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(data, tc)

	return err
}

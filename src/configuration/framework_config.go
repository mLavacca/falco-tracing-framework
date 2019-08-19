package configuration

import (
	"io/ioutil"
	"time"

	"gopkg.in/yaml.v2"
)

type TracerConfigurations struct {
	Record           RecordConfiguration            `yaml:"record,omitempty"`
	OfflineReport    OfflineReportConfiguration     `yaml:"offline_report,omitempty"`
	OnlineReport     OnlineReportConfiguration      `yaml:"online_report,omitempty"`
	RulesBreakers    []RulesBreakersConfiguration   `yaml:"rules_breakers"`
	BreakingProfiles []BreakingProfileConfiguration `yaml:"breaking_profiles"`
}

type RecordConfiguration struct {
	ProgConfig      ProgConfiguration `yaml:"prog_config,omitempty"`
	BreakingProfile string            `yaml:"breaking_profile,omitempty"`
}

type OfflineReportConfiguration struct {
	ProgConfig             ProgConfiguration `yaml:"prog_config,omitempty"`
	OutputFile             string            `yaml:"output_file,omitempty"`
	OutputFoldedStacktrace string            `yaml:"output_folded_stacktrace,omitempty"`
	Iterations             int               `yaml:"iterations,omitempty"`
}

type OnlineReportConfiguration struct {
	ProgConfig   ProgConfiguration `yaml:"prog_config,omitempty"`
	OutputFile   string            `yaml:"output_file,omitempty"`
	PullInterval time.Duration     `yaml:"pull_interval,omitempty"`
}

type ProgConfiguration struct {
	ProgBin  string   `yaml:"prog_bin,omitempty"`
	ProgArgs []string `yaml:"prog_args,omitempty"`
}

type RulesBreakersConfiguration struct {
	Rule   string `yaml:"rule,omitempty"`
	RuleId int    `yaml:"rule_id,omitempty"`
}

type BreakingProfileConfiguration struct {
	Name             string  `yaml:"name,omitempty"`
	Sequence         [][]int `yaml:"sequence,omitempty"`
	RollbackSequence []int   `yaml:"rollback_sequence,omitempty"`
	Ratio            int     `yaml:"ratio,omitempty"`
	Limit            int     `yaml:"limit,omitempty"`
	Duration         int     `yaml:"duration,omitempty"`
}

func (tc *TracerConfigurations) UnmarshalYAML(configFile string) error {

	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(data, tc)

	return err
}

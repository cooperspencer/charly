package types

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"time"

	"gopkg.in/yaml.v3"
)

type Repo struct {
	ID     int
	Name   string
	Url    string
	SshUrl string
	User   string
	Branch string
	Commit string
}

type Config struct {
	Configuration struct {
		DBFile string `yaml:"db-file"`
		Cron   string `yaml:"cron"`
	} `yaml:"configuration"`
	Github VCS `yaml:"github"`
	Gitlab VCS `yaml:"gitlab"`
	Gitea  VCS `yaml:"gitea"`
	Gogs   VCS `yaml:"gogs"`
}

type Repos struct {
	ID         int               `yaml:"id"`
	URL        string            `yaml:"url"`
	User       string            `yaml:"user"`
	Repo       string            `yaml:"repo"`
	Token      string            `yaml:"token"`
	Branch     string            `yaml:"branch"`
	WorkingDir string            `yaml:"working-dir"`
	Script     string            `yaml:"script"`
	Variables  map[string]string `yaml:"variables"`
}

type VCS struct {
	Token string  `yaml:"token"`
	URL   string  `yaml:"url"`
	Repos []Repos `yaml:"repos"`
}

func (c *Config) Get(path string) error {
	cfgdata, err := ioutil.ReadFile(path)

	if err != nil {
		return err
	}

	err = yaml.Unmarshal(cfgdata, &c)

	if err != nil {
		return err
	}

	return nil
}

func (c Config) MissingCronSpec() bool {
	return c.Configuration.Cron == ""
}

func ParseCronSpec(spec string) (cron.Schedule, error) {
	sched, err := cron.ParseStandard(spec)

	if err != nil {
		log.Error().Str("spec", spec).Msg(err.Error())
	}

	return sched, err
}

func (c Config) GetNextRun() (*time.Time, error) {
	if c.MissingCronSpec() {
		return nil, fmt.Errorf("cron unspecified")
	}
	parsedSched, err := ParseCronSpec(c.Configuration.Cron)
	if err != nil {
		return nil, err
	}
	next := parsedSched.Next(time.Now())
	return &next, nil
}

func (c Config) HasValidCronSpec() bool {
	if c.MissingCronSpec() {
		return false
	}

	_, err := ParseCronSpec(c.Configuration.Cron)

	return err == nil
}

package types

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"time"

	"gopkg.in/yaml.v3"
)

type Repos struct {
	URL        string            `yaml:"url"`
	Auth       Auth              `yaml:"auth"`
	Branch     string            `yaml:"branch"`
	WorkingDir string            `yaml:"working-dir"`
	Script     string            `yaml:"script"`
	Variables  map[string]string `yaml:"variables"`
}

type Auth struct {
	SSHKeyfile     string `yaml:"ssh-keyfile"`
	SSHKeyPassword string `yaml:"ssh-key-password"`
	Token          string `yaml:"token"`
	Username       string `yaml:"username"`
	Password       string `yaml:"password"`
}

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
	VCSRepos `yaml:",inline"`
}

type VCSRepos map[string]VCS

type VCS struct {
	Auth       Auth    `yaml:"auth"`
	WorkingDir string  `yaml:"working-dir"`
	Repos      []Repos `yaml:"repos"`
}

func (a *Auth) Fill(auth *Auth) {
	if a.Username == "" {
		a.Username = auth.Username
	}
	if a.Password == "" {
		a.Password = auth.Password
	}
	if a.Token == "" {
		a.Token = auth.Token
	}
	if a.SSHKeyfile == "" {
		a.SSHKeyfile = a.SSHKeyfile
	}
	if a.SSHKeyPassword == "" {
		a.SSHKeyPassword = auth.SSHKeyPassword
	}
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

package main

import (
	"charly/db"
	"charly/github"
	"charly/gitlab"
	"charly/scripts"
	"charly/types"
	"fmt"
	"os"
	"path/filepath"

	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	debug      = kingpin.Flag("debug", "enable debug mode").Bool()
	configfile = kingpin.Arg("configfile", "the configuration file").Default("conf.yaml").String()
	exPath     string
	d          db.DB
	version    = "0.0.1"
)

func dorepo(repo types.Repos, token, vcs string) {
	if repo.Token == "" {
		repo.Token = token
	}
	log.Debug().Str("vcs", vcs).Msg(fmt.Sprint(repo))

	if repo.WorkingDir == "" {
		repo.WorkingDir = exPath
	}
	log.Info().Str("vcs", vcs).Msgf("checking %s...", repo.Repo)
	r := types.Repo{}
	switch vcs {
	case "github":
		re, err := github.GetRepo(repo.User, repo.Repo, repo.Branch, repo.Token)
		if err != nil {
			log.Fatal().Str("vcs", vcs).Str("branch", repo.Branch).Err(err)
		}
		r = re
	case "gitlab":
		re, err := gitlab.GetRepo(repo.Branch, repo.Token, repo.URL, repo.ID)
		if err != nil {
			log.Fatal().Str("vcs", vcs).Str("branch", repo.Branch).Err(err)
		}
		r = re
	}

	log.Debug().Str(vcs, vcs).Msg(fmt.Sprint(r))
	dbrepo, err := d.GetRepo(r)
	if err != nil {
		if err.Error() == "repos bucket doesn't exist" {
			log.Debug().Str("vcs", vcs).Str("branch", r.Branch).Err(err)
		} else {
			log.Fatal().Str("vcs", vcs).Str("branch", r.Branch).Err(err)
		}
	}

	if dbrepo.Commit != r.Commit {
		log.Info().Str("vcs", vcs).Str("branch", r.Branch).Msgf("%s was updated", r.Name)
		if repo.Script != "" {
			log.Info().Str("vcs", vcs).Str("branch", r.Branch).Msgf("run script for %s", r.Name)
			err = scripts.RunScript(r, repo)
			if err != nil {
				log.Fatal().Str("vcs", vcs).Str("branch", r.Branch).Err(err)
			}
		}
		err = d.InsertRepo(r)
		if err != nil {
			log.Fatal().Str("vcs", vcs).Str("branch", r.Branch).Err(err)
		}
	}
}

func logNextRun(conf *types.Config) {
	nextRun, err := conf.GetNextRun()
	if err == nil {
		log.Info().
			Str("next", nextRun.String()).
			Str("cron", conf.Configuration.Cron).
			Msg("Next cron run")
	}
}

func PlaysForever() {
	wait := make(chan struct{})
	for {
		<-wait
	}
}

func main() {
	kingpin.Version(version)
	kingpin.Parse()
	timeformat := "2006-01-02T15:04:05Z07:00"
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: timeformat,
	})

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	ex, err := os.Executable()
	if err != nil {
		log.Fatal().Err(err)
	}
	exPath = filepath.Dir(ex)

	config := types.Config{}
	err = config.Get(*configfile)
	if err != nil {
		log.Panic().Err(err)
	}
	if config.Configuration.DBFile == "" {
		config.Configuration.DBFile = "charly.db"
	}
	log.Debug().Str("config", *configfile).Msg(fmt.Sprint(config))
	d, err = db.New(config.Configuration.DBFile, 0664)
	if err != nil {
		log.Fatal().Err(err)
	}

	if config.HasValidCronSpec() {
		c := cron.New()
		logNextRun(&config)

		c.AddFunc(config.Configuration.Cron, func() {
			RunRepos(config)
		})
		PlaysForever()
	} else {
		RunRepos(config)
	}
}

func RunRepos(config types.Config) {
	for _, repo := range config.Github.Repos {
		dorepo(repo, config.Github.Token, "github")
	}
	for _, repo := range config.Gitlab.Repos {
		if repo.URL == "" {
			repo.URL = config.Gitlab.URL
		}
		if repo.Repo == "" {
			repo.Repo = fmt.Sprintf("ID %d", repo.ID)
		}
		dorepo(repo, config.Gitlab.Token, "gitlab")
	}
}

package main

import (
	"charly/db"
	"charly/git"
	"charly/scripts"
	"charly/types"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5/plumbing"
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

func getName(location string) string {
	if f, err := os.Stat(location); !os.IsNotExist(err) {
		return f.Name()
	} else {
		lastelement := location[strings.LastIndex(location, "/")+1:]
		if strings.HasSuffix(lastelement, ".git") {
			return lastelement[:strings.LastIndex(lastelement, ".git")]
		} else {
			return lastelement
		}
	}
}

func checkrepo(repo types.Repos, re types.Repo, data []*plumbing.Reference) {
	log.Debug().Str("url", repo.URL).Str("branch", repo.Branch).Msg("gathering info for the repo...")
	for _, d := range data {
		if d.Name().Short() == repo.Branch {
			re = types.Repo{
				ID:     0,
				Name:   getName(repo.URL),
				Url:    repo.URL,
				User:   repo.Auth.Username,
				Branch: d.Name().Short(),
				Commit: d.Hash().String(),
			}
			log.Debug().Str("url", repo.URL).Str("branch", repo.Branch).Msg(fmt.Sprint(re))
			break
		}
	}
	dbrepo, err := d.GetRepo(re)
	if err != nil {
		if err.Error() == "repos bucket doesn't exist" {
			log.Debug().Str("branch", re.Branch).Err(err)
		} else {
			log.Fatal().Str("branch", re.Branch).Err(err)
		}
	}

	if dbrepo.Commit != re.Commit {
		log.Info().Str("branch", re.Branch).Msgf("%s was updated", re.Name)
		if repo.Script.Code != "" {
			log.Info().Str("branch", re.Branch).Msgf("run script for %s", re.Name)
			err = scripts.RunScript(re, repo)
			if err != nil {
				log.Fatal().Str("branch", re.Branch).Err(err)
			}
		}
		err = d.InsertRepo(re)
		if err != nil {
			log.Fatal().Str("branch", re.Branch).Err(err)
		}
	}
}

func dorepo(repo types.Repos) {
	log.Info().Str("url", repo.URL).Msg("checking for updates...")
	re := types.Repo{}
	data, err := git.GetRepo(repo)
	if err != nil {
		log.Error().Str("url", repo.URL).Msg(err.Error())
		return
	}
	if repo.Branch == "" {
		if repo.AllBranches {
			log.Info().Str("url", repo.URL).Msg("checking all branches...")
		} else {
			log.Info().Str("url", repo.URL).Msg("checking for default branch...")
		}
		for _, d := range data {
			if repo.AllBranches {
				if d.Name().IsBranch() {
					fmt.Println(d.Name().Short())
					repo.Branch = d.Name().Short()
					log.Debug().Str("url", repo.URL).Msgf("branch is %s...", repo.Branch)
					checkrepo(repo, re, data)
				}
			} else {
				if d.Hash().IsZero() {
					repo.Branch = d.Target().Short()
					log.Debug().Msgf("branch is %s", repo.Branch)
					log.Debug().Str("url", repo.URL).Msgf("default branch is %s...", repo.Branch)
					checkrepo(repo, re, data)
					break
				}
			}
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
			logNextRun(&config)
		})
		c.Start()
		PlaysForever()
	} else {
		RunRepos(config)
	}
}

func RunRepos(config types.Config) {
	for vcs, repos := range config.VCSRepos {
		log.Info().Msgf("running for %s...", vcs)
		for _, repo := range repos.Repos {
			repo.Auth.Fill(&repos.Auth)
			if repo.WorkingDir == "" {
				repo.WorkingDir = repos.WorkingDir
			}
			if repo.Script.Template != "" {
				if repo.Script.Code == "" {
					repo.Script.Code = config.Scripts[repo.Script.Template]
				} else {
					log.Warn().Msg("values are set in template and code! will use the value from code!")
				}
			}
			dorepo(repo)
		}
	}
}

package github

import (
	"charly/types"
	"context"
	"github.com/google/go-github/v42/github"
	"golang.org/x/oauth2"
)

func GetRepo(user, repo, branch, token string) (types.Repo, error) {
	client := github.NewClient(nil)
	if token == "" {
		client = github.NewClient(nil)
	} else {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		)
		tc := oauth2.NewClient(context.TODO(), ts)
		client = github.NewClient(tc)
	}
	origin, _, err := client.Repositories.Get(context.Background(), user, repo)
	if err != nil {
		return types.Repo{}, err
	}

	if branch == "" {
		branch = origin.GetDefaultBranch()
	}

	r, _, err := client.Repositories.GetBranch(context.Background(), user, repo, branch, true)
	return types.Repo{0, repo, origin.GetCloneURL(), origin.GetSSHURL(), user, r.GetName(), r.GetCommit().GetSHA()}, err
}

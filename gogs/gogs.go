package gogs

import (
	"charly/types"
	"github.com/gogs/go-gogs-client"
)

func GetRepo(user, repo, branch, token, url string) (types.Repo, error) {
	r := types.Repo{}
	client := gogs.NewClient(url, token)

	gr, err := client.GetRepo(user, repo)
	if err != nil {
		return r, err
	}

	if branch == "" {
		branch = gr.DefaultBranch
	}

	gb, err := client.GetRepoBranch(user, repo, branch)
	if err != nil {
		return r, err
	}

	r = types.Repo{
		ID:     0,
		Name:   gr.Name,
		Url:    gr.CloneURL,
		SshUrl: gr.SSHURL,
		User:   gr.Owner.UserName,
		Branch: gb.Name,
		Commit: gb.Commit.ID,
	}

	return r, err
}

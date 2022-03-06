package gitea

import (
	"charly/types"
	"code.gitea.io/sdk/gitea"
)

func GetRepo(user, repo, branch, token, url string) (types.Repo, error) {
	if url == "" {
		url = "https://gitea.com"
	}
	r := types.Repo{}
	client := &gitea.Client{}
	var err error
	if token != "" {
		client, err = gitea.NewClient(url, gitea.SetToken(token))
	} else {
		client, err = gitea.NewClient(url)
	}

	if err != nil {
		return r, err
	}

	gr, _, err := client.GetRepo(user, repo)
	if err != nil {
		return r, err
	}

	if branch == "" {
		branch = gr.DefaultBranch
	}

	gb, _, err := client.GetRepoBranch(user, repo, branch)
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

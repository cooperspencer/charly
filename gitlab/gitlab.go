package gitlab

import (
	"charly/types"
	"github.com/xanzy/go-gitlab"
)

func GetRepo(branch, token, url string, id int) (types.Repo, error) {
	client := &gitlab.Client{}
	r := types.Repo{}
	var err error
	if url == "" {
		url = "https://gitlab.com"
		client, err = gitlab.NewClient(token)
	} else {
		client, err = gitlab.NewClient(token, gitlab.WithBaseURL(url))
	}
	if err != nil {
		return r, err
	}

	project, _, err := client.Projects.GetProject(id, &gitlab.GetProjectOptions{})
	if err != nil {
		return r, err
	}
	if branch == "" {
		branch = project.DefaultBranch
	}

	gb, _, err := client.Branches.GetBranch(id, branch)
	if err != nil {
		return r, nil
	}

	r = types.Repo{
		ID:     id,
		Name:   project.Name,
		Url:    project.HTTPURLToRepo,
		SshUrl: project.SSHURLToRepo,
		User:   project.Namespace.Name,
		Branch: project.DefaultBranch,
		Commit: gb.Commit.ID,
	}

	return r, err
}

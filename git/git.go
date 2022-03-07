package git

import (
	"charly/types"
	"github.com/go-git/go-git/v5/plumbing"
	"os"
	"path"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
)

func GetRepo(repo types.Repos) ([]*plumbing.Reference, error) {
	var auth transport.AuthMethod
	var err error
	if strings.HasPrefix(repo.URL, "http://") || strings.HasPrefix(repo.URL, "https://") {
		if repo.Auth.Token != "" {
			auth = &http.BasicAuth{
				Username: "xyz",
				Password: repo.Auth.Token,
			}
		} else if repo.Auth.Username != "" && repo.Auth.Password != "" {
			auth = &http.BasicAuth{
				Username: repo.Auth.Username,
				Password: repo.Auth.Password,
			}
		}
	} else {
		if repo.Auth.SSHKeyfile == "" {
			home := os.Getenv("HOME")
			repo.Auth.SSHKeyfile = path.Join(home, ".ssh", "id_rsa")
		}
		auth, err = ssh.NewPublicKeysFromFile("git", repo.Auth.SSHKeyfile, repo.Auth.SSHKeyPassword)
		if err != nil {
			return nil, err
		}
	}

	rem := git.NewRemote(nil, &config.RemoteConfig{Name: "origin", URLs: []string{repo.URL}})
	data, err := rem.List(&git.ListOptions{Auth: auth})
	if err != nil {
		return nil, err
	}

	return data, nil
}

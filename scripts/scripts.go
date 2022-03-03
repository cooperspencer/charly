package scripts

import (
	"bufio"
	"charly/types"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"time"
)

func RunScript(repo types.Repo, repos types.Repos) error {
	os.Chdir(repos.WorkingDir)
	for k, v := range repos.Variables {
		os.Setenv(k, v)
	}

	os.Setenv("COMMIT", repo.Commit)
	os.Setenv("BRANCH", repo.Branch)
	os.Setenv("URL", repo.Url)
	os.Setenv("SSHURL", repo.SshUrl)
	os.Setenv("USER", repo.User)
	os.Setenv("REPO", repo.Name)

	f, err := ioutil.TempFile("", fmt.Sprintf("charly-%d.sh", time.Now().Unix()))
	defer os.Remove(f.Name())
	if err != nil {
		return err
	}
	_, err = f.WriteString(repos.Script)
	if err != nil {
		return err
	}

	cmd := exec.Command("/bin/sh", f.Name())
	stderr, _ := cmd.StderrPipe()
	err = cmd.Run()
	if err != nil {
		return err
	}

	stdread := bufio.NewReader(stderr)
	outputbytes, _ := ioutil.ReadAll(stdread)
	if len(string(outputbytes)) > 0 {
		return errors.New(string(outputbytes))
	}

	return nil
}

package sys

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

func RunExecutable(path string, gameserver string, cookie string) (err error) {
	log.WithFields(log.Fields{
		"path":       path,
		"gameserver": gameserver,
		"cookie":     cookie,
	}).Trace("starting toontown rewritten")

	dir := filepath.Dir(path)
	log.WithField("dir", dir).Trace("changing cwd")

	os.Chdir(dir)

	env := os.Environ()
	env = append(env, fmt.Sprintf("TTR_GAMESERVER=%s", gameserver))
	env = append(env, fmt.Sprintf("TTR_PLAYCOOKIE=%s", cookie))

	cmd := &exec.Cmd{
		Path:   path,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Env:    env,
	}

	if err = cmd.Start(); err != nil {
		return fmt.Errorf("could not start toontown rewritten: %s", err)
	}

	return
}

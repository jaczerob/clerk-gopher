package sys

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

func RunExecutable(path string, gameserver string, cookie string, pipe bool) (err error) {
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
		Path: path,
		Env:  env,
	}

	if pipe {
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			return fmt.Errorf("error obtaining stdout pipe: %s", err)
		}

		if err = cmd.Start(); err != nil {
			return fmt.Errorf("could not start toontown rewritten: %s", err)
		}

		in := bufio.NewScanner(stdout)
		for in.Scan() {
			log.Trace(in.Text())
		}

	} else {
		if err = cmd.Start(); err != nil {
			return fmt.Errorf("could not start toontown rewritten: %s", err)
		}
	}

	return
}

package executer

import (
	"fmt"
	"strings"

	"github.com/ahmadwaleed/chore/pkg/config"
	"github.com/ahmadwaleed/chore/pkg/ssh"
	"golang.org/x/sync/errgroup"
)

type Parallel struct {
	config *ssh.Config
}

func (e Parallel) Run(task config.Task, callback OnRunCallback) error {
	g := new(errgroup.Group)
	for i := range task.Hosts {
		srv := task.Hosts[i]
		g.Go(func() error {
			if srv.IsLocalhost() {
				if err := New("local").Run(task, callback); err != nil {
					return err
				}

				return nil
			}

			cmd, err := NewCmd()
			if err != nil {
				return err
			}
			cmd.Env = append(cmd.Env, task.Env.ToStringVars()...)

			out := cmd.Output()
			out.Host = srv.Host
			out.Command = strings.Join(task.Commands, " ")

			bash, err := findExecutable("bash")
			if err != nil {
				return err
			}

			ssh, err := findExecutable("ssh")
			if err != nil {
				return err
			}

			script := fmt.Sprintf(
				"%s %s '%s' -se <<EOF\n set -e\n%s\nEOF",
				ssh, srv.Target(), bash, strings.Join(task.Commands, "\n"),
			)
			err = cmd.Run(script)
			callback(out)

			return err
		})
	}

	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}

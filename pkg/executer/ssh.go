package executer

import (
	"fmt"
	"strings"

	"github.com/ahmadwaleed/chore/pkg/config"
)

type SSH struct{}

func (e SSH) Run(task config.Task, callback OnRunCallback) error {
	for _, srv := range task.Hosts {
		if srv.IsLocalhost() {
			if err := New("local").Run(task, callback); err != nil {
				return err
			}

			continue
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
	}

	return nil
}

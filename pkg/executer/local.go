package executer

import (
	"fmt"
	"strings"

	"github.com/ahmadwaleed/chore/pkg/config"
)

type Local struct{}

func (e Local) Run(task config.Task, callback OnRunCallback) error {
	cmd, err := NewCmd()
	if err != nil {
		return err
	}
	cmd.Env = append(cmd.Env, task.Env.ToStringVars()...)
	out := cmd.Output()
	out.Host = "localhost"
	out.Command = strings.Join(task.Commands, " ")

	bash, err := findExecutable("bash")
	if err != nil {
		return err
	}

	script := fmt.Sprintf("%s -se <<EOF\n set -e\n%s\nEOF", bash, strings.Join(task.Commands, "\n"))
	err = cmd.Run(script)
	callback(out)

	return err
}

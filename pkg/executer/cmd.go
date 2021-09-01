package executer

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"mvdan.cc/sh/expand"
	"mvdan.cc/sh/interp"
	"mvdan.cc/sh/syntax"
)

type CmdOutput struct {
	Host    string
	Command string
	Stdout  bytes.Buffer
	Stderr  bytes.Buffer
}

func (o *CmdOutput) Display() {
	fmt.Println(string("\033[34m"), fmt.Sprintf("[%s](%s):", o.Host, o.Command), string("\033[0m"))
	fmt.Println(o.Stdout.String())

	if o.Stderr.String() != "" {
		fmt.Println(string("\033[31m"), fmt.Sprintf("[%s]:", o.Host), string("\033[0m"))
		fmt.Fprintf(os.Stderr, "%s", o.Stderr.String())
	}
}

type Cmd struct {
	Env    []string
	runner *interp.Runner
	output *CmdOutput
}

func (c *Cmd) Run(cmd string) error {
	file, err := syntax.NewParser().Parse(strings.NewReader(cmd), "")
	if err != nil {
		return err
	}

	c.runner.Env = expand.ListEnviron(c.Env...)

	if err := c.runner.Run(context.TODO(), file); err != nil {
		return err
	}

	return nil
}

func (c *Cmd) Output() *CmdOutput {
	return c.output
}

func NewCmd() (*Cmd, error) {
	c := &Cmd{}
	c.output = &CmdOutput{}

	runner, err := interp.New(
		interp.Env(expand.ListEnviron(c.Env...)),
		interp.StdIO(nil, &c.output.Stdout, &c.output.Stderr),
	)

	if err != nil {
		return nil, err
	}

	c.runner = runner

	return c, nil
}

func findExecutable(name string) (string, error) {
	exec, err := exec.LookPath(name)
	if err != nil {
		return "", fmt.Errorf("Could not find %s executable. %v", name, err)
	}

	return exec, nil
}

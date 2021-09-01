package cli

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/ahmadwaleed/chore/pkg/config"
	"github.com/ahmadwaleed/chore/pkg/executer"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const defaultConfigFile = "taskfile.yaml"

func NewRunCommand() *cobra.Command {
	opts := NewRunOptions()

	c := &cobra.Command{
		Use:   "run",
		Short: "Run tasks",
		Long:  "Run the tasks defined in Commet file.",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if err := run(args[0], opts); err != nil {
				log.Fatalf(err.Error())
			}
		},
	}

	opts.BindFlags(c.Flags())

	return c
}

type RunOptions struct {
	Task     string
	IsBucket bool
	Continue bool
	DryRun   bool
	Parallel bool
	Path     string
	Filename string
}

func NewRunOptions() *RunOptions {
	return &RunOptions{}
}

func (opt *RunOptions) BindFlags(flags *pflag.FlagSet) {
	flags.BoolVar(&opt.Continue, "continue", false, "Continue running even if a task fails")
	flags.BoolVar(&opt.DryRun, "dry-run", false, "Dump Bash script for inspection")
	flags.BoolVar(&opt.IsBucket, "bucket", false, "Run the bucket of tasks")
	flags.BoolVar(&opt.Parallel, "parallel", false, "Run task concurrently on servers")
	flags.StringVar(&opt.Path, "path", "", "The path to the Commet.yaml file")
	flags.StringVar(&opt.Filename, "filename", defaultConfigFile, "The name of the Commet file")
}

func (opt *RunOptions) configFile() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	if opt.Path == "" {
		return filepath.Join(dir, opt.Filename), nil
	}

	return filepath.Join(opt.Path, opt.Filename), nil
}

func run(task string, opts *RunOptions) error {
	file, err := opts.configFile()
	if err != nil {
		return err
	}

	if _, err := os.Stat(file); err != nil {
		return err
	}

	conf, err := config.NewConfig(file)
	if err != nil {
		return err
	}

	if !opts.IsBucket {
		if err := runTask(task, conf.Tasks, opts); err != nil {
			return err
		}
	} else {
		if err := runBucket(task, conf.Buckets, conf.Tasks, opts); err != nil {
			return err
		}
	}

	return nil
}

func runBucket(name string, buckets []*config.Bucket, tasks []*config.Task, opts *RunOptions) error {
	var found bool
	for _, bucket := range buckets {
		if bucket.Name == name {
			found = true
			for _, task := range bucket.Tasks {
				if err := runTask(task, tasks, opts); err != nil {
					return err
				}
			}
		}
	}

	if !found {
		return fmt.Errorf("could not find any bucket with name: %s", name)
	}

	return nil
}

func runTask(name string, tasks []*config.Task, opts *RunOptions) error {
	var found bool
	for _, task := range tasks {
		if task.Name == name {
			found = true
			if opts.DryRun {
				for _, script := range task.Commands {
					fmt.Println(script)
				}
				continue
			}

			var exec executer.Executer
			if opts.Parallel {
				exec = executer.New("parallel")
			} else {
				exec = executer.New("ssh")
			}

			if err := exec.Run(*task, func(o *executer.CmdOutput) { o.Display() }); err != nil {
				fmt.Fprintf(os.Stderr, "could not run task: %v\n", err)
			}
		}
	}

	if !found {
		return fmt.Errorf("could not find any task with name: %s", name)
	}

	return nil
}

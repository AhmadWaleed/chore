package cli

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var host string
var filename string

func NewInitCommand() *cobra.Command {
	c := &cobra.Command{
		Use:   "init",
		Short: "Create new config file",
		Long:  "Create a new chore config file in the current directory.",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			if err := createFile(); err != nil {
				log.Fatalf("could not initialize config file: %v", err)
			}
		},
	}

	c.Flags().StringVarP(&filename, "name", "", "taskfile", "Config file name default (taskfile).")

	return c
}

func path() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return filepath.Join(dir, filename), nil
}

func createFile() error {
	p, _ := path()

	file, err := os.Create(fmt.Sprintf("%s.yaml", p))
	if err != nil {
		return err
	}

	_, err = file.WriteString(`
servers:
  - localhost@127.0.0.1

tasks:
  - name: deploy
    run: 
      - cd /path/to/site
      - git pull origin main
`)

	if err != nil {
		return err
	}

	return nil
}

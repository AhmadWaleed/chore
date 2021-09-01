package main

import (
	"log"

	"github.com/ahmadwaleed/chore/pkg/cli"
)

func main() {
	cmd := cli.NewCommand()
	if err := cmd.Execute(); err != nil {
		log.Fatalf("could run command: %v", err)
	}
}

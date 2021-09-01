package executer

import (
	"github.com/ahmadwaleed/chore/pkg/config"
)

type OnRunCallback func(*CmdOutput)
type Executer interface {
	Run(config.Task, OnRunCallback) error
}

func New(e string) Executer {
	switch e {
	case "ssh":
		return SSH{}
	case "parallel":
		return Parallel{}
	default:
		return Local{}
	}
}

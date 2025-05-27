package runner

import (
	"github.com/ing-bank/kaet/internal/infratools"
	"github.com/ing-bank/kaet/pkg/types"
)

type Runner struct {
	batch           bool
	Client          *types.ClientConnections
	closed          bool
	interactive     bool
	noRunSafe       bool
	ExecutionMemory infratools.ExecutionMemory // memory stack
}

func (r *Runner) IsInteractive() bool {
	return r.interactive
}

func (r *Runner) IsBatchExecution() bool {
	return r.batch
}

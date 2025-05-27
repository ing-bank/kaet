package types

import (
	"encoding/json"
	"reflect"
	"runtime"

	"github.com/ing-bank/kaet/internal/infratools"
)

const DEFAULT_USER_AGENT = "kaet"

type ExploitFunction func(*ClientConnections, *infratools.MemoryItem, *ExploitOptions) *ExploitResult

func (ef *ExploitFunction) String() string {
	p := reflect.ValueOf(ef).Pointer()
	rf := runtime.FuncForPC(p)
	return rf.Name()
}

type Exploit interface {
	PreExecution(*ClientConnections, *infratools.MemoryItem, *ExploitOptions)
	Execution(*ClientConnections, *infratools.MemoryItem, *ExploitOptions) *ExploitResult
	PostExecution(*ClientConnections, *infratools.MemoryItem, *ExploitOptions, *ExploitResult)
	String() string
}

type ExploitResult struct {
	Data         map[string]any
	NewExecution *infratools.MemoryItem
	Successful   bool
}

type ExploitOptions struct {
	Command string             `json:"command"`
	Pod     *PodExploitOptions `json:"pod_options"`
	TTY     *PodExploitTTY     `json:"tty_options"`
}

func (eo *ExploitOptions) String() string {
	marshalBytes, err := json.Marshal(eo)
	if err != nil {
		return ""
	}
	return string(marshalBytes)
}

type PodExploitOptions struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

type PodExploitTTY struct {
	Command string `json:"command"` // could be bash, ash, zsh, ...
}

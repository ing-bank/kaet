package infratools

import (
	"encoding/json"
	"strings"
)

type ExecutionMemory []*MemoryItem

func (ss *ExecutionMemory) String() string {
	ssItemsStringList := make([]string, 0)

	for _, i := range *ss {
		ssItemsStringList = append(ssItemsStringList, i.String())
	}

	return strings.Join(ssItemsStringList, " ")
}

type MemoryItem struct {
	Kubernetes *MemoryItemKubernetes `json:"kubernetes"`
	Exploit    *StepStackItemExploit `json:"exploit"`
}

func (ssi *MemoryItem) String() string {
	sb := &strings.Builder{}

	if ssi.Kubernetes != nil {
		sb.WriteRune('[')
		sb.WriteString("kubernetes: {" + ssi.Kubernetes.String() + "}")
		sb.WriteRune(']')
	}

	return sb.String()
}

func (ssi *MemoryItem) JSON() string {
	b, err := json.Marshal(ssi)
	if err != nil {
		return "{}"
	}
	return string(b)
}

type MemoryItemKubernetes struct {
	BearerToken          string      `json:"bearer_token,omitempty"`
	CurrentPermissionSet interface{} `json:"-"`
	Namespace            string      `json:"namespace"`
	URL                  string      `json:"url"`
	UserAgent            string      `json:"user_agent"`
	InsecureTLS          bool        `json:"insecure_tls"`
}

type StepStackItemExploit struct {
	// map[namespace]pod_name
	ExploredPods map[string]string `json:"explored_pods"`
}

func (ssik *MemoryItemKubernetes) String() string {
	sb := &strings.Builder{}

	sb.WriteString(ssik.URL)

	sb.WriteString(" [")
	sb.WriteString(ssik.Namespace)
	sb.WriteRune(']')

	sb.WriteString(" [")
	btHead := ssik.BearerToken
	btTail := ssik.BearerToken
	if len(ssik.BearerToken) > 5 {
		btHead = ssik.BearerToken[0:6]
		btTail = ssik.BearerToken[len(ssik.BearerToken)-7:]
	}
	sb.WriteString(btHead + "..." + btTail)
	sb.WriteRune(']')

	return sb.String()
}

func (ss *ExecutionMemory) Push(ssi *MemoryItem) {
	*ss = append(*ss, ssi)
}

func (ss *ExecutionMemory) Pop() *MemoryItem {
	if ss.IsEmpty() {
		return nil
	}

	top := (*ss)[len(*ss)-1]
	*ss = (*ss)[:len(*ss)-1]
	return top
}

func (ss *ExecutionMemory) IsEmpty() bool {
	return len(*ss) == 0
}

func NewExecutionMemory() ExecutionMemory {
	return make(ExecutionMemory, 0)
}

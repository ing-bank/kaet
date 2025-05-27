package runner

import (
	"github.com/ing-bank/kaet/internal/infratools"
	"github.com/ing-bank/kaet/pkg/clouds/kubernetes"
	"github.com/ing-bank/kaet/pkg/types"
)

func newClientConnectionFromOptions(options *types.Options) *types.ClientConnections {
	connections := &types.ClientConnections{}

	if options.Kubernetes != nil {
		connections.Kubernetes = kubernetes.NewClientFromOptions(options)
	}

	return connections
}

func (r *Runner) GetClientConnections() *types.ClientConnections {
	return r.Client
}

func (r *Runner) ResetClientConnections(item *infratools.MemoryItem) {
	if item.Kubernetes != nil {
		r.Client.Kubernetes = kubernetes.NewClientFromOptions(&types.Options{
			Kubernetes: &types.OptionsKubernetes{
				BaseURL:         item.Kubernetes.URL,
				IgnoreTLS:       item.Kubernetes.InsecureTLS,
				Namespace:       item.Kubernetes.Namespace,
				AuthToken:       item.Kubernetes.BearerToken,
				CustomUserAgent: item.Kubernetes.UserAgent,
			},
		})
	}
}

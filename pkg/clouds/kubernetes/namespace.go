package kubernetes

import (
	"context"

	"github.com/ing-bank/kaet/pkg/types"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func ListNamespaces(cc *types.ClientConnections) ([]string, error) {
	nsList := make([]string, 0)

	clusterNamespaces, err := cc.Kubernetes.CoreV1().Namespaces().List(context.TODO(), v1.ListOptions{})
	if err != nil {
		return nsList, err
	}

	for _, ns := range clusterNamespaces.Items {
		nsList = append(nsList, ns.Name)
	}

	return nsList, nil
}

package types

type Runner interface {
	Run()
	Close()
}

type ClientConnections struct {
	Kubernetes *KubernetesClient
	OpenShift  any
	Azure      any
}

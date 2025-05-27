package types

import (
	"net/url"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
)

type KubernetesClient struct {
	BaseURL     string
	BearerToken string
	Namespace   string
	InsecureTLS bool
	UserAgent   string
	options     *OptionsKubernetes
	*kubernetes.Clientset
}

func (kc *KubernetesClient) SetOptions(kOptions *OptionsKubernetes) {
	kc.options = kOptions
}

func (kc *KubernetesClient) NewSPDYExecutor(method string, url *url.URL) (remotecommand.Executor, error) {
	rc := &rest.Config{
		Host:        kc.options.BaseURL,
		BearerToken: kc.options.AuthToken,
		UserAgent:   kc.options.CustomUserAgent,
		TLSClientConfig: rest.TLSClientConfig{
			Insecure: kc.options.IgnoreTLS,
		},
	}

	return remotecommand.NewSPDYExecutor(rc, method, url)
}

func (kc *KubernetesClient) NewWebSocketExecutor(method string, url *url.URL) (remotecommand.Executor, error) {
	rc := &rest.Config{
		Host:        kc.options.BaseURL,
		BearerToken: kc.options.AuthToken,
		UserAgent:   kc.options.CustomUserAgent,
		TLSClientConfig: rest.TLSClientConfig{
			Insecure: kc.options.IgnoreTLS,
		},
	}
	return remotecommand.NewWebSocketExecutor(rc, method, url.String())
}

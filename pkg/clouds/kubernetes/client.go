package kubernetes

import (
	"errors"

	utilz "github.com/ing-bank/kaet/pkg/clouds/kubernetes/utils"
	"github.com/ing-bank/kaet/pkg/types"
	"github.com/projectdiscovery/gologger"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func NewClientFromOptions(options *types.Options) (lc *types.KubernetesClient) {
	lc = &types.KubernetesClient{
		BaseURL:     options.Kubernetes.BaseURL,
		BearerToken: options.Kubernetes.AuthToken,
		Namespace:   options.Kubernetes.Namespace,
		InsecureTLS: options.Kubernetes.IgnoreTLS,
		UserAgent:   options.Kubernetes.CustomUserAgent,
	}
	lc.SetOptions(options.Kubernetes)

	var client *kubernetes.Clientset
	var err error

	// always try first to create the custom client
	client, err = getCustomClient(lc, options)
	if err != nil {
		gologger.Error().
			Str("code_location", "clouds/kubernetes/client").
			Str("error", err.Error()).
			Msg("kubernetes custom client\n")
	}

	if client == nil {
		client, err = getKubeConfigClient(lc, options)
		if err != nil {
			gologger.Error().
				Str("code_location", "clouds/kubernetes/client").
				Str("error", err.Error()).
				Msg("kubernetes kubeconfig client\n")
		}
	}

	if client == nil {
		client, err = getInPodClient(lc, options)
		if err != nil {
			gologger.Error().
				Str("code_location", "clouds/kubernetes/client").
				Str("error", err.Error()).
				Msg("kubernetes in pod client\n")
		}
	}

	if client == nil {
		panic("invalid kubernetes options. could not create client")
	}

	lc.Clientset = client

	return
}

func Ping(cc *types.ClientConnections) bool {
	openApiSchema, err := cc.Kubernetes.OpenAPISchema()
	if err != nil {
		gologger.Error().
			Str("code_location", "clouds/kubernetes/client").
			Str("error", err.Error()).
			Msg("could not connect to kubernetes API\n")
		return false
	}

	return openApiSchema.Info != nil
}

func getCustomClient(lc *types.KubernetesClient, o *types.Options) (*kubernetes.Clientset, error) {
	if o.Kubernetes == nil {
		return nil, errors.New("invalid kubernetes custom client options")
	}

	if o.Kubernetes.BaseURL == "" || o.Kubernetes.AuthToken == "" {
		return nil, errors.New("invalid configurations for kubernetes custom client")
	}

	config := &rest.Config{
		BearerToken: o.Kubernetes.AuthToken,
		Host:        o.Kubernetes.BaseURL,
		TLSClientConfig: rest.TLSClientConfig{
			Insecure: o.Kubernetes.IgnoreTLS,
		},
		UserAgent: o.Kubernetes.CustomUserAgent,
	}
	setConfigCustomOptions(config, o)
	httpClient, err := rest.HTTPClientFor(config)
	if err != nil {
		return nil, err
	}

	client, err := kubernetes.NewForConfigAndClient(config, httpClient)
	if err != nil {
		return nil, err
	}

	lc.UserAgent = config.UserAgent

	return client, nil
}

func getKubeConfigClient(lc *types.KubernetesClient, o *types.Options) (*kubernetes.Clientset, error) {
	if o.Kubernetes == nil {
		return nil, errors.New("invalid kubernetes kubeconfig options")
	}

	if o.Kubernetes.KubeConfigPath == "" {
		return nil, errors.New("invalid configuration for kubernetes kubeconfig client")
	}

	if o.Kubernetes.AuthToken != "" {
		gologger.Warning().
			Str("code_location", "clouds/kubernetes/client").
			Msg("-token will be ignored. using kubeconfig configuration")
	}

	config, err := clientcmd.BuildConfigFromFlags(o.Kubernetes.BaseURL, o.Kubernetes.KubeConfigPath)
	if err != nil {
		return nil, err
	}
	setConfigCustomOptions(config, o)

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	lc.BearerToken = config.BearerToken
	lc.Namespace = utilz.GetNamespaceFromSAToken(lc.BearerToken)
	lc.InsecureTLS = config.Insecure
	lc.UserAgent = config.UserAgent

	return client, nil
}

func getInPodClient(lc *types.KubernetesClient, o *types.Options) (*kubernetes.Clientset, error) {
	// it falls back to in pod config
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	setConfigCustomOptions(config, o)

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	lc.BaseURL = config.Host
	lc.BearerToken = config.BearerToken
	lc.Namespace = utilz.GetNamespaceFromSAToken(lc.BearerToken)
	lc.InsecureTLS = config.Insecure
	lc.UserAgent = config.UserAgent

	return client, nil
}

func setConfigCustomOptions(cfg *rest.Config, o *types.Options) {
	if cfg == nil || o == nil || o.Kubernetes == nil {
		return
	}

	if o.Kubernetes.CustomUserAgent != "" {
		cfg.UserAgent = o.Kubernetes.CustomUserAgent
	} else {
		cfg.UserAgent = types.DEFAULT_USER_AGENT
	}

	if o.Kubernetes.IgnoreTLS {
		cfg.TLSClientConfig = rest.TLSClientConfig{
			Insecure: o.Kubernetes.IgnoreTLS,
		}
	}

	cfg.QPS = 400
	cfg.Burst = 400
}

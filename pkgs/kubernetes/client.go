package kubernetes

import (
	"github.com/rs/zerolog/log"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/flowcontrol"
)

type Client struct {
	ClientSet *kubernetes.Clientset
	Rest      *rest.Config
}

func NewClient(opts *Options) (*Client, error) {
	var cfg *rest.Config
	var err error

	if opts.Server != "" {
		// always true, otherwise certificate authentication will throw an error
		opts.InsecureTLS = true
		cfg = &rest.Config{
			Host:        opts.Server,
			BearerToken: opts.Token,
			TLSClientConfig: rest.TLSClientConfig{
				Insecure: opts.InsecureTLS,
			},
		}
	} else if opts.Token != "" && opts.Kubeconfig != "" {
		cfg, err = clientcmd.BuildConfigFromFlags("", opts.Kubeconfig)

		cfg = &rest.Config{
			Host:        cfg.Host,
			BearerToken: opts.Token,
			UserAgent:   opts.UserAgent,
		}
		opts.InsecureTLS = true
	} else if opts.Kubeconfig != "" {
		cfg, err = clientcmd.BuildConfigFromFlags("", opts.Kubeconfig)
	} else {
		cfg, err = rest.InClusterConfig()
	}
	configureKubernetesClient(opts, cfg)

	if err != nil {
		log.Error().Msgf("failed to get valid kubernetes client. error: %s\n", err)
		return nil, err
	}

	configureKubernetesClient(opts, cfg)
	cs, err := kubernetes.NewForConfig(cfg)
	return &Client{cs, cfg}, err
}

func ClientFromToken(opts *Options, base *rest.Config, token string) (*Client, error) {
	// TODO: maybe change the `base` to a server string
	cfg := rest.CopyConfig(base)
	cfg.BearerToken = token

	configureKubernetesClient(opts, cfg)
	cs, err := kubernetes.NewForConfig(cfg)
	return &Client{cs, cfg}, err
}

func configureKubernetesClient(opts *Options, cfg *rest.Config) {
	if cfg == nil {
		return
	}

	if opts.NoRateLimit {
		cfg.RateLimiter = flowcontrol.NewFakeAlwaysRateLimiter()
	}
	if opts.UserAgent != "" {
		cfg.UserAgent = opts.UserAgent
	}

	if opts.InsecureTLS {
		cfg.TLSClientConfig.Insecure = opts.InsecureTLS
		cfg.CAData = nil
		cfg.CAFile = ""
	}
}

func (c *Client) Server() string {
	return c.Rest.Host
}

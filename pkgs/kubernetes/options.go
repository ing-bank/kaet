package kubernetes

type Options struct {
	InsecureTLS bool
	Kubeconfig  string
	NoRateLimit bool
	Server      string
	Token       string
	UserAgent   string
}

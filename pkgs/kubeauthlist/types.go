package kubeauthlist

import (
	v1 "k8s.io/api/authorization/v1"
)

type Options struct {
	AuthToken   string
	InsecureTLS bool
	Kubeconfig  string
	Namespace   string
	NoRateLimit bool
	Server      string
	UserAgent   string
}

type ProcessingPropertiesCtxType struct{}

type ProcessingProperties struct {
	ResultChan chan *v1.SelfSubjectAccessReview
	OutputChan chan *v1.SelfSubjectAccessReview
}

type Result struct {
	Token          string
	Namespace      string
	PermissionList []*v1.SelfSubjectAccessReview
}

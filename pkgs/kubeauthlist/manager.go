package kubeauthlist

import (
	"context"
	"sync"

	"github.com/ing-bank/kaet/pkgs/kubernetes"
	"github.com/rs/zerolog/log"
	authv1 "k8s.io/api/authorization/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type AuthList struct {
	KubeClient *kubernetes.Client
	Options    *Options
}

func New(opts *Options) (*AuthList, error) {
	kubeClient, err := kubernetes.NewClient(&kubernetes.Options{
		InsecureTLS: opts.InsecureTLS,
		Kubeconfig:  opts.Kubeconfig,
		NoRateLimit: opts.NoRateLimit,
		Server:      opts.Server,
		Token:       opts.AuthToken,
		UserAgent:   opts.UserAgent,
	})
	if err != nil {
		log.Error().Err(err).Msg("failed to create kubernetes client")
		return nil, err
	}

	return &AuthList{
		KubeClient: kubeClient,
		Options:    opts,
	}, nil
}

func NewFromKubeClient(kc *kubernetes.Client, ns string) (*AuthList, error) {
	return &AuthList{
		KubeClient: kc,
		Options:    &Options{Namespace: ns},
	}, nil
}

func GetDefaultCtx() context.Context {
	processingChannels := ProcessingProperties{
		ResultChan: make(chan *authv1.SelfSubjectAccessReview, 100),
		OutputChan: nil,
	}

	return context.WithValue(context.TODO(), ProcessingPropertiesCtxType{}, processingChannels)
}

func (al *AuthList) Start(ctx context.Context) (*Result, error) {
	groupList, resourceList, err := al.KubeClient.GetServerResources()
	if err != nil {
		log.Error().Err(err).Msg("failed to get groups and resources")
		return nil, nil
	}

	accessResultList := make([]*authv1.SelfSubjectAccessReview, 0)
	processingProperties := ctx.Value(ProcessingPropertiesCtxType{}).(ProcessingProperties)

	var resultWg sync.WaitGroup
	if processingProperties.ResultChan != nil {
		resultWg.Add(1)
		go func(c chan *authv1.SelfSubjectAccessReview, r *[]*authv1.SelfSubjectAccessReview) {
			defer resultWg.Done()

			for accessResult := range c {
				*r = append(*r, accessResult)
			}

		}(processingProperties.ResultChan, &accessResultList)
	}

	for _, resource := range resourceList {
		group := getResourceGroup(resource, groupList)
		al.Analyze(ctx, al.Options.Namespace, group, resource)
	}

	close(processingProperties.ResultChan)
	resultWg.Wait()
	return &Result{
		Token:          al.KubeClient.Rest.BearerToken,
		Namespace:      al.Options.Namespace,
		PermissionList: accessResultList,
	}, nil
}

func (al *AuthList) Finish(ctx context.Context) {
	processingChannels := ctx.Value(ProcessingPropertiesCtxType{}).(ProcessingProperties)

	if processingChannels.OutputChan != nil {
		close(processingChannels.OutputChan)
	}
}

func getResourceGroup(r *v1.APIResourceList, groupList []*v1.APIGroup) *v1.APIGroup {
	for _, group := range groupList {
		for _, gForDiscovery := range group.Versions {
			if r.GroupVersion == gForDiscovery.GroupVersion {
				return group
			}
		}
	}
	return nil
}

package kubeauthlist

import (
	"context"

	"github.com/ing-bank/kaet/pkgs/kubernetes"
	"github.com/rs/zerolog/log"
	errors2 "k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (al *AuthList) Analyze(ctx context.Context, ns string, group *v1.APIGroup,
	apiResource *v1.APIResourceList) {

	pp := ctx.Value(ProcessingPropertiesCtxType{}).(ProcessingProperties)

	for _, httpVerb := range kubernetes.ApiVerbs {
		for _, resource := range apiResource.APIResources {
			accessResponse, err := al.KubeClient.SelfSubjectAccessReview(ctx, ns, httpVerb, group,
				resource)
			errors2.ReasonForError(err)
			if err != nil {
				log.Debug().Err(err).
					Str("group_name", group.Name).
					Str("group_version", group.PreferredVersion.Version).
					Str("resource_name", resource.Name).
					Msg("failed to retrieve SelfSubjectAccessReview")
				continue
			}

			if pp.ResultChan != nil {
				pp.ResultChan <- accessResponse
			}
			if pp.OutputChan != nil {
				pp.OutputChan <- accessResponse
			}
		}
	}
}

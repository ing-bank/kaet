package kubernetes

import (
	"context"
	"strconv"

	"github.com/ing-bank/kaet/pkg/types"
	"github.com/projectdiscovery/gologger"
	v1 "k8s.io/api/authorization/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func AuthCanI(cc *types.ClientConnections, verb, resource, ns string) bool {
	canI := false
	defer func(v, r, n string, ci *bool) {
		gologger.Debug().
			Str("code_location", "clouds/kubernetes/auth").
			Str("api_verb", v).
			Str("resource", r).
			Str("namespace", n).
			Str("can_i", strconv.FormatBool(*ci)).
			Msg("auth can-i result")
	}(verb, resource, ns, &canI)

	accessReview, err := cc.Kubernetes.AuthorizationV1().SelfSubjectAccessReviews().
		Create(
			context.TODO(),
			&v1.SelfSubjectAccessReview{
				TypeMeta: metaV1.TypeMeta{
					Kind:       "SelfSubjectAccessReview",
					APIVersion: "authorization.k8s.io/v1",
				},
				Spec: v1.SelfSubjectAccessReviewSpec{
					ResourceAttributes: &v1.ResourceAttributes{
						Verb:      verb,
						Resource:  resource,
						Namespace: ns,
					},
				},
				ObjectMeta: metaV1.ObjectMeta{
					Namespace: ns,
				},
			},
			metaV1.CreateOptions{},
		)
	if err != nil {
		gologger.Error().
			Str("code_location", "clouds/kubernetes/auth").
			Str("api_verb", verb).
			Str("resource", resource).
			Str("namespace", ns).
			Str("error", err.Error()).
			Msg("could not get SelfSubjectAccessReview\n")
		return canI
	}

	canI = accessReview.Status.Allowed

	return canI
}

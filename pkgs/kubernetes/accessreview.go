package kubernetes

import (
	"context"
	"strings"

	authv1 "k8s.io/api/authorization/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (c *Client) SelfSubjectAccessReview(ctx context.Context, ns, verb string, group *v1.APIGroup,
	resource v1.APIResource) (*authv1.SelfSubjectAccessReview, error) {
	resourceNameParts := strings.Split(resource.Name, "/")

	resourceName := resourceNameParts[0]
	subResourceName := ""
	if len(resourceNameParts) >= 2 {
		subResourceName = resourceNameParts[1]
	}

	accessReviewRequest := &authv1.SelfSubjectAccessReview{
		Spec: authv1.SelfSubjectAccessReviewSpec{
			ResourceAttributes: &authv1.ResourceAttributes{
				Namespace:     ns,
				Verb:          verb,
				Group:         group.Name,
				Version:       group.PreferredVersion.Version,
				Resource:      resourceName,
				Subresource:   subResourceName,
				Name:          resourceName,
				FieldSelector: nil,
				LabelSelector: nil,
			},
		},
	}

	return c.ClientSet.
		AuthorizationV1().
		SelfSubjectAccessReviews().
		Create(
			ctx,
			accessReviewRequest,
			v1.CreateOptions{},
		)

}

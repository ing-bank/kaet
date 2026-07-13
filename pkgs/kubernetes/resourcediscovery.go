package kubernetes

import (
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (c *Client) GetServerResources() ([]*v1.APIGroup, []*v1.APIResourceList, error) {
	return c.ClientSet.DiscoveryClient.ServerGroupsAndResources()
}

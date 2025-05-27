package discovery

import (
	"fmt"
	"strings"

	"github.com/ing-bank/kaet/internal/infratools"
	"github.com/ing-bank/kaet/pkg/clouds/kubernetes"
	"github.com/ing-bank/kaet/pkg/types"
	kalRunner "github.com/ing-bank/kal/pkg/runner"
	kalTypes "github.com/ing-bank/kal/pkg/types"
	"github.com/projectdiscovery/gologger"
	"golang.org/x/exp/maps"
)

func GetMapOfResourcesPermissions(
	cc *types.ClientConnections,
	state *infratools.MemoryItem,
	namespace string,
) map[*kubernetes.Resource][]kubernetes.APIVerb {
	gologger.Debug().
		Str("code_location", "kaet/discovery/resource").
		Str("namespace", namespace).
		Str("state", state.JSON()).
		Msg("starting enumeration of resource permissions")

	res := make(map[*kubernetes.Resource][]kubernetes.APIVerb, 0)

	authorizationRunner := kalRunner.FromOptions(&kalTypes.Options{
		Kubernetes: &kalTypes.KubernetesOptions{
			ApiToken:    state.Kubernetes.BearerToken,
			Namespace:   namespace,
			ServerURL:   state.Kubernetes.URL,
			InsecureTLS: state.Kubernetes.InsecureTLS,
			NoRateLimit: true,
		},
		Output: &kalTypes.OutputOptions{
			JSON:    types.IsJSONOutput(),
			NoColor: types.IsNoColorOutput(),
		},
		Silent:  types.IsSilentOutput(),
		Verbose: false,
		NoLogs:  false,
	})
	authorizationListResult := authorizationRunner.Exec()
	gologger.Debug().
		Str("resource_permissions_raw", fmt.Sprintf("%#v", authorizationListResult)).
		Msg("finished authorization enumeration\n")

	for composedResourceName, permissions := range authorizationListResult {
		r := &kubernetes.Resource{}

		splitResourceName := strings.Split(composedResourceName, "/")

		r.Name = splitResourceName[0]

		// composed resource name has a version
		if len(splitResourceName) > 1 {
			r.GroupVersion = splitResourceName[1]
		}

		// composed resource name has a sub-resource
		if len(splitResourceName) > 2 {
			r.SubResource = splitResourceName[2]
		}

		res[r] = make([]kubernetes.APIVerb, 0)
		for _, permission := range permissions {
			res[r] = append(res[r], kubernetes.APIVerb(permission))
		}
	}

	gologger.Debug().
		Str("code_location", "kaet/discovery/resource").
		Str("namespace", namespace).
		Str("state", state.JSON()).
		Str("resource_permissions", fmt.Sprintf("%v", maps.Values(res))).
		Msg("finished enumeration of resource permissions")

	return res
}

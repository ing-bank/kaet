package kubernetes

const INTERNAL_KUBERNETES_API_BASE_URL = "https://kubernetes.default.svc"

var CONTROL_NAMESPACES []string = appendMultipleSlices(
	VANILLA_CONTROL_NAMESPACES,
	AKS_CONTROL_NAMESPACES,
	GKE_CONTROL_NAMESPACES,
	AWS_CONTROL_NAMESPACES,
	OPENSHIFT_CONTROL_NAMESPACES,
)

var VANILLA_CONTROL_NAMESPACES []string = []string{
	"kube-node-lease",
	"kube-public",
	"kube-system",
	"local-path-storage",
}
var AKS_CONTROL_NAMESPACES []string = []string{
	"app-routing-system",
	"gatekeeper-system",
}
var GKE_CONTROL_NAMESPACES []string = []string{}
var AWS_CONTROL_NAMESPACES []string = []string{}
var OPENSHIFT_CONTROL_NAMESPACES []string = []string{}

func appendMultipleSlices(s ...[]string) []string {
	res := make([]string, 0)

	for _, item := range s {
		res = append(res, item...)
	}

	return res
}

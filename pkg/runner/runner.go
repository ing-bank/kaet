package runner

import (
	"fmt"
	"slices"
	"strconv"

	"github.com/ing-bank/kaet/internal/infratools"
	"github.com/ing-bank/kaet/pkg/clouds/kubernetes"
	kubeExploits "github.com/ing-bank/kaet/pkg/clouds/kubernetes/exploits"
	k8sUtils "github.com/ing-bank/kaet/pkg/clouds/kubernetes/utils"
	"github.com/ing-bank/kaet/pkg/kaet/discovery"
	"github.com/ing-bank/kaet/pkg/types"
	"github.com/projectdiscovery/gologger"
)

func FromOptions(options *types.Options) types.Runner {
	runner := &Runner{
		Client:          newClientConnectionFromOptions(options),
		ExecutionMemory: infratools.NewExecutionMemory(),
		noRunSafe:       options.Kubernetes.NoRunSafe,
	}

	return runner
}

func (r *Runner) Close() {
	if r.closed {
		return
	}
	r.closed = true
}

func (r *Runner) Run() {
	initialMemoryItem := &infratools.MemoryItem{
		Kubernetes: &infratools.MemoryItemKubernetes{
			BearerToken: r.Client.Kubernetes.BearerToken,
			Namespace:   r.Client.Kubernetes.Namespace,
			URL:         r.Client.Kubernetes.BaseURL,
			InsecureTLS: r.Client.Kubernetes.InsecureTLS,
			UserAgent:   r.Client.Kubernetes.UserAgent,
		},
	}

	if initialMemoryItem.Kubernetes.Namespace == "" {
		initialMemoryItem.Kubernetes.Namespace = k8sUtils.GetNamespaceFromSAToken(r.Client.Kubernetes.BearerToken)
	}

	r.ExecutionMemory.Push(initialMemoryItem)

	r.recursiveRun()
}

func (r *Runner) recursiveRun() {
	if r.ExecutionMemory.IsEmpty() {
		return
	}

	currentState := r.ExecutionMemory.Pop()
	r.ResetClientConnections(currentState)
	clientConnections := r.GetClientConnections()
	gologger.Debug().Str("state", currentState.JSON()).Msg("starting recursive run\n")

	gologger.Info().Msg("starting exploration\n")

	canIConnectToKubernetes := kubernetes.Ping(clientConnections)
	if !canIConnectToKubernetes {
		gologger.
			Fatal().
			Str("state", currentState.JSON()).
			Msg("could not connect to kubernetes with the provided configuration\n")
	}

	canIListNamespaces := kubernetes.AuthCanI(
		clientConnections,
		string(kubernetes.API_VERB_LIST),
		"namespaces",
		"",
	)
	gologger.Debug().Str("result", strconv.FormatBool(canIListNamespaces)).Msg("can I list all namespaces?\n")

	namespaceList := make([]string, 0)
	if canIListNamespaces {
		gologger.Debug().Msg("listing all namespaces\n")

		namespaceListResponse, err := kubernetes.ListNamespaces(clientConnections)
		if err != nil {
			gologger.Error().Str("error", err.Error()).Msg("could not list all namespaces\n")
		} else {
			namespaceList = namespaceListResponse
		}

	}

	// always add the token namespace to the namespace list, if it doesn't exist
	if !slices.Contains(namespaceList, currentState.Kubernetes.Namespace) {
		gologger.Debug().Msg("adding service account token namespace to the exploration list\n")
		namespaceList = append(namespaceList, currentState.Kubernetes.Namespace)
	}

	// always safe, remove control namespaces
	if !r.noRunSafe {
		// remove control namespaces so it won't crash the cluster
		namespaceList = slices.DeleteFunc(namespaceList, func(ns string) bool {
			return slices.Contains(kubernetes.CONTROL_NAMESPACES, ns)
		})
	}

	gologger.Info().Str("namespace_quantity", strconv.Itoa(len(namespaceList))).Msg("found namespace(s) to explore\n")
	gologger.Debug().Str("namespace_list", fmt.Sprintf("%v", namespaceList)).Msg("namespaces found\n")

	for _, namespace := range namespaceList {
		r.exploreNamespace(namespace, currentState)
	}
}

func (r *Runner) exploreNamespace(namespace string, state *infratools.MemoryItem) {
	gologger.Info().Str("namespace", namespace).Msg("starting namespace exploration\n")

	resourcePermissionsMap := discovery.GetMapOfResourcesPermissions(r.GetClientConnections(), state, namespace)

	for resource, permissionList := range resourcePermissionsMap {
		gologger.Debug().
			Str("resource", resource.String()).
			Str("permission_list", fmt.Sprintf("%v", permissionList)).
			Msg("exploiting resource\n")

		r.exploitResource(state, resource, permissionList, namespace)
	}

	gologger.Info().Str("namespace", namespace).Msg("finished namespace exploration\n")
}

func (r *Runner) exploitResource(
	state *infratools.MemoryItem,
	resource *kubernetes.Resource,
	resourceVerbs []kubernetes.APIVerb,
	namespace string,
) {
	gologger.Debug().
		Str("resource", resource.String()).
		Msg("starting resource exploration\n")

	possibleExploits := kubeExploits.GetPossibleExploits(resource, resourceVerbs)

	if len(possibleExploits) == 0 {
		gologger.Info().
			Str("resource", resource.String()).
			Msg(types.AU.Green("no valid exploits found").String() + "\n")
		return
	} else {
		gologger.Info().
			Str("resource", resource.String()).
			Str("exploit_quantity", strconv.Itoa(len(possibleExploits))).
			Msg("found possible exploitation path(s)\n")
	}

	for _, exploit := range possibleExploits {
		gologger.Info().
			Str("resource", resource.String()).
			Str("exploit_name", exploit.String()).
			Msg(types.AU.Red("exploiting resource").String() + "\n")

		exploitOptions := &types.ExploitOptions{
			Pod: &types.PodExploitOptions{
				Namespace: namespace,
			},
		}

		clientConnections := r.GetClientConnections()

		exploit.PreExecution(clientConnections, state, exploitOptions)
		resourceExploitResult := exploit.Execution(clientConnections, state, exploitOptions)
		exploit.PostExecution(clientConnections, state, exploitOptions, resourceExploitResult)

		gologger.Debug().Str("result", fmt.Sprintf("%#v", resourceExploitResult)).Msg("exploit result\n")

		// if after all the executions the result has a `NewStep`, it will:
		// 1. add the current step to the memory stack
		// 2. add the new step in the memory stack
		// 3. start a recursive execution
		if resourceExploitResult.NewExecution != nil {
			// save step permissions and config
			gologger.Debug().
				Str("code_location", "runner/runner").
				Str("state", state.JSON()).
				Msg("adding current state to memory\n")
			r.ExecutionMemory.Push(state)
			gologger.Debug().
				Str("code_location", "runner/runner").
				Str("state", resourceExploitResult.NewExecution.JSON()).
				Msg("adding next state to memory\n")
			r.ExecutionMemory.Push(resourceExploitResult.NewExecution)
			r.recursiveRun()

			state = r.ExecutionMemory.Pop()
			gologger.Debug().
				Str("code_location", "runner/runner").
				Str("state", state.JSON()).
				Msg("removing previous state from memory\n")
			r.ResetClientConnections(state)
		}
	}

	gologger.Info().Str("resource", resource.String()).
		Msg("finished resource exploration\n")
}

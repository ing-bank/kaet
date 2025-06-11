# Using KAET

This walkthrough details how to use KAET locally, in a [Docker](https://docs.docker.com/) container, applying manually to a [Kubernetes Cluster](https://kubernetes.io/), or using [Helm Charts](https://helm.sh/).

## Local

```bash linenums="1"
> kaet -h

#######################################
#                                     #
#  ██╗  ██╗ █████╗ ███████╗████████╗  #
#  ██║ ██╔╝██╔══██╗██╔════╝╚══██╔══╝  #
#  █████╔╝ ███████║█████╗     ██║     #
#  ██╔═██╗ ██╔══██║██╔══╝     ██║     #
#  ██║  ██╗██║  ██║███████╗   ██║     #
#  ╚═╝  ╚═╝╚═╝  ╚═╝╚══════╝   ╚═╝     #
#    Kubernetes Auto Exploit Tool     #
#######################################


Usage:
  kaet [flags]

Flags:
KUBERNETES OPTIONS:
   -k8s-url string                    kubernetes API base url (default "https://kubernetes.default.svc")
   -serviceaccounttoken, -sat string  kubernetes service account token
   -k, -ignore-tls                    ignore TLS
   -ua, -user-agent string            custom user agent (default "KAET")
   -n, -namespace string              kubernetes namespace
   -safe                              do not explore control namespaces
   -kubeconfig string                 absolute path to kubeconfig file (default "/home/kaet/.kube/config")

EXECUTION OPTIONS:
   -batch             accept all default responses
   -it, -interactive  interactive execution

OUTPUT OPTIONS:
   -v, -verbose    verbose output
   -s, -silent     silent output
   -j, -json       json output
   -nc, -no-color  colorful output
```

## Helm Charts

Using Helm charts, you can deploy KAET with a single command. The following command pulls a Helm chart from Github Package Registry and publishes it to a Kubernetes cluster.

```bash linenums="1"
helm upgrade kaet oci://ghcr.io/ing-bank/kaet-helm:0.1.0 --version 0.1.0 -n your_namespace --set 'serviceAccounts={the-service-account-name-you-want-to-test}'
```

After running this command, Helm will create a Kubernetes Job in the cluster that will execute KAET. For each Service Account Name provided, a Job object will be created.

### Finding Created Jobs

- Jobs

```bash linenums="1"
kubectl get jobs -n your_namespace

NAME       COMPLETIONS   DURATION   AGE
kaet-job   1/1           5m33s      8m39s
```

- Pods

```bash linenums="1"
kubectl get pods -n your_namespace

NAME                                       READY   STATUS      RESTARTS   AGE
kaet-job-5cdmk                             0/1     OOMKilled   0          2m57s
kaet-job-wpr4n                             1/1     Running     0          103s
```

## Docker

```bash linenums="1"
docker run --rm ghcr.io/ing-bank/kaet:latest -h

#######################################
#                                     #
#  ██╗  ██╗ █████╗ ███████╗████████╗  #
#  ██║ ██╔╝██╔══██╗██╔════╝╚══██╔══╝  #
#  █████╔╝ ███████║█████╗     ██║     #
#  ██╔═██╗ ██╔══██║██╔══╝     ██║     #
#  ██║  ██╗██║  ██║███████╗   ██║     #
#  ╚═╝  ╚═╝╚═╝  ╚═╝╚══════╝   ╚═╝     #
#    Kubernetes Auto Exploit Tool     #
#######################################


Usage:
  kaet [flags]

Flags:
KUBERNETES OPTIONS:
   -k8s-url string                    kubernetes API base url (default "https://kubernetes.default.svc")
   -serviceaccounttoken, -sat string  kubernetes service account token
   -k, -ignore-tls                    ignore TLS
   -ua, -user-agent string            custom user agent (default "KAET")
   -n, -namespace string              kubernetes namespace
   -safe                              do not explore control namespaces
   -kubeconfig string                 absolute path to kubeconfig file (default "/home/kaet/.kube/config")

EXECUTION OPTIONS:
   -batch             accept all default responses
   -it, -interactive  interactive execution

OUTPUT OPTIONS:
   -v, -verbose    verbose output
   -s, -silent     silent output
   -j, -json       json output
   -nc, -no-color  colorful output
```

## Manually Creating a Kubernetes Job

```bash linenums="1"
kubectl apply -f - <<EOF
apiVersion: batch/v1
kind: Job
metadata:
  name: kaet
  namespace: your-namespace
spec:
  selector: {}
  backoffLimit: 3
  template:
    metadata:
      name: kaet-job
    spec:
      serviceAccount: the-service-account-name-you-want-to-test
      restartPolicy: Never
      containers:
        - name: kaet
          image: ghcr.io/ing-bank/kaet:latest
          args: [""]
          resources:
            limits:
              memory: '100Mi'
              cpu: '100m'
EOF
```

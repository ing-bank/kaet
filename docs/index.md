<img width=273 src="./imgs/logo.jpeg" alt="KAET logo" align=left />

# Kubernetes Auto Exploit Tool (KAET)

`KAET`: an automation that analyzes weaknesses in Role-Based Access Controls (RBAC) in Kubernetes Clusters. This tool uses a set of known attacks on misconfigurations and loose permissions in RBAC controls, finding attack paths based on initial access to the cluster.

Kubernetes Clusters have a large number of [Roles](https://kubernetes.io/docs/reference/access-authn-authz/rbac/#role-example) and [Cluster Roles](https://kubernetes.io/docs/reference/access-authn-authz/rbac/#clusterrole-example), making it not feasible for humans to test all possible combinations and verify what a malicious actor can do with those permissions. Therefore, we need an automation to perform this evaluation and provide feedback. `KAET` can do it all! In this case, `KAET` actively tests all possible attacks, based on initial access inside or outside the cluster.

In that case, based on the initial access, `KAET` enumerates all current permissions using [KAL](https://github.com/ing-bank/kal). Each permission rule uses loose permissions and misconfigurations to exploit the Kubernetes Cluster and its workloads.

Examples of loose permissions are:

- Reading Secrets of other namespaces
- A Service Account Token being able to start a new Kubernetes POD from inside the cluster
- Executing remote commands in other PODs, using the `pods/exec` resource

## Main Features

- No additional role required to run KAET
- Self-contained Kubernetes exploitation tool
- Evaluate Role-Based Access Control (RBAC)
- Executable in Zero-privileged environments
- Non-interactive usage
- Multiple options for customized execution

## Contributing

Contributions are more than welcome! Please see our [contribution guidelines first](https://github.com/ing-bank/kaet/blob/develop/CONTRIBUTING.md).

## License

You can check our licensing scheme [here](https://github.com/ing-bank/kaet/blob/develop/LICENSE).

## Tools that Inspired KAET

- [BOtB](https://github.com/brompwnie/botb)
- [Kubedestroyer](https://github.com/Rolix44/Kubestroyer)
- [KubeHound](https://kubehound.io/)
- [KubiScan](https://github.com/cyberark/KubiScan)
- [peirates](https://github.com/inguardians/peirates)
- [rbac-police](https://github.com/PaloAltoNetworks/rbac-police)
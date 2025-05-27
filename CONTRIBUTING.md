# Contributing to KAET

`KAET` is an automation used to analyze weaknesses on Role-Based Access Controls (RBAC) in Kubernetes Clusters. This tool uses a set of known attacks of misconfigurations and loose permissions in RBAC controls finding attack paths based on initial access to the cluster.

Kubernetes Clusters have a great number of [Roles](https://kubernetes.io/docs/reference/access-authn-authz/rbac/#role-example) and [Cluster Roles](https://kubernetes.io/docs/reference/access-authn-authz/rbac/#clusterrole-example) making it not feasible for humans to test all possible combinations and verify what a malicious actor is able to do with those permissions. Therefore, we need an automation to perform that analysis and provide feedback. `KAET` can do it all! In this case, `KAET` actively tests all possible attacks, based on an initial access inside the cluster, or outside the cluster.

In that case, based on the initial access, `KAET` enumerates all current permissions using [KAL](https://github.com/ing-bank/kal) and for each permission rule it tries to exploit misconfiguration and loose permissions.

We're very much open to contributions but there are some things to keep in mind:

- Discuss the feature and implementation you want to add on Github before you write a PR for it. On disagreements, maintainer(s) will have the final word.
- Features need a somewhat general use case. If the use case is very niche it will be hard for us to consider maintaining it.
- If you’re going to add a feature, consider if you could help out in the maintenance of it.
- When issues or pull requests are not going to be resolved or merged, they should be closed as soon as possible. This is kinder than deciding this after a long period. Our issue tracker should reflect work to be done.

That said, there are many ways to contribute to KAET, including:

- Contribution to code
- Improving the documentation
- Reviewing merge requests
- Investigating bugs
- Reporting issues

<!-- Starting out with open source? See the guide [How to Contribute to Open Source](https://opensource.guide/how-to-contribute/) and have a look at [our issues labelled *good first issue*](https://github.com/ing-bank/probatus/issues?q=is%3Aissue+is%3Aopen+label%3A%22good+first+issue%22). -->


Planned improvements and changes are discussed in GitHub issues. Feel free to open a discussion.

Roadmap and future features are managed through GitHub projects.

## Standards

- Go latest version
- Formatter: [gofumpt](https://github.com/mvdan/gofumpt)
- Linter: [golangci-lint](https://github.com/golangci/golangci-lint)

### Branch structure

`KAET` follows the [GitFlow Workflow](https://danielkummer.github.io/git-flow-cheatsheet/). 

- **master**: release code
- **develop**: on going feature development
- **feature/<your_feature>**: create based on **develop** branch, contains new features and improvements
- **bugfix/<your_bugfix>**: create based on **develop** branch, contains bug fixes

## Feature requests

If you have an improvement idea on how we can improve KAET, please check our [discussions](https://github.com/ing-bank/kaet/issues) and [projects](https://github.com/ing-bank/kaet/projects?query=is%3Aopen) to see if there are similar ideas or feature requests. If there is no similar idea, feel free to create a [new discussion](https://github.com/ing-bank/kaet/issues/new) with your idea. Use the **Feature request** template.

## Pull requests

Help out improving and fixing bugs in KAET by sending a Pull Request. Make sure that before contributing you create a fork of KAET and develop in your source version. After completing the development, open a [Pull Request](https://github.com/ing-bank/kaet/pulls) following the **Pull Request** issue template.

Make sure no tests are broken and your code is linted:

```sh
make test
make lint
```

### Increase you PR approval

- Write tests
- Add documentation
- Write a [good commit message](https://www.conventionalcommits.org/).
<!--
Please create an issue to collect feedback prior to feature additions. Please also reference that issue in any PRs.
If possible try to keep PRs scoped to one feature, and add tests for new features.
-->

### Description:
Explain the purpose of the PR. Add an explanation of the use case or bug that it fix.

### Checklist:
* [ ] Tests passing (`go test -timeout=5m ./...`)?
* [ ] Lint passing (`golangci-lint run -c .golangci.yaml`)? This requires [golangci-lint](https://golangci-lint.run/welcome/install/#local-installation)
* [ ] Basic security test passing (`semgrep scan . --exclude=readme.md`)? This requires [semgrep](https://semgrep.dev/docs/getting-started/quickstart)

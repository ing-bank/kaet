# Installing KAET

## Requirements

- [Go](https://go.dev/doc/install) >= 1.20.0
- [Git](https://git-scm.com/downloads)
- [Docker](https://docs.docker.com/engine/install/) (*optional*)
- [Azure CLI](https://learn.microsoft.com/en-us/cli/azure/install-azure-cli) (*optional*)
- [Helm Chart](https://helm.sh/docs/intro/install/) (*optional*)

## From Source

You can install KAET by cloning it's source code.

```bash linenums="1"
git clone --depth 1 https://github.com/ing-bank/kaet
cd kaet
go install
```

## Docker

This setup will pull KAET's latest Docker image and will make it available to use.

```bash linenums="1"
docker pull ghcr.io/ing-bank/kaet:latest
```


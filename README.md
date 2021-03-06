# TFUI

:warning: This is a WIP and helps us to learn Go, so please advice on any best practices and ignore bugs (or even submit issues or pull requests).

`tfui` is a server and command line client to review [Hashicorp Terraform](https://www.terraform.io/) plans in a succinct manner.
Plans can be uploaded with a CLI and seen on a browser.

The full architectural documentation lives in [docs/](./docs/).

## Usage

### Server

The server can be deployed with our [Helm chart](./deploy/chart/).

For the configuration of the server, environment variables are available:

- `APPLICATION_TOKEN` to authenticate to the API (default: not set)
- `BASE_DIR` as the directory of the plan file storage (default: `./plans`)
- `PORT` for which port the application should run on (default: `8080`)

```bash
tfui server
```

### Client

Environment variables available are:

- `TFUI_ADDR` as location of the TFUI server (default: `http://localhost:8080`)
- `TFUI_TOKEN` to authenticate to the API (default: not set)

```bash
tfui is a tool to manage the Terraform UI server, e.g. upload plans, or reset the server.

Usage:
  tfui [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  config      Manages local config
  help        Help about any command
  plan        Commands to manage the Terraform plans
  server      Start the server
  system      Commands to manage the system

Flags:
  -h, --help      help for tfui
  -v, --version   version for tfui

Use "tfui [command] --help" for more information about a command.
```

### Integrating into CI/CD workflow

Requirements:

1. The environment variables `REGION`, `PROJECT` and `WORKSPACE` have to be set.
2. The upload must run from a directory within a Git repository, so that the current commit hash can be parsed and added to the plan's meta information.

```bash
terraform plan -out infra.tfplan
terraform show -json infra.tfplan > infraplan.json
tfui plan upload -f infraplan.json
```

## Developing, building & running

### Make Targets

| Make Target    | Description                                          |
|----------------|------------------------------------------------------|
| `build`        | Build the binary                                     |
| `format`       | Auto-format the code to conform with common Go style |
| `lint`         | Run the linter to enforce best practices             |
| `test`         | Run all tests                                        |
| `release`      | Cross-compile the binary for OS X and Linux          |
| `docker-build` | Build docker container                               |
| `docker-run`   | Run docker container                                 |
| `helm-deploy`  | Deploys the Helm chart into a K8s cluster            |

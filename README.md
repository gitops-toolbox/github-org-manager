# Github org manager

## Setup

This cli tool uses cobra/viper to load params.
It will look for a `~/.github-org-manager` file in your home directory.

To start using the tool you should create the above file with the following content:

```yaml
---
github-token: <github token>
```

Add new command:

```bash
~/go/bin/cobra add <command>
```

## Run tests

```
GITHUB_TOKEN=<github_token> go test gitops-toolbox/github-org-manager/lib
```
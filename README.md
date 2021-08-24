# gh repo-explore

A [gh](https://github.com/cli/cli) extension for interactively exploring a repo without cloning the entire repo.

## installation

```sh
gh extension install samcoe/gh-repo-explore
```

## usage

```sh
# explore a repo
gh repo-explore samcoe/gh-repo-explore

# explore specific branch of a repo
gh repo-explore samcoe/gh-repo-explore --branch trunk
```

Supports `--hostname` flag to explore repos on hosts other than `github.com`.

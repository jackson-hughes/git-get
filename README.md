# git-get

## About
`git-get` is a simple utility that clones git repositories into a go get style directory structure.

For example, this command:

    git-get git@github.com:jackson-hughes/git-get.git

Produces:

    <GIT_GET_DIR>/github.com/jackson-hughes/git-get

## Installation

Install directly with Go:

    go install github.com/jackson-hughes/git-get@latest

Alternatively, clone the repository and build or install from source. The [Taskfile](https://taskfile.dev/) provides build and install commands, e.g:

    task install

## Config

`git-get` is configured via environment variables:

| Variable | Required | Description |
| --- | --- | --- |
| `GIT_GET_DIR` | Yes | Root directory that repositories are cloned into |
| `GIT_GET_DEBUG` | No | Set to `true` to enable debug logging |

I set `GIT_GET_DIR` in my `~/.zshrc`, e.g:

    export GIT_GET_DIR="$HOME/Projects"

# git-get

## About
`git-get` is a simple utility that clones git repositories into a go get style directory structure.

For example, this command:

    git-get git@github.com:jackson-hughes/git-get.git

Produces:

    <GIT_GET_DIR>/github.com/jackson-hughes/git-get

## Installation

Clone the repository and either build the binary or install the binary to `$GOBIN`.

The Makefile provides build and install targets, e.g:

    make install

Or simply:

    go install

## Config

Set the root directory that `git-get` will use by setting the `GIT_GET_DIR` environment variable. I set mine in my `~/.zshrc`, e.g:

    export GIT_GET_DIR="$HOME/Projects"

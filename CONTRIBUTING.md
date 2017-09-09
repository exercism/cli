# Contributing Guide

First, thank you! :tada:
Exercism would be impossible without people like you being willing to spend time and effort making things better.

## Dependencies

You'll need Go version 1.7 or higher. Follow the directions on http://golang.org/doc/install

## Development

If you've never contributed to a Go project before this is going to feel a little bit foreign.

The TL;DR is: **don't clone your fork**, and it matters where on your filesystem the project gets cloned to.

If you don't care how and why and just want something that works, follow these steps:

1. [fork this repo][fork]
1. `go get github.com/exercism/cli/exercism`
1. `cd $GOPATH/src/github.com/exercism/cli` (or `cd %GOPATH%/src/github.com/exercism/cli` on Windows)
1. `git remote set-url origin https://github.com/<your-github-username>/cli`
1. `go get -u github.com/golang/dep/cmd/dep`
1. `dep ensure`

Then make the change as usual, and submit a pull request. Please provide tests for the changes where possible.

If you care about the details, check out the blog post [Contributing to Open Source Repositories in Go][contrib-blog] on the Splice blog.

## Running the Tests

To run the tests locally on Linux or MacOS, use

```
go test $(go list ./... | grep -v vendor)
```

On Windows, the command is more painful (sorry!):

```
for /f "" %G in ('go list ./... ^| find /i /v "/vendor/"') do @go test %G
```

As of Go 1.9 this is simplified to `go test ./...`.

## Manual Testing against Exercism

You can build whatever is in your local, working copy of the CLI without overwriting your existing Exercism
CLI installation by using the `go build` command:

```
go build -o testercism exercism/main.go
```

This assumes that you are standing at the root of the exercism/cli repository checked out locally, and it will put a binary named `testercism` in your current working directory.

You can call it whatever you like, but `exercism` would conflict with the directory that is already there.

Then you call it with `./testercism`.

You can always put this in your path if you want to run it from elsewhere on your system.

We highly recommend spinning up a local copy of Exercism to test against so that you can mess with the database (and so you don't accidentally break stuff for yourself in production).

[TODO: link to the nextercism repo installation instructions, and explain how to reconfigure the CLI]

### Building for All Platforms

In order to cross-compile for all platforms, run `bin/build-all`. The binaries
will be built into the `release` directory.

[fork]: https://github.com/exercism/cli/fork
[contrib-blog]: https://splice.com/blog/contributing-open-source-git-repositories-go/

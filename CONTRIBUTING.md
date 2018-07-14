# Contributing Guide

First, thank you! :tada:
Exercism would be impossible without people like you being willing to spend time and effort making things better.

## Dependencies

You'll need Go version 1.8 or higher. Follow the directions on http://golang.org/doc/install

You will also need `dep`, the Go dependency management tool. Follow the directions on https://golang.github.io/dep/docs/installation.html

## Development

If you've never contributed to a Go project before this is going to feel a little bit foreign.

The TL;DR is: **don't clone your fork**, and it matters where on your filesystem the project gets cloned to.

If you don't care how and why and just want something that works, follow these steps:

1. [fork this repo on the Github webpage][fork]
1. `go get github.com/exercism/cli/exercism`
1. `cd $GOPATH/src/github.com/exercism/cli` (or `cd %GOPATH%/src/github.com/exercism/cli` on Windows)
1. `git remote rename origin upstream`
1. `git remote add origin git@github.com:<your-github-username>/cli.git`
1. `git checkout -b development`
1. `git push -u origin development` (setup where you push to, check it works)
1. `go get -u github.com/golang/dep/cmd/dep`
   * depending on your setup, you may need to install `dep` by following the instructions in the [`dep` repo](https://github.com/golang/dep)
1. `dep ensure`
1. `git update-index --assume-unchanged Gopkg.lock` (prevent your dep changes being committed)

Then make changes as usual and submit a pull request. Please provide tests for the changes where possible.

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

If you want to test your changes while doing your everyday exorcism work you could do:

On Unices:

- `cd $GOPATH/src/github.com/exercism/cli/exercism && go build -o exercism main.go && mv exercism ~/bin`

On Windows:

- ?? TODO

### Building for All Platforms

In order to cross-compile for all platforms, run `bin/build-all`. The binaries
will be built into the `release` directory.

[fork]: https://github.com/exercism/cli/fork
[contrib-blog]: https://splice.com/blog/contributing-open-source-git-repositories-go/

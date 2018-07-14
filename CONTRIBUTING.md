# Contributing Guide

First, thank you! :tada:
Exercism would be impossible without people like you being willing to spend time and effort making things better.

## Dependencies

You'll need Go version 1.9 or higher. Follow the directions on http://golang.org/doc/install

You will also need `dep`, the Go dependency management tool. Follow the directions on https://golang.github.io/dep/docs/installation.html

## Development

If you've never contributed to a Go project before this is going to feel a little bit foreign.

The TL;DR is: **don't clone your fork**, and it matters where on your filesystem the project gets cloned to.

If you don't care how and why and just want something that works, follow these steps:

1. [fork this repo on the GitHub webpage][fork]
1. `go get github.com/exercism/cli/exercism`
1. `cd $GOPATH/src/github.com/exercism/cli` (or `cd %GOPATH%\src\github.com\exercism\cli` on Windows)
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

To run the tests locally

```
go test ./...
```

## Manual Testing against Exercism

To test your changes while doing everyday Exercism work you
can build using the following instructions. Any name may be used for the
binary (e.g. `testercism`) - by using a name other than `exercism` you
can have different profiles under `~/.config` and avoid possibly
damaging your real Exercism submissions, or test different tokens, etc.

On Unices:

- `cd $GOPATH/src/github.com/exercism/cli/exercism && go build -o testercism main.go`
- `./testercism -h`

On Windows:

- `cd /d %GOPATH%\src\github.com\exercism\cli`
- `go build -o testercism.exe exercism\main.go`
- `testercism.exe —h`

### Building for All Platforms

In order to cross-compile for all platforms, run `bin/build-all`. The binaries
will be built into the `release` directory.

[fork]: https://github.com/exercism/cli/fork
[contrib-blog]: https://splice.com/blog/contributing-open-source-git-repositories-go/

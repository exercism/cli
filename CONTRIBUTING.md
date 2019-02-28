# Contributing Guide

First, thank you! :tada:
Exercism would be impossible without people like you being willing to spend time and effort making things better.

## Documentation
* [Exercism Documentation Repository](https://github.com/exercism/docs)
* [Exercism Glossary](https://github.com/exercism/docs/blob/master/about/glossary.md)
* [Exercism Architecture](https://github.com/exercism/docs/blob/master/about/architecture.md)

## Dependencies

You'll need Go version 1.11 or higher. Follow the directions on http://golang.org/doc/install

You will also need to be familiar with the Go `modules` dependency management system. Refer to the [modules wiki page](https://github.com/golang/go/wiki/Modules) to learn more.

## Development

A typical development workflow looks like this:

1. [fork this repo on the GitHub webpage][fork]
1. `cd /path/to/the/development/directory`
1. `git clone https://github.com/<your-github-username>/cli.git`
1. `cd cli`
1. `git remote add upstream https://github.com/exercism/cli.git`
1. Optionally: `git config user.name <your-github-username>` and `git config user.email <your-github-email>` 
1. `git checkout -b <development-branch-name>`
1. `git push -u origin <development-branch-name>` (setup where you push to, check it works)

Then make your desired changes and submit a pull request. Please provide tests for the changes where possible.

Please note that if your development directory is located inside the `GOPATH`, you would need to set the `GO111MODULE=on` environment variable, in order to be able to use the `modules` system. 

If you wish to learn how to contribute to the Go projects without the `modules`, check out the blog post [Contributing to Open Source Repositories in Go][contrib-blog] on the Splice blog.

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

- `cd /path/to/the/development/directory/cli && go build -o testercism main.go`
- `./testercism -h`

On Windows:

- `cd /d \path\to\the\development\directory\cli`
- `go build -o testercism.exe exercism\main.go`
- `testercism.exe —h`

### Building for All Platforms

In order to cross-compile for all platforms, run `bin/build-all`. The binaries
will be built into the `release` directory.

[fork]: https://github.com/exercism/cli/fork
[contrib-blog]: https://splice.com/blog/contributing-open-source-git-repositories-go/

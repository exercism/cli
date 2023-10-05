# Contributing Guide

First, thank you! :tada:
Exercism would be impossible without people like you being willing to spend time and effort making things better.

## Documentation
* [Exercism Documentation Repository](https://github.com/exercism/docs)

## Dependencies

You'll need Go version 1.20 or higher. Follow the directions on http://golang.org/doc/install

## Development

This project uses Go's [`modules` dependency management](https://github.com/golang/go/wiki/Modules) system.

To contribute [fork this repo on the GitHub webpage][fork] and clone your fork.
Make your desired changes and submit a pull request.
Please provide tests for the changes where possible.

Please note that if your development directory is located inside the `GOPATH`, you need to set the `GO111MODULE=on` environment variable.

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

- `cd /path/to/the/development/directory/cli && go build -o testercism ./exercism/main.go`
- `./testercism -h`

On Windows:

- `cd /d \path\to\the\development\directory\cli`
- `go build -o testercism.exe exercism\main.go`
- `testercism.exe —h`

### Releasing a new CLI version
Consult the [release documentation](RELEASE.md).

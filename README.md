[![Build Status](https://travis-ci.org/exercism/cli.png?branch=master)](https://travis-ci.org/exercism/cli)
[![Supporting 24 Pull Requests](https://img.shields.io/badge/Supporting-24%20Pull%20Requests-red.svg?style=flat)](http://24pullrequests.com)

# Exercism Command-Line Client

The CLI provides a way to do the problems on
[exercism.io](http://exercism.io).

**Important**: If you're looking for instructions on how to install the CLI. Please read [Installing the CLI](http://exercism.io/cli)

This CLI ships as a binary with no additional runtime requirements. This means
that if you're doing the Haskell problems on exercism you don't need a working
Python or Ruby environment simply to fetch and submit exercises.

## Dependencies

Go version 1.6 or higher

## Installing Go

Follow the directions on http://golang.org/doc/install

## Development

1. fork this repo
1. `go get github.com/exercism/cli/exercism`
1. `cd $GOPATH/src/github.com/exercism/cli`
1. `git remote set-url origin https://github.com/<your-github-username>/cli`
1. `go get -t ./...`
1. Make the change.
1. Submit a pull request.

Please provide tests for the changes where possible.

To run the tests locally, use `go test ./...`

At the moment the CLI commands are not tested, so if you're adding a new
command don't worry too hard about tests.

## Building

To build the binary for your platform run

```
go install github.com/exercism/cli/exercism
```

or

```
go build -o out/exercism exercism/main.go
```

The resulting binary can be found in `out/exercism` (Linux, Mac OS X) or `out/exercism.exe` (Windows).

In order to cross-compile for all platforms, run `bin/build-all`. The binaries
will be built into the `release` directory.

## Domain Concepts

- **Language** is the name of a programming language. E.g. C++ or Objective-C or JavaScript.
- **Track ID** is a normalized, url-safe identifier for a language track. E.g. `cpp` or `objective-c` or `javascript`.
- **Problem** is an exercism exercise.
- **Problem Slug** is a normalized, url-safe identifier for a problem.
- **Iteration** is a solution that a user has written for a particular problem in a particular language track. A user may have several iterations for the same problem.



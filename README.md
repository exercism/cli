[![Build Status](https://travis-ci.org/exercism/cli.png?branch=master)](https://travis-ci.org/exercism/cli)

# Exercism Command-Line Client

The CLI provides a way to do the problems on
[exercism.io](http://exercism.io).

This CLI ships as a binary with no additional runtime requirements. This means
that if you're doing the Haskell problems on exercism you don't need a working
Python or Ruby environment simply to fetch and submit exercises.

## Installing Go

Follow the directions on http://golang.org/doc/install

## Development

1. fork this repo
1. `go get github.com/exercism/cli/exercism`
1. `cd $GOPATH/src/github.com/exercism/cli`
1. `git remote set-url origin https://github.com/<your-github-username>/cli`
1. `go get`
1. Make sure $GOPATH/bin is on your path (you may need something like `export PATH=$PATH:/projects/goprojects/bin`)
1. Open a separate terminal window to your project directory and run the command `glitch`
1. Make the change.
1. Submit a pull request.

Please provide tests for the changes where possible.

At the moment the CLI commands are not tested, so if you're adding a new
command don't worry too hard about tests.

## Building

To build the binary for your platform run

```
bin/build
```

The resulting binary can be found in `out/exercism` (Linux, Mac OS X) or `out/exercism.exe` (Windows).

In order to cross-compile for all platforms, run `bin/build-all`. The binaries
will be built into the `release` directory.

## Troubleshooting

```plain
app.Run(os.Args) used as value
```

This error is due to a breaking change between the 0.x version of the `codegangsta/cli` library and the `1.x` version of the library.

To fix it update the `codegangsta/cli` dependency:

```plain
$ go get -u github.com/codegangsta/cli
```

## Using Glitch

If you'd like to run lint, vet, and the tests on every change, install Levi Cook's `glitch` library:

1. `go get github.com/levicook/glitch`
1. `go install github.com/levicook/glitch`
1. Ensure that you have `go vet`
1. Run it with `glitch`.

### Troubleshooting.

When you `glitch`, do you get stymied like this?

```shell
# github.com/exercism/cli
api_test.go:7:2: cannot find package "github.com/stretchr/testify/assert" in any of: ...
FAIL	github.com/exercism/cli [setup failed]
```

You may need to

```shell
$ go get github.com/stretchr/testify/assert
$ go install github.com/stretchr/testify/assert
```

Now you should be able to run `glitch`.

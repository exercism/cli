[![Build Status](https://travis-ci.org/exercism/cli.png?branch=master)](https://travis-ci.org/exercism/cli)

Goals
===========

Provide developers an easy way to work with [exercism.io](http://exercism.io) that doesn't require a 
Ruby environment.

Installing Go
=============

### On Mac OS X

You may get away with ```brew install go --cross-compile-common``` unless you have the latest XCode, which does not ship with gcc.

If have the latest XCode, try ```brew install go --cross-compile-common --without-cgo```.

If that throws an error, try ```brew install go --cross-compile-common --with-llvm```.

Development
===========
1. `cd $GOPATH`
1. git clone git@github.com:exercism/cli.git src/github.com/exercism/cli
1. cd src/github.com/exercism/cli
1. go get
1. go get github.com/levicook/glitch
1. go install github.com/levicook/glitch
1. Make sure $GOPATH/bin is on your path (you may need something like `export PATH=$PATH:/projects/goprojects/bin`)
1. Open a separate terminal window to your project directory and run the command `glitch`
1. Write a test.
1. Watch test fail.
1. Make test pass.
1. Submit a pull request.

Building
========
1. Run ```bin/build``` and the binary for your platform will be built into the out directory.
1. Run ```bin/build-all``` and the binaries for OSX, Linux and Windows will be built into the release directory.

Troubleshooting
===============

```plain
app.Run(os.Args) used as value
```

This error is due to a breaking change between the 0.x version of the `codegangsta/cli` library and the `1.x` version of the library.

To fix it update the `codegangsta/cli` dependency:

```plain
$ go get -u github.com/codegangsta/cli
```

New to go?  Missing packages for glitch?
-------------------------------------------

### missing assertion library

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
$ glitch
```

### `go vet` on MacOS

Depending on your brew installation of `go` you may not yet have the `vet` command available.  You may also need `hg` (mercurial) to get rolling.

Here's a sample (trimmed) output from a successful installation of `vet` on MacOSX 10.8.5 with XCode 5.0.2

```shell
% brew install go --cross-compile-common

... installation output ...

% go vet
go tool: no such tool "vet"; to install:
	go get code.google.com/p/go.tools/cmd/vet

% go get code.google.com/p/go.tools/cmd/vet
go: missing Mercurial command. See http://golang.org/s/gogetcmd

% brew install hg
% go get code.google.com/p/go.tools/cmd/vet
% go vet

```

Now you should be able to run `glitch`.



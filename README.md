[![Build Status](https://travis-ci.org/msgehard/go-exercism.png?branch=master)](https://travis-ci.org/msgehard/go-exercism)

Goals
===========

Provide developers an easy way to work with [exercism.io](http://exercism.io) that doesn't require a 
Ruby environment.

Installing Go
=============

### On Mac OS X

You may get away with ```brew install go --cross-compile-common``` unless you have the latest XCode, which does not ship with gcc.

If have the latest XCode, try ```brew install go --cross-compile-common --without-cgo```.

If that throws an error, try ```brew install go --crosscompile-commone --with-llvm```.

Development
===========
1. Fork and clone into your `$GOPATH/src`
1. `go get`
1. `go install github.com/levicook/glitch`
1. Open a separate terminal window to your project directory and run the command `glitch`
1. Write a test.
1. Watch test fail.
1. Make test pass.
1. Submit a pull request.

Building
========
1. Run ```bin/build``` and the binary for your platform will be built into the out directory.
1. Run ```bin/build-all``` and the binaries for OSX, Linux and Windows will be built into the release directory.

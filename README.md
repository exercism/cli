[![Build Status](https://travis-ci.org/msgehard/go-exercism.png?branch=master)](https://travis-ci.org/msgehard/go-exercism)

Goals
===========

Provide developers an easy way to work with [exercism.io](http://exercism.io) that doesn't require a 
Ruby environment.

This tool is in beta testing. All of the major functionality has been implemented. Please help
us work out the kinks by using it to access the site. 

If you want something more stable to access exercism.io, please
see the [ruby gem](https://github.com/kytrinyx/exercism).

Development
===========
1. Install Go ```brew install go --cross-compile-common``` or the command appropriate for your platform.
1. Fork and clone.
1. Run ```git submodule update --init --recursive```
1. Write a test.
1. Run ``` bin/test ``` and watch test fail.
1. Make test pass.
1. Submit a pull request.

Building
========
1. Run ```bin/build``` and the binary for your platform will be built into the out directory.
1. Run ```bin/build-all``` and the binaries for OSX, Linux and Windows will be built into the release directory.

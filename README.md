[![Build Status](https://travis-ci.org/msgehard/go-exercism.png?branch=master)](https://travis-ci.org/msgehard/go-exercism)

Goals
===========

Provide non-Ruby developers an easy way to work with [exercism.io](http://exercism.io).

This tool is under heavy development. If you want something more stable to access exercism.io, please
see the [ruby gem](https://github.com/kytrinyx/exercism).

Development
===========
1. Install Go ```brew install go``` or the command appropriate for your platform.
1. Fork and clone.
1. Run ```git submodule update --init --recusive```
1. Write a test.
1. Run ``` bin/test ``` and watch test fail.
1. Make test pass.
1. Submit a pull request.

Building
========
1. Run ```./bin/build```
1. The binary will be built into the out directory.

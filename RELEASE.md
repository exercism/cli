# Cutting a CLI Release

## Bootstrap Cross-Compilation for Go

**This only has to be done once.**

Change directory to the go source. Then run the bootstrap command for
each operating system and architecture.

```plain
$ cd `which go`/../../src
$ sudo GCO_ENABLED=0 GOOS=windows GOARCH=386 ./make.bash --no-clean
$ sudo GCO_ENABLED=0 GOOS=darwin GOARCH=386 ./make.bash --no-clean
$ sudo GCO_ENABLED=0 GOOS=linux GOARCH=386 ./make.bash --no-clean
$ sudo GCO_ENABLED=0 GOOS=windows GOARCH=amd64 ./make.bash --no-clean
$ sudo GCO_ENABLED=0 GOOS=darwin GOARCH=amd64 ./make.bash --no-clean
$ sudo GCO_ENABLED=0 GOOS=linux GOARCH=amd64 ./make.bash --no-clean
$ sudo GCO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=6 ./make.bash --no-clean
$ sudo GCO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=5 ./make.bash --no-clean
```

## Update the Changelog

Make sure all the recent changes are reflected in the "next release" section
of the Changelog. Make this a separate commit from bumping the version.

You can view changes using the /compare/ view:
https://github.com/exercism/cli/compare/$PREVIOUS_RELEASE...master

## Bump the version

Edit the `Version` constant in `cmd/version.go`, and edit the Changelog.

All the changes in the "next release" section should be moved to a new section
that describes the version number, and gives it a date.

The "next release" section should contain only "Your contribution here".

_Note: It's useful to add the version to the commit message when you bump it: e.g. `Bump version to v2.3.4`._

## Generate the Binaries

```plain
$ rm release/*
$ CGO_ENABLED=0 bin/build-all
```

## Cut Release on GitHub

Go to [the exercism/cli "new release" page](https://github.com/exercism/cli/releases/new).

Describe the release, select a specific commit to target, name the version `v{VERSION}`, where
VERSION matches the value of the `Version` constant.

Upload all the binaries from `release/*`.

Paste the release text and describe the new changes (`tail -n +57 RELEASE.md | head -n 16 | pbcopy`):

```
### Exercism Command-Line Interface (CLI)

Exercism takes place in two places: the discussions happen on the website, and you work on exercises locally. The CLI bridges the gap, allowing you to fetch exercises and submit solutions to the site.

This is a stand-alone binary, which means that you don't need to install any particular language or environment in order to use it.

To install, download the archive that matches your operating system and architecture, unpack the archive, and put the binary somewhere on your path.

You will need to configure the CLI with your [Exercism API Key](http://exercism.io/account/key) before submitting.

For more detailed instructions, see the [CLI page on Exercism](http://exercism.io/cli).

#### Recent changes

* ABC...
* XYZ...
```

## Update Homebrew

This is helpful for the (many) Mac OS X users.

First, get a copy of the latest tarball of the source code:

```
cd ~/tmp && wget https://github.com/exercism/cli/archive/vX.Y.Z.tar.gz
```

Get the SHA256 of the tarball:

```
shasum -a 256 vX.Y.Z.tar.gz
```

Update the homebrew formula:

```
cd $(brew --repository)
git checkout master
brew update
brew bump-formula-pr --strict exercism --url=https://github.com/exercism/cli/archive/vX.Y.Z.tar.gz --sha256=$SHA
```

For more information see their [contribution guidelines](https://github.com/Homebrew/homebrew/blob/master/share/doc/homebrew/How-To-Open-a-Homebrew-Pull-Request-(and-get-it-merged).md#how-to-open-a-homebrew-pull-request-and-get-it-merged).

## Update the docs site

If there are any significant changes, we should describe them on
[exercism.io/cli]([https://exercism.io/cli).

The codebase lives at [exercism/website-copy](https://github.com/exercism/website-copy) in `pages/cli.md`.

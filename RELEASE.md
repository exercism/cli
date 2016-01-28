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

## Bump the version

Edit the `Version` constant in `exercism/main.go`, and edit the Changelog.

All the changes in the "next release" section should be moved to a new section
that describes the version number, and gives it a date.

The "next release" section should contain only "Your contribution here".

_Note: It's useful to add the version to the commit message when you bump it: e.g. `Bump version to v2.3.4`.

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

Past the release text and describe the new changes:

```
### Exercism Command-Line Interface (CLI)

Exercism takes place in two places: the discussions happen on the website, and you work on exercises locally. The CLI bridges the gap, allowing you to fetch exercises and submit solutions to the site.

This is a stand-alone binary, which means that you don't need to install any particular language or environment in order to use it.

To install, download the archive that matches your operating system and architecture, unpack the archive, and put the binary somewhere on your path.

You will need to configure the CLI with your [Exercism API Key](http://exercism.io/account/key) before submitting.

For more detailed instructions, see the [CLI page on Exercism](http://exercism.io/cli).

#### Recent Changes

* ABC...
* XYZ...
```

Calculate the `SHA1` checksum for the two mac builds:

```plain
$ openssl sha1 release/exercism-mac-32bit.tgz
$ openssl sha1 release/exercism-mac-64bit.tgz
```

Update the homebrew-binary/exercism.rb formula.

- version
- urls
- `SHA1` checksums

If you are on a mac, you can test the formula by running:

```plain
$ brew unlink exercism && brew install ./exercism.rb
```

Then submit a pull request to homebrew-binary. Note that they're very, very careful about their
commit history. If you made multiple commits, squash them. They don't merge using the button on
GitHub, they merge by making sure that the branch is rebased on to the most recent master, and
then they do a fast-forward merge. That means that the PR will be red when it's closed, not purple.

Also, don't bother trying to DRY out the formula, they prefer having explicit, hard-coded values.

## Update the Docs Site

If there are any significant changes, we should describe them on
[cli.exercism.io](http://cli.exercism.io/).

The codebase lives at [exercism/cli-www](https://github.com/exercism/cli-www).

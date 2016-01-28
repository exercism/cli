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

## Update Homebrew

This is helpful for the (many) Mac OS X users.

Fork and clone [homebrew/homebrew-binary](https://github.com/homebrew/homebrew-binary/fork).

Add the upstream repository:

```plain
$ git remote add upstream git@github.com:Homebrew/homebrew-binary.git
```

If you already had this cloned, ensure that you are entirely up-to-date with the upstream master:

```plain
$ git fetch upstream && git checkout master && git reset upstream/master && git push -f origin master
```

Check out a feature branch, where X.Y.Z is the actual version number.

```plain
$ git checkout -b exercism-vX.Y.Z
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

Update the formula:

```
cd $(brew --repository)
git checkout master
git remote add YOUR_USERNAME git@github.com:YOUR_USERNAME/homebrew.git
brew update
git checkout -b exercism-vX.Y.Z
brew edit exercism
# update sha256 and tarball url
brew audit exercism
brew install exercism
brew test exercism
git commit -m "exercism X.Y.Z"
git push --set-upstream YOUR_USERNAME exercism-vX.Y.Z
```

Then go to https://github.com/Homebrew/homebrew and create pull request.

Note that they really don't want any verbose commit messages or PR descriptions when all you're doing is bumping a version.

For more information see their [contribution guidelines](https://github.com/Homebrew/homebrew/blob/master/share/doc/homebrew/How-To-Open-a-Homebrew-Pull-Request-(and-get-it-merged).md#how-to-open-a-homebrew-pull-request-and-get-it-merged).

## Update the Docs Site

If there are any significant changes, we should describe them on
[cli.exercism.io](http://cli.exercism.io/).

The codebase lives at [exercism/cli-www](https://github.com/exercism/cli-www).

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

## Bump the version

Edit the `Version` constant in `exercism/main.go`.

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

Calculate the `SHA1` checksum for the two mac builds:

```plain
$ openssl sha1 release/exercism-mac-32bit.tgz
$ openssl sha1 release/exercism-mac-64bit.tgz
```

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

## Update the Docs Site

If there are any significant changes, we should describe them on
[cli.exercism.io](http://cli.exercism.io/).

The codebase lives at [exercism/cli-www](https://github.com/exercism/cli-www).

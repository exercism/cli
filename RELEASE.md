# Cutting a CLI Release

The Exercism CLI uses [GoReleaser](https://goreleaser.com) to automate the
release process. 

## Requirements

1. [Install GoReleaser](https://goreleaser.com/install/)
1. [Install snapcraft](https://snapcraft.io/docs/snapcraft-overview)
1. [Setup GitHub token](https://goreleaser.com/environment/#github-token)
1. Have a gpg key installed on your machine - it is [used for signing the artifacts](https://goreleaser.com/sign/)

## Confirm / Update the Changelog

Make sure all the recent changes are reflected in the "next release" section of the CHANGELOG.md file.  All the changes in the "next release" section should be moved to a new section that describes the version number, and gives it a date.

You can view changes using the /compare/ view:
https://github.com/exercism/cli/compare/$PREVIOUS_RELEASE...master

GoReleaser supports the [auto generation of a changelog](https://goreleaser.com/customization/#customize-the-changelog), however we would need to customize the output to meet our standards (not including refactors, test updates, etc). We should also consider using [the release notes feature](https://goreleaser.com/customization/#custom-release-notes).

## Bump the version

Edit the `Version` constant in `cmd/version.go`

_Note: It's useful to add the version to the commit message when you bump it: e.g. `Bump version to v2.3.4`._

In the future we will probably want to replace the hardcoded `Version` constant with [main.version](https://goreleaser.com/environment/#using-the-main-version). Here is a [stack overflow post on injecting to cmd/version.go](https://stackoverflow.com/a/47510909).

Commit this change on a branch along with the CHANGELOG updates in a single commit, and create a PR for merge to master.

## Cut a release

```bash
# Test run
goreleaser --skip-publish --snapshot --rm-dist

# Create a new tag on the master branch and push it
git tag -a v3.0.16 -m "Trying out GoReleaser"
git push origin v3.0.16

# Build and release
goreleaser --rm-dist

# You must be logged into snapcraft to publish a new snap
snapcraft login

# Push to snapcraft
for f in `ls dist/*.snap`; do snapcraft push --release=stable $f; done

# [TODO] Push to homebrew
```

## Cut Release on GitHub

Run [exercism-cp-archive-hack.sh](https://gist.github.com/ekingery/961650fca4e2233098c8320f32736836) which takes the new archive files and renames them to match the old naming scheme for backward compatibility. Until mid to late 2020, we will need to manually upload the backward-compatible archive files generated in `/tmp/exercism_tmp_upload`.

The generated archive files should be uploaded to the [draft release page created by GoReleaser](https://github.com/exercism/cli/releases). Describe the release, select a specific commit to target, paste the following release text, and describe the new changes.

```
To install, follow the interactive installation instructions at https://exercism.io/cli-walkthrough
---

[describe changes in this release]
```

 Lastly, test and publish the draft


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

For more information see [How To Open a Homebrew Pull Request](https://docs.brew.sh/How-To-Open-a-Homebrew-Pull-Request).

## Update the docs site

If there are any significant changes, we should describe them on
[exercism.io/cli](https://exercism.io/cli).

The codebase lives at [exercism/website-copy](https://github.com/exercism/website-copy) in `pages/cli.md`.

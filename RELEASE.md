# Cutting a CLI Release

The Exercism CLI uses [GoReleaser](https://goreleaser.com) to automate the
release process.

## Requirements

1. [Install GoReleaser](https://goreleaser.com/install/)
1. [Install snapcraft](https://snapcraft.io/docs/snapcraft-overview)
1. [Setup GitHub token](https://goreleaser.com/scm/github/)
1. Have a gpg key installed on your machine - it is [used for signing the artifacts](https://goreleaser.com/customization/sign/)

## Confirm / Update the Changelog

Make sure all the recent changes are reflected in the "next release" section of the CHANGELOG.md file. All the changes in the "next release" section should be moved to a new section that describes the version number, and gives it a date.

You can view changes using the /compare/ view:
https://github.com/exercism/cli/compare/$PREVIOUS_RELEASE...main

GoReleaser supports the [auto generation of a changelog](https://goreleaser.com/customization/#customize-the-changelog) we will want to customize to meet our standards (not including refactors, test updates, etc). We should also consider using [the release notes feature](https://goreleaser.com/customization/#custom-release-notes).

## Bump the version

Edit the `Version` constant in `cmd/version.go`

_Note: It's useful to add the version to the commit message when you bump it: e.g. `Bump version to v2.3.4`._

In the future we will probably want to replace the hardcoded `Version` constant with [main.version](https://goreleaser.com/cookbooks/using-main.version). Here is a [stack overflow post on injecting to cmd/version.go](https://stackoverflow.com/a/47510909).

Commit this change on a branch along with the CHANGELOG updates in a single commit, and create a PR for merge to main.

## Cut a release

```bash
# Test run
goreleaser --skip-publish --snapshot --clean

# Create a new tag on the main branch and push it
git tag -a v3.0.16 -m "Trying out GoReleaser"
git push origin v3.0.16

# Build and release
goreleaser --clean

# You must be logged into snapcraft to publish a new snap
snapcraft login

# Push to snapcraft
for f in `ls dist/*.snap`; do snapcraft push --release=stable $f; done

# [TODO] Push to homebrew
```

## Cut Release on GitHub

At this point, Goreleaser will a created a draft PR at https://github.com/exercism/cli/releases/tag/vX.Y.Z.
On that page, update the release description to:

```
To install, follow the interactive installation instructions at https://exercism.io/cli-walkthrough
---

[describe changes in this release]
```

Lastly, test and publish the draft

## Update Homebrew

Next, we'll submit a PR to Homebrew to update the Exercism formula (which is how macOS users usually download the CLI):

```
cd /tmp && curl -O https://github.com/exercism/cli/archive/vX.Y.Z.tar.gz
cd $(brew --repository)
git checkout master
brew update
brew bump-formula-pr --strict exercism --url=https://github.com/exercism/cli/archive/vX.Y.Z.tar.gz --sha256=$(shasum -a 256 /tmp/vX.Y.Z.tar.gz)
```

For more information see [How To Open a Homebrew Pull Request](https://docs.brew.sh/How-To-Open-a-Homebrew-Pull-Request).

## Update the docs site

If there are any significant changes, we should describe them on
[exercism.io/cli](https://exercism.io/cli).

The codebase lives at [exercism/website-copy](https://github.com/exercism/website-copy) in `pages/cli.md`.

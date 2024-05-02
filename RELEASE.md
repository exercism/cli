# Cutting a CLI Release

The Exercism CLI uses [GoReleaser](https://goreleaser.com) to automate the release process.

## Requirements

1. [Install GoReleaser](https://goreleaser.com/install/)
1. [Setup GitHub token](https://goreleaser.com/scm/github/)
1. Have a gpg key installed on your machine - it is [used for signing the artifacts](https://goreleaser.com/customization/sign/)

## Confirm / Update the Changelog

Make sure all the recent changes are reflected in the "next release" section of the CHANGELOG.md file.
All the changes in the "next release" section should be moved to a new section that describes the version number, and gives it a date.

You can view changes using the /compare/ view:
https://github.com/exercism/cli/compare/$PREVIOUS_RELEASE...main

## Bump the version

1. Create a branch for the new version
1. Edit the `Version` constant in `cmd/version.go`
1. Update the `CHANGELOG.md` file
1. Commit the updated version
1. Create a PR

_Note: It's useful to add the version to the commit message when you bump it: e.g. `Bump version to v2.3.4`._

## Cut a release

Once the version bump PR has been merged, run the following commands:

```bash
VERSION=$(sed -n -E 's/^const Version = "([0-9]+\.[0-9]+\.[0-9]+)"$/\1/p' cmd/version.go)
TAG_NAME="v${VERSION}"
GPG_FINGERPRINT="<GPG FINGERPRINT>"

# Test run
goreleaser --skip-publish --snapshot --clean

# Create a new tag on the main branch and push it
git tag -a "${TAG_NAME}" -m "Trying out GoReleaser"
git push origin "${TAG_NAME}"
```

Brew tap is now managed by `.goreleaser.yml` so no need to update it manually.
GoReleaser can generate and publish a homebrew-tap recipe into a repository
automatically. See [GoReleaser docs](https://goreleaser.com/customization/homebrew/)
for more details.

## Cut Release on GitHub

At this point, Goreleaser will have created a draft release at https://github.com/exercism/cli/releases/tag/vX.Y.Z.
On that page, update the release description to:

```
To install, follow the interactive installation instructions at https://exercism.io/cli-walkthrough
---

[modify the generated release-notes to describe changes in this release]
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

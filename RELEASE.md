# Cutting a CLI Release

The Exercism CLI uses [GoReleaser](https://goreleaser.com) to automate the release process.

## Requirements

1. [Install GoReleaser](https://goreleaser.com/install/)
1. [Setup GitHub token](https://goreleaser.com/scm/github/)
1. Have a gpg key installed on your machine - it is [used for signing the artifacts](https://goreleaser.com/customization/sign/)

## Bump the version

1. Create a branch for the new version
1. Bump the `Version` constant in `cmd/version.go`
1. Update the `CHANGELOG.md` file to include a section for the new version and its changes.
   Hint: you can view changes using the compare view: https://github.com/exercism/cli/compare/$PREVIOUS_RELEASE...main.
1. Commit the updated files
1. Create a PR

_Note: It's useful to add the version to the commit message when you bump it: e.g. `Bump version to v2.3.4`._

## Cut a release

Once the version bump PR has been merged, run the following command to cut a release:

```shell
GPG_FINGERPRINT="<THE_GPG_FINGERPRINT>" ./bin/release.sh
```

## Cut Release on GitHub

Once the `./bin/release.sh` command finishes, the [release workflow](https://github.com/exercism/cli/actions/workflows/release.yml) will automatically run.
This workflow will create a draft release at https://github.com/exercism/cli/releases/tag/vX.Y.Z.
Once created, go that page to update the release description to:

```
To install, follow the interactive installation instructions at https://exercism.org/cli-walkthrough
---

[modify the generated release-notes to describe changes in this release]
```

Lastly, test and then publish the draft.

## Homebrew

Homebrew will automatically bump the version, no manual action is required.

## Update the docs site

If there are any significant changes, we should describe them on
[exercism.org/cli](https://exercism.org/cli).

The codebase lives at [exercism/website-copy](https://github.com/exercism/website-copy) in `pages/cli.md`.

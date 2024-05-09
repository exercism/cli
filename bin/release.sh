#!/usr/bin/env bash

set -euo pipefail

if [[ -z "${GPG_FINGERPRINT}" ]]; then
    echo "GPG_FINGERPRINT environment variable is not set"
    exit 1
fi

echo "Syncing repo with latest main..."
git checkout main
git pull

VERSION=$(sed -n -E 's/^const Version = "([0-9]+\.[0-9]+\.[0-9]+)"$/\1/p' cmd/version.go)
TAG_NAME="v${VERSION}"

echo "Verify release can be built..."
goreleaser --skip=publish --snapshot --clean

echo "Pushing tag..."
git tag -a "${TAG_NAME}" -m "Release ${TAG_NAME}"
git push origin "${TAG_NAME}"

echo "Tag pushed"
echo "The release CI workflow will automatically create a draft release."
echo "Once created, edit the release notes and publish it."

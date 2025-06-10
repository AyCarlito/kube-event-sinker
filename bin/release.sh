#!/bin/bash

# Creates a release branch in the repository using the latest tag.

latest_tag=$(git describe --tags --abbrev=0)
release_branch="releases/release-${latest_tag}"

git checkout -b "${release_branch}"
git push --set-upstream origin "${release_branch}"
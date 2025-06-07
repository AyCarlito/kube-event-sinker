#!/bin/bash

# Return a SemVer compliant version for the current commit.
# - major, minor and patch numbers are the latest git tag (with the "v" prefix removed).
# "-build" is appended in the case of:
#       - local builds on all branches
#       - CI builds on non-release branches
# build is defined as:
#       - dev.$BRANCH_NAME.$BUILD_ID (CI)
#       - dev.branch.sha (Local)

major_minor_patch=$(git describe --tags --abbrev=0 | cut -c2-)

current_branch="${BRANCH_NAME}"
if [ -z "${BRANCH_NAME}" ]; then
    current_branch="$(git symbolic-ref --short HEAD)" || current_branch="detached" # detached HEAD
fi

# If the branch is prefixed with a number corresponding to an issue, use it.
# Otherwise we use what is remaining after stripping the non alphanumeric characters.
branch_number_prefix="${current_branch%%[!0-9]*}"
readonly branch_number_prefix
if [ "${branch_number_prefix}" ]; then
    branch_prefix="${branch_number_prefix}"
else
    branch_prefix="${current_branch//[^[:alnum:]]/}"
fi

# SHA for local builds.
if [ -z "${BUILD_ID}" ]; then
    current_sha="$(git rev-parse --short HEAD)"
    build="dev.${branch_prefix}.${current_sha}"
else
    # Build ID for CI builds.
    build="dev.${branch_prefix}.${BUILD_ID}"
fi

# Special case when building on a release branch
version="${major_minor_patch}-${build}"
if [ "${current_branch:0:7}" = "release" ]; then
    version="${major_minor_patch}"
fi

echo "${version}"
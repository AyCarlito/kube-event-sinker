#!/bin/bash

GIT_ROOT="$(git rev-parse --show-toplevel)"
readonly GIT_ROOT

PULL_REQUEST_BASE="https:\/\/github.com\/AyCarlito\/kube-event-sinker\/pull"
readonly PULL_REQUEST_BASE

latest_tag=$(git describe --tags --abbrev=0)
latest_tag_commit=$(git rev-list -n 1 "${latest_tag}")
previous_tag=$(git describe --abbrev=0 --tags "$(git rev-list --tags --skip=1 --max-count=1)")
previous_tag_commit=$(git rev-list -n 1 "${previous_tag}")

changelog="${GIT_ROOT}/CHANGELOG/CHANGELOG-$latest_tag.md"
readonly changelog

# Retrieve all the commits between the previous tag and latest tag.
commits=$(git log --oneline --no-decorate --pretty=format:"%s" "${previous_tag_commit}".."${latest_tag_commit}")

breaking_changes=$(echo "${commits}" | grep -E "^BREAKING(.*):" | awk -F':' '{print "-" $2}')
features=$(echo "${commits}" | grep -E "^feat(.*):" | awk -F':' '{print "-" $2}')
fixes=$(echo "${commits}" | grep -E "^fix(.*):" | awk -F':' '{print "-" $2}')

mkdir -p CHANGELOG
echo "# ${latest_tag}" > "${changelog}"

if [ -n "${breaking_changes}" ]; then
    printf "\nBREAKING CHANGES:\n\n%s\n" "${breaking_changes}" >> "${changelog}"
fi

if [ -n "${features}" ]; then
    printf "\nFeatures:\n\n%s\n" "${features}" >> "${changelog}"
fi

if [ -n "${fixes}" ]; then
    printf "\nFixes:\n\n%s\n" "${fixes}" >> "${changelog}"
fi

# Consider the pattern (#19), we replace it with "[#19](https://github.com/AyCarlito/kube-event-sinker/pull/19)".
sed -i -E "s|(#([0-9]+))|[#\2](${PULL_REQUEST_BASE}/\2)|g" "${changelog}"
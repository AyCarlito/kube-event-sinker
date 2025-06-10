#!/bin/bash

# Pushes a new tag as needed to the remote.
# If the latest tag is associated with the current commit, no action is taken.
# If the latest tag is associated with a previous commit:
#   - Retrieve all commits between the current commit and commit associated with the latest tag.
#       - If any commit contains the string "BREAKING CHANGE"; increment the major, reset the minor and patch.
#       - If any commit contains the string "feat"; increment the minor and reset the patch version.
#       - if any commit contains the string "fix"; increment the patch version.
#   - The above 3 actions are mutually exclusive with one another.
#   - The new tag (if there is one) is then pushed to the remote.

# Do nothing if the latest tag corresponds to the current commit.
if git describe --tags --exact-match &>/dev/null; then
    echo "Latest commit is tagged."
    exit 0
fi


latest_tag=$(git describe --tags --abbrev=0)
latest_tag_commit=$(git rev-list -n 1 "${latest_tag}")


major_minor_patch=$(echo "${latest_tag}" | cut -c2-)
major=$(echo "${major_minor_patch}" | cut --delimiter=. --fields=1)
minor=$(echo "${major_minor_patch}" | cut --delimiter=. --fields=2)
patch=$(echo "${major_minor_patch}" | cut --delimiter=. --fields=3)


# Retrieve all the commits between the current commit and the commit corresponding to the latest tag.
# The major, minor or patch version will be incremented based on the commit titles.
commits=$(git log --oneline --no-decorate --pretty=format:"%s" "${latest_tag_commit}"..HEAD)
if echo "$commits" | grep -E "^BREAKING(.*):" &>/dev/null; then
    ((major++))
    minor=0
    patch=0
elif echo "$commits" | grep -E "^feat(.*):" &>/dev/null; then
    ((minor++))
    patch=0
elif echo "$commits" | grep -E "^fix(.*):" &>/dev/null; then
    ((patch++))
fi


new_tag="v${major}.${minor}.${patch}"
if [[ "$latest_tag" == "$new_tag" ]]; then
    echo "Generated tag is the same as latest."
    exit 0
fi

echo "Pushing new tag: ${new_tag}"
git tag "${new_tag}"
git push origin "${new_tag}"
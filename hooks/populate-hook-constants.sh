#!/bin/bash

# Extract version from constants.go
export constantsVersion=$(sed -n 's/^const VERSION = "\([0-9\.]*\)".*/\1/p' supertokens/constants.go)
export constantsVersionXy=$(sed -n 's/^const VERSION = "\([0-9]*\.[0-9]*\).*/\1/p' supertokens/constants.go)

# Go modules use the same source for version, so these are identical
export setupVersion="$constantsVersion"
export setupVersionXy="$constantsVersionXy"

export newestVersion="$constantsVersion"

# Target branch of the PR.
if [[ "$GITHUB_BASE_REF" != "" ]]; then
    export targetBranch="$GITHUB_BASE_REF"
else
    export targetBranch=$(git branch --show-current 2> /dev/null) || export targetBranch="(unnamed branch)"
fi
export targetBranch=${targetBranch##refs/heads/}

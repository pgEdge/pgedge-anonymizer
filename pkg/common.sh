#!/usr/bin/env bash
# common.sh - Common environment variables

# Default PostgreSQL version and derived values
export PG_VERSION="${PG_VERSION:-17.7}"
export PG_MAJOR_VERSION="$(echo "$PG_VERSION" | cut -d. -f1)"

export ANONYMIZER_REPO="https://github.com/pgEdge/pgedge-anonymizer.git"
export ANONYMIZER_BRANCH="${COMPONENT_BRANCH:-v1.0.0}"
export ANONYMIZER_VERSION=${COMPONENT_VERSION:-1.0.0}
export ANONYMIZER_BUILDNUM=${COMPONENT_BUILDNUM:-1}

# DEB only: move a pre-release pretag (e.g. BUILDNUM='beta3_1') into the
# upstream VERSION with a leading '~' (1.0.0~beta3, BUILDNUM=1) so '~' sorts
# pre-releases BELOW stable in dpkg/reprepro. Downloads use the tag
# (ANONYMIZER_BRANCH), not VERSION, so this never affects the source URL.
if command -v apt-get &>/dev/null; then
    if [[ "$ANONYMIZER_BUILDNUM" == *_* ]]; then
        ANONYMIZER_PRETAG="${ANONYMIZER_BUILDNUM%%_*}"
        export ANONYMIZER_VERSION="${ANONYMIZER_VERSION}~${ANONYMIZER_PRETAG}"
        ANONYMIZER_BUILDNUM="${ANONYMIZER_BUILDNUM#*_}"
    fi
fi

export REPO_TYPE="${REPO_TYPE:-daily}"

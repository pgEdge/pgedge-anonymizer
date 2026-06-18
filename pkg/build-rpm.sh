#!/bin/bash
set -euo pipefail

RHEL="$(rpm --eval %rhel)"
ARCH=$(uname -m)
if [ "$ARCH" = "aarch64" ]; then
  ARCH="arm64"
fi

# Release assets are named with the full tag version (e.g. 1.0.0-beta2);
# ANONYMIZER_VERSION is the suffix-stripped value the spec expects in SOURCES.
TAG_VERSION="${ANONYMIZER_BRANCH#v}"
ARTIFACT_DIR="${ARTIFACT_DIR:-$(pwd)/release-artifacts}"
RELEASE_URL="https://github.com/pgEdge/pgedge-anonymizer/releases/download/${ANONYMIZER_BRANCH}"
RAW_URL="https://raw.githubusercontent.com/pgEdge/pgedge-anonymizer/${ANONYMIZER_BRANCH}"

# stage <canonical-local-name> <release-asset-or-raw-name> <dest> <url-base>
# Prefer the workflow-staged file under release-artifacts/ (so simulate_tag
# branch tests work and package cells don't depend on the GH release); else
# download from the given URL base.
stage() {
  local local_name="$1" remote_name="$2" dest="$3" url_base="$4"
  if [ -f "${ARTIFACT_DIR}/${local_name}" ]; then
    cp "${ARTIFACT_DIR}/${local_name}" "${dest}"
  else
    wget -q "${url_base}/${remote_name}" -O "${dest}"
  fi
}

prepare() {
  setup_dnf_build_env

  echo "Copying packaging files..."
  cp "${COMPONENT_NAME}/rpm/anonymizer.spec" ~/rpmbuild/SPECS/

  echo "Staging source tarball + docs into SOURCES..."
  stage "anonymizer.tar.gz" "pgedge-anonymizer_${TAG_VERSION}_Linux_${ARCH}.tar.gz" \
        ~/rpmbuild/SOURCES/pgedge-anonymizer_${ANONYMIZER_VERSION}_Linux_${ARCH}.tar.gz "${RELEASE_URL}"
  # LICENCE.md isn't in the GoReleaser archive (its files glob is LICENSE*).
  stage "LICENCE.md" "LICENCE.md" ~/rpmbuild/SOURCES/LICENCE.md "${RAW_URL}"
  stage "pgedge-anonymizer-patterns.yaml" "pgedge-anonymizer-patterns.yaml" \
        ~/rpmbuild/SOURCES/pgedge-anonymizer-patterns.yaml "${RAW_URL}"
  # Packaged default config ships from this repo, not the release.
  cp "${COMPONENT_NAME}"/common/pgedge-anonymizer.yaml ~/rpmbuild/SOURCES/

  # This function is for debugging purpose if you have your own keys. GH workflow does not need it.
  #import_gpg_keys

  echo "🔧 Installing RPM build dependencies..."
  dnf builddep -y \
    --define "anonymizer_version ${ANONYMIZER_VERSION}" \
    --define "anonymizer_buildnum ${ANONYMIZER_BUILDNUM}" \
    --define "arch ${ARCH}" \
    ~/rpmbuild/SPECS/anonymizer.spec
}

build() {
  echo "Building RPM and SRPM..."
  QA_RPATHS=$(( 0xffff )) rpmbuild -ba ~/rpmbuild/SPECS/anonymizer.spec \
    --define "anonymizer_version ${ANONYMIZER_VERSION}" \
    --define "anonymizer_buildnum ${ANONYMIZER_BUILDNUM}" \
    --define "arch ${ARCH}"
}

post_build() {
  echo "📤 Copying built RPMs to /output..."
  mkdir -p /output
  cp -v ~/rpmbuild/RPMS/*/*.rpm /output/ || echo "No binary RPMs found"
  cp -v ~/rpmbuild/SRPMS/*.src.rpm /output/ || echo "No SRPM found"

  sign_rpms /output/*.rpm
  validate_signatures /output/*.rpm
}

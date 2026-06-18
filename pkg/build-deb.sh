#!/usr/bin/env bash
set -euo pipefail

# Environment variables
BUILD_DIR="/tmp/pg_deb_build"

CWD="$(pwd)"

export DEBIAN_FRONTEND=noninteractive
ARCH=$(uname -m)
if [ "$ARCH" = "aarch64" ]; then
  ARCH="arm64"
fi

# Release assets are named with the full tag version (e.g. 1.0.0-beta2).
# ANONYMIZER_VERSION may carry a '~beta…' suffix for DEB ordering, so it must
# NOT be used to build the download URL — use the tag instead.
TAG_VERSION="${ANONYMIZER_BRANCH#v}"
ARTIFACT_DIR="${ARTIFACT_DIR:-${CWD}/release-artifacts}"
SRC_DIR="${BUILD_DIR}/pgedge-anonymizer-${TAG_VERSION}"
RELEASE_URL="https://github.com/pgEdge/pgedge-anonymizer/releases/download/${ANONYMIZER_BRANCH}"

# stage <canonical-local-name> <remote-name> <dest> <url-base>
stage() {
  local local_name="$1" remote_name="$2" dest="$3" url_base="$4"
  if [ -f "${ARTIFACT_DIR}/${local_name}" ]; then
    cp "${ARTIFACT_DIR}/${local_name}" "${dest}"
  else
    wget -q "${url_base}/${remote_name}" -O "${dest}"
  fi
}

prepare() {

  setup_apt_build_env

  # This function is for debugging purpose if you have your own keys. GH workflow does not need it.
  #import_gpg_keys

  echo "Resetting build workspace at ${SRC_DIR}..."
  rm -rf "$SRC_DIR"
  mkdir -p "$SRC_DIR"

  echo "Staging source tarball..."
  stage "anonymizer.tar.gz" "pgedge-anonymizer_${TAG_VERSION}_Linux_${ARCH}.tar.gz" \
        "${BUILD_DIR}/anonymizer.tar.gz" "${RELEASE_URL}"
  tar -C "$SRC_DIR" -xzf "${BUILD_DIR}/anonymizer.tar.gz"

  echo "Moving Debian packaging into source directory..."
  cp -rp "${CWD}/${COMPONENT_NAME}/deb/debian" "$SRC_DIR/"
  cp "${CWD}/${COMPONENT_NAME}"/common/pgedge-anonymizer.yaml "$SRC_DIR/debian/"

  echo "Staging LICENCE.md + patterns from the repo checkout..."
  cp "${CWD}/LICENCE.md" "$SRC_DIR/LICENCE.md"
  cp "${CWD}/pgedge-anonymizer-patterns.yaml" "$SRC_DIR/debian/pgedge-anonymizer-patterns.yaml"

  echo "Installing build dependencies..."
  cd "$SRC_DIR"
  sudo apt-get update
  sudo apt-get build-dep -y .
}

build() {

  cd "$SRC_DIR"
  echo "Building Debian package..."
  DISTRO=$(lsb_release -cs)
  rm -f debian/changelog
cat > debian/changelog <<EOF
pgedge-anonymizer (${ANONYMIZER_VERSION}-${ANONYMIZER_BUILDNUM}.${DISTRO}) ${DISTRO}; urgency=medium

  * Update pgedge-anonymizer package.

 -- pgEdge Build Team <support@pgedge.com>  $(date -R)
EOF

  dpkg-buildpackage -us -uc -b
}

post_build() {
  echo "Copying .deb packages to output..."
  sudo mkdir -p "/output"
  rename_ddeb_packages $BUILD_DIR
  sudo cp "$BUILD_DIR"/*.deb "/output" || echo "No .deb packages found."
}

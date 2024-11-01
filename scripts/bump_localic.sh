#!/bin/sh
# sh scripts/bump_localic.sh

NEW_VERSION=v8.8.0

find . -type f -name "*.md"

findAndReplace() {
    find . -type f -name "$1" -not -path "*node_modules*" -exec sed -i "$2" {} \;
}

findAndReplace "*.yml" "s/v[0-9]\.[0-9]\.[0-9]\/local-ic/$NEW_VERSION\/local-ic/g"
findAndReplace "*.yml.optional" "s/v[0-9]\.[0-9]\.[0-9]\/local-ic/$NEW_VERSION\/local-ic/g"
findAndReplace "Dockerfile" "s/v[0-9]\.[0-9]\.[0-9]\/local-ic/$NEW_VERSION\/local-ic/g"

GIT_REPO="https://github.com/strangelove-ventures/interchaintest"
GIT_REPO=$(echo $GIT_REPO | sed 's/\//\\\//g') # https:\/\/github.com\/strangelove-ventures\/interchaintest
findAndReplace "Makefile" "s/v[0-9]\.[0-9]\.[0-9]\/local-ic/$NEW_VERSION\/local-ic/g"

# replace https://github.com/strangelove-ventures/interchaintest/releases/download/v8.7.0/local-ic in makefile
findAndReplace "Makefile" "s/https:\/\/github.com\/strangelove-ventures\/interchaintest\/releases\/download\/v[0-9]\.[0-9]\.[0-9]\/local-ic/https:\/\/github.com\/strangelove-ventures\/interchaintest\/releases\/download\/$NEW_VERSION\/local-ic/g"

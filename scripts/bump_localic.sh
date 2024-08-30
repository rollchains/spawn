#!/bin/sh
# sh scripts/bump_localic.sh

NEW_VERSION=v8.7.0

find . -type f -name "*.yml" -exec sed -i "s/v[0-9]\.[0-9]\.[0-9]\/local-ic/$NEW_VERSION\/local-ic/g" {} \;
find . -type f -name "*.yml.optional" -exec sed -i "s/v[0-9]\.[0-9]\.[0-9]\/local-ic/$NEW_VERSION\/local-ic/g" {} \;
find . -type f -name "Dockerfile" -exec sed -i "s/v[0-9]\.[0-9]\.[0-9]\/local-ic/$NEW_VERSION\/local-ic/g" {} \;

GIT_REPO="https://github.com/strangelove-ventures/interchaintest"
GIT_REPO=$(echo $GIT_REPO | sed 's/\//\\\//g') # https:\/\/github.com\/strangelove-ventures\/interchaintest
find . -type f -name "Makefile" -exec sed -i "s/v[0-9]\.[0-9]\.[0-9] $GIT_REPO/$NEW_VERSION $GIT_REPO/g" {} \;

# replace https://github.com/strangelove-ventures/interchaintest/releases/download/v8.7.0/local-ic in makefile
find . -type f -name "Makefile" -exec sed -i "s/https:\/\/github.com\/strangelove-ventures\/interchaintest\/releases\/download\/v[0-9]\.[0-9]\.[0-9]\/local-ic/https:\/\/github.com\/strangelove-ventures\/interchaintest\/releases\/download\/$NEW_VERSION\/local-ic/g" {} \;

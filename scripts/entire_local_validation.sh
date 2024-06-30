#!/bin/bash
set -m  ## set job control
#
# sh scripts/entire_local_validation.sh
# Original:
# This script takes what is normally done via github CI and runs it locally to ensure matrix variations work as expected
# This includes;
# - Docker
# - Unit Test
# - Make Commands
# - Linting
# - Integration Test
#
# Pre Checks:
# - reference docs/TAGGED_RELEASE_CHECKLIST.md
# - https://github.com/rollchains/spawn/issues/86
# - scripts/bump_localic.sh

# The new version which spawn will be
NEW_SPAWN_VERSION=v0.50.3
find . -type f -name "README.md" -exec sed -i "s/git checkout v0\.[0-9]*\.[0-9]*/git checkout $NEW_SPAWN_VERSION/g" {} \;

make install

LOGS_DIR=${HOME}/spawnlogs
mkdir -p $LOGS_DIR

# docker network prune ??

# TEst TODO:
  # create new module(s) + ibc middleware

  # setup_network builds, waits, and runs the integration tests for a given network
# this allows for multi-process execution in an enviroment via screen jobs
function setup_network() {
    name=$1

    # Builds docker image, pushes to background. Starts installing binary
    # Once completed, pops local-image back to foreground
    #
    # (cd $name; make local-image) & (cd $name; make install) && fg
    screen -S $name-build -dm bash -c "(cd $name; make local-image) && (cd $name; make install)"
    while screen -ls | grep -q $name-build; do
        echo "Waiting for $name-build to finish"
        sleep 3
    done

    ictests=$(find $name -type f -name "interchaintest-e2e.yml" -exec grep -o 'ictest-[a-zA-Z]*' {} \;) && echo $ictests

    for test in $ictests; do
        logFile=$logs/$test.log
        screenName=$name-$test

        echo "Running $screenName"
        echo $screenName > $logFile
        screen -S $screenName -dm bash -c "cd $name; make $test > $logFile 2>&1"
    done
}

# create a new function base
function base_standard_network() {
    name=spawntestbase

    logs=$LOGS_DIR/$name && mkdir -p $logs

    spawn new $name --bypass-prompt --bech32=abcd --bin=$(echo name)d --denom=uxyz --org=strangelove

    setup_network $name
}

function minimal_network() {
    name=spawntestminimal

    logs=$LOGS_DIR/$name && mkdir -p $logs

    # Copy all from spawn new `--disable` list
    spawn new $name --consensus=proof-of-stake --disable=tokenfactory,globalfee,ibc-packetforward,ibc-ratelimit,cosmwasm,wasm-light-client,optimistic-execution,ignite-cli --bech32=abcd --bin=$(echo $name)d  --denom=uxyz --org=strangelove

    setup_network $name
}



function validate_all_test() {
    find $LOGS_DIR -type f -name "*.log" -exec grep -l "FAIL" {} \;
    find $LOGS_DIR -type f -name "*.log" -exec grep -l "No rule to make target" {} \; # make file error
}

base_standard_network &
minimal_network &

# Wait for processes to complete
while lines=$(screen -ls | grep  "[0-9]\." | wc -l) && [ $lines -gt 0 ]; do
    echo "Waiting for $lines screens to finish"
    sleep 3
done

echo "All screens finished, validating outputs"
validate_all_test

killall screen || true

# Reminder here to rm go binaries named spawntest*d & docker imagesC
# rm -rf spawntest*/
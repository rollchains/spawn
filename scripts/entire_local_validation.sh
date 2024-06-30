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

# TEst TODO:
  # create new module(s) + ibc middleware

# create a new function base
function base_standard_network() {
    name=spawntestbase
    cd $name

    logs=$LOGS_DIR/$name
    mkdir -p $logs

    # spawn new $name --bypass-prompt --bech32=abcd --bin=$(name)d --denom=uxyz --org=strangelove

    # Builds docker image, pushes to background. Starts installing binary
    # Once completed, pops local-image back to foreground
    # make local-image & make install && fg

    ictests=$(find . -type f -name "interchaintest-e2e.yml" -exec grep -o 'ictest-[a-zA-Z]*' {} \;) && echo $ictests

    for test in $ictests; do
        logFile=$logs/$test.log
        echo "" > $logFile
        screenName=$name-$test
        echo "Running $screenName"
        screen -S $screenName -dm bash -c "make $test > $logFile 2>&1"
    done
}

function validate_all_test() {
    find $LOGS_DIR -type f -name "*.log" -exec grep -l "FAIL" {} \;
}

base_standard_network

# wait until there are no screens running
while screen -ls | grep -q "[0-9]\."; do
    echo "Waiting for screens to finish"
    sleep 5
done

echo "All screens finished, validating outputs"
validate_all_test

# rm go binaries named spawntest*d & docker images
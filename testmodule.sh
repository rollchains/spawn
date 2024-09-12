make install
# rm -rf myproject

make template-staking
cd myproject
spawn module new mynsibc --ibc-module --log-level=debug
make proto-gen

code .

# run testnet
make local-image
local-ic start self-ibc


# connect via relayer
source <(curl -s https://raw.githubusercontent.com/strangelove-ventures/interchaintest/main/local-interchain/bash/source.bash)
API_ADDR="http://localhost:8080"

ICT_POLL_FOR_START $API_ADDR 50

# only 1
CHANNELS=`ICT_RELAYER_CHANNELS $API_ADDR "localchain-1"` && echo "CHANNELS: $CHANNELS"

ICT_RELAYER_EXEC $API_ADDR "localchain-1" "rly tx connect localchain-1_localchain-2 --src-port=mynsibc --dst-port=mynsibc --order=unordered --version=mynsibc-1"

# 2 open :D
CHANNELS=`ICT_RELAYER_CHANNELS $API_ADDR "localchain-1"` && echo "CHANNELS: $CHANNELS"


# perform ibc action


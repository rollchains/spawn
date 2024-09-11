make install
# rm -rf myproject

make template-staking
cd myproject
spawn module new mynsibc --ibc-module --log-level=debug
spawn module new mynsibc2 --ibc-module --log-level=debug

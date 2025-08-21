#!/usr/bin/env bash
nohup anvil > anvil.log 2>&1 &
sleep 1
cd forge && forge script script/Counter.s.sol --rpc-url http://localhost:8545 --private-key 0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80  --broadcast
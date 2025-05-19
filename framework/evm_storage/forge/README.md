## Commands

```
anvil &
forge script script/Deploy.s.sol:Deploy --rpc-url http://127.0.0.1:8545 --broadcast

forge script script/Counter.s.sol --rpc-url http://localhost:8545 --private-key 0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80  --broadcast

cast storage 0x5FbDB2315678afecb367f032d93F642f64180aa3 0x1 --rpc-url http://localhost:8545

cast call 0x5FbDB2315678afecb367f032d93F642f64180aa3 \
 "number()(uint256)" \
 --rpc-url http://localhost:8545
 
 cast call 0x5FbDB2315678afecb367f032d93F642f64180aa3 \
 "values(uint256)(uint256)" 0 \
 --rpc-url http://localhost:8545
 
  cast call 0x5FbDB2315678afecb367f032d93F642f64180aa3 \
 "scores(address)(uint256)" 0x5FbDB2315678afecb367f032d93F642f64180aa3 \
 --rpc-url http://localhost:8545

cast send 0x5FbDB2315678afecb367f032d93F642f64180aa3 \
  "increment()" \
  --private-key 0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80 \
  --rpc-url http://localhost:8545
  
 cast --to-uint256 200 | cast --to-bytes32
 
 cast rpc anvil_setStorageAt \
  0x5FbDB2315678afecb367f032d93F642f64180aa3 \
  0x0 \
  0x000000000000000000000000000000000000000000000000000000000000002c

```
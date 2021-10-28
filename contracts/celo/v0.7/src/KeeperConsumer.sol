pragma solidity 0.7.6;

import "./KeeperCompatibleInterface.sol";


contract KeeperConsumer is KeeperCompatibleInterface {
    uint public counter;
    uint public immutable interval;
    uint public lastTimeStamp;


    constructor(uint updateInterval) public {
        interval = updateInterval;
        lastTimeStamp = block.timestamp;
        counter = 0;
    }

    function checkUpkeep(bytes calldata checkData) external override returns (bool upkeepNeeded, bytes memory performData) {
        return (true, checkData);
    }

    function performUpkeep(bytes calldata performData) external override {
        counter = counter + 1;
    }
}


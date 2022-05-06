pragma solidity ^0.7.0;

import "./KeeperCompatibleInterface.sol";
import "./KeeperBase.sol";

contract KeeperConsumer is KeeperCompatibleInterface, KeeperBase {
    uint public counter;
    uint public immutable interval;
    uint public lastTimeStamp;
    uint[] array = [2,4,5];


    constructor(uint updateInterval) public {
        interval = updateInterval;
        lastTimeStamp = block.timestamp;
        counter = 0;
    }

    function checkUpkeep(bytes calldata checkData) 
    external 
    override
    cannotExecute
    returns (bool upkeepNeeded, bytes memory performData) {
        if (array[counter]%2==0) {
            return (true, checkData);
        }else{
            counter = 0;
            return (false, checkData);
        }
    }

    function performUpkeep(bytes calldata performData) external override {
        counter = counter + 1;
    }
}


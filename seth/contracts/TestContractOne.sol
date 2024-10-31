// SPDX-License-Identifier: MIT
pragma solidity 0.8.19;

import "./AbstractContractWithEvent.sol";

contract unique_event_two is AbstractContractWithEvent {
    function executeFirstOperation(int256 x, int256 y) public returns (int256) {
        emit NonUniqueEvent(x, y);
        return x + y;
    }
}

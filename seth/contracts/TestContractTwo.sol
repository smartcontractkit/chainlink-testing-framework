// SPDX-License-Identifier: MIT
pragma solidity 0.8.19;

import "./AbstractContractWithEvent.sol";

contract TestContractTwo is AbstractContractWithEvent {
    function executeSecondOperation(int256 x, int256 y) public returns (int256) {
        emit NonUniqueEvent(x, y);
        return x + y;
    }
}

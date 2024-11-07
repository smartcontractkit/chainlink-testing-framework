// SPDX-License-Identifier: MIT
pragma solidity 0.8.19;

interface DebugContractCallback {
    function callbackMethod(int x) external returns (int);
}

contract NetworkDebugSubContract {
    int256 storedData;

    /*
        Basic types events
        1 index by default is selector signature, topic 0
    */

    event NoIndexEventString(string str);
    event NoIndexEvent(address sender);
    event OneIndexEvent(uint indexed a);
    event TwoIndexEvent(uint256 indexed roundId, address indexed startedBy);
    event ThreeIndexEvent(uint256 indexed roundId, address indexed startedBy, uint256 startedAt);
    event UniqueSubDebugEvent();

    /* Struct events */

    struct Account {
        string name;
        uint64 balance;
        uint dailyLimit;
    }
    event NoIndexStructEvent(Account a);

    /* Errors */

    error CustomErr(uint256 available, uint256 required);

    function trace(int256 x, int256 y) public returns (int256) {
        y = y + 2;
        emit TwoIndexEvent(uint256(y), address(msg.sender));
        return x + y;
    }

    function traceWithCallback(int256 x, int256 y) public returns (int256) {
        emit TwoIndexEvent(uint256(y), address(msg.sender));
        int256 response = DebugContractCallback(msg.sender).callbackMethod(y);
        emit OneIndexEvent(uint(response));
        return response;
    }

    function traceOneInt(int256 x) public returns (int256 r) {
        emit NoIndexEvent(msg.sender);
        return x + 3;
    }

    function traceUniqueEvent() public {
        emit UniqueSubDebugEvent();
    }

    function alwaysRevertsCustomError(uint256 x, uint256 y) public {
        revert CustomErr({
            available: x,
            required: y
        });
    }

    function pay() public payable {}
}

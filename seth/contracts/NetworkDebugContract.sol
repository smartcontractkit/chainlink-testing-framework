// SPDX-License-Identifier: MIT
pragma solidity 0.8.19;

import "./NetworkDebugSubContract.sol";

contract NetworkDebugContract {
    int256 public storedData;

    mapping(address => int256) public storedDataMap;
    mapping(int256 => int256) public counterMap;

    NetworkDebugSubContract public subContract;

    uint256 private data;

    constructor(address subAddr) {
        subContract = NetworkDebugSubContract(subAddr);
        data = 256;
    }

    /*
        Basic types events
        1 index by default is selector signature, topic 0
    */

    event NoIndexEventString(string str);
    event NoIndexEvent(address sender);
    event OneIndexEvent(uint indexed a);
    event IsValidEvent(bool success);
    event TwoIndexEvent(uint256 indexed roundId, address indexed startedBy);
    event ThreeIndexEvent(uint256 indexed roundId, address indexed startedBy, uint256 indexed startedAt);
    event ThreeIndexAndOneNonIndexedEvent(uint256 indexed roundId, address indexed startedBy, uint256 indexed startedAt, string dataId);
    event CallbackEvent(int256 indexed a);

    /* Struct events */

    struct Account {
        string name;
        uint64 balance;
        uint dailyLimit;
    }

    event NoIndexStructEvent(Account a);

    /* Errors */

    error CustomErr(uint256 available, uint256 required);
    error CustomErrNoValues();
    error CustomErrWithMessage(string message);

    /* Getters/Setters */
    function setMap(int256 x) public returns (int256 value) {
        storedDataMap[msg.sender] = x;
        return x;
    }

    function getMap() public view returns (int256 data) {
        return storedDataMap[msg.sender];
    }

    function set(int256 x) public returns (int256 value) {
        storedData = x;
        return x;
    }

    function addCounter(int256 idx, int256 x) public returns (int256 value) {
        counterMap[idx] += x;
        return x;
    }

    function getCounter(int256 idx) public view returns (int256 data) {
        return counterMap[idx];
    }

    function resetCounter(int256 idx) public {
        counterMap[idx] = 0;
    }

    function get() public view returns (int256 data) {
        return storedData;
    }

    function trace(int256 x, int256 y) public returns (int256) {
        subContract.trace(x, y);
        emit TwoIndexEvent(uint256(y), address(msg.sender));
        return x + y;
    }

    function traceDifferent(int x, int256 y) public returns (int256) {
        subContract.traceOneInt(y);
        emit OneIndexEvent(uint(x));
        return x + y;
    }

    function validate(int x, int256 y) public returns (bool) {
        emit IsValidEvent(x > y);
        return x > y;
    }

    function traceWithValidate(int x, int256 y) public payable returns (int256) {
        if (validate(x, y)) {
            subContract.trace(x, y);
            emit TwoIndexEvent(uint256(y), address(msg.sender));
            return x + y;
        }

        revert CustomErrWithMessage("first int was not greater than second int");
    }

    function traceYetDifferent(int x, int256 y) public returns (int256) {
        subContract.trace(x, y);
        emit TwoIndexEvent(uint256(y), address(msg.sender));
        return x + y;
    }

    /* Events */

    function emitNoIndexEventString() public {
        emit NoIndexEventString("myString");
    }

    function emitNoIndexEvent() public {
        emit NoIndexEvent(msg.sender);
    }

    function emitOneIndexEvent() public {
        emit OneIndexEvent(83);
    }

    function emitTwoIndexEvent() public {
        emit TwoIndexEvent(1, address(msg.sender));
    }

    function emitThreeIndexEvent() public {
        emit ThreeIndexEvent(1, address(msg.sender), 3);
    }

    function emitFourParamMixedEvent() public {
        emit ThreeIndexAndOneNonIndexedEvent(2, address(msg.sender), 3, "some id");
    }

    function emitNoIndexStructEvent() public {
        emit NoIndexStructEvent(Account("John", 5, 10));
    }

    /* Reverts */

    function alwaysRevertsRequire() public {
        require(false, "always revert error");
    }

    function alwaysRevertsAssert() public {
        assert(false);
    }

    function alwaysRevertsCustomError() public {
        revert CustomErr({
            available: 12,
            required: 21
        });
    }

    function alwaysRevertsCustomErrorNoValues() public {
        revert CustomErrNoValues();
    }

    /* Inputs/Outputs */

    function emitNamedInputsOutputs(uint256 inputVal1, string memory inputVal2) public returns (uint256 outputVal1, string memory outputVal2) {
        return (inputVal1, inputVal2);
    }

    function emitInputsOutputs(uint256 inputVal1, string memory inputVal2) public returns (uint256, string memory) {
        return (inputVal1, inputVal2);
    }

    function emitInputs(uint256 inputVal1, string memory inputVal2) public {
        return;
    }

    function emitOutputs() public returns (uint256, string memory) {
        return (31337, "outputVal1");
    }

    function emitNamedOutputs() public returns (uint256 outputVal1, string memory outputVal2) {
        return (31337, "outputVal1");
    }

    function emitInts(int first, int128 second, uint third) public returns (int, int128 outputVal1, uint outputVal2) {
        return (first, second, third);
    }

    function emitAddress(address addr) public returns (address) {
        return addr;
    }

    function emitBytes32(bytes32 input) public returns (bytes32 output) {
        return input;
    }

    function processUintArray(uint256[] memory input) public returns (uint256[] memory) {
        uint256[] memory output = new uint256[](input.length);
        for (uint i = 0; i < input.length; i++) {
            output[i] = input[i] + 1;
        }
        return output;
    }

    function processAddressArray(address[] memory input) public returns (address[] memory) {
        return input;
    }

    // struct with dynamic fields
    struct Data {
        string name;
        uint256[] values;
    }

    function processDynamicData(Data calldata data) public returns (Data calldata) {
        return data;
    }

    function processFixedDataArray(Data[3] calldata data) public returns (Data[2] memory) {
        Data[2] memory output;

        output[0] = data[0];
        output[1] = data[1];

        return output;
    }

    struct NestedData {
        Data data;
        bytes dynamicBytes;
    }

    function processNestedData(NestedData calldata data) public returns (NestedData memory) {
        return data;
    }

    /* Overload of processNestedData */
    function processNestedData(Data calldata data) public returns (NestedData memory) {
        bytes32 hashedData = keccak256(abi.encodePacked(data.name));

        bytes memory convertedData = new bytes(32);
        for (uint i = 0; i < 32; i++) {
            convertedData[i] = hashedData[i];
        }

        return NestedData(data, convertedData);
    }

    function pay() public payable {}

    /* Fallback and receive */

    event EtherReceived(address sender, uint amount);

    fallback() external payable {
        emit EtherReceived(msg.sender, msg.value);
    }

    event Received(address caller, uint amount, string message);

    receive() external payable {
        emit Received(msg.sender, msg.value, "Received Ether");
    }

    /* Enums */

    enum Status { Pending, Active, Completed, Cancelled }
    Status public currentStatus;

    event CurrentStatus(Status indexed status);

    function setStatus(Status status) public returns (Status) {
        currentStatus = status;
        emit CurrentStatus(currentStatus);
        return status;
    }

    function traceSubWithCallback(int256 x, int256 y) public returns (int256) {
        y = y + 2;
        subContract.traceWithCallback(x, y);
        emit TwoIndexEvent(1, address(msg.sender));
        return x + y;
    }

    function callRevertFunctionInSubContract(uint256 x, uint256 y) public {
        subContract.alwaysRevertsCustomError(x, y);
    }


    function callRevertFunctionInTheContract() public {
        alwaysRevertsCustomError();
    }

    /* Static call */
    function getData() external view returns (uint256) {
        return data;
    }

    function performStaticCall() external view returns (uint256) {
        address self = address(this);

        // Perform a static call to getData function
        (bool success, bytes memory returnData) = self.staticcall(
            abi.encodeWithSelector(this.getData.selector)
        );

        require(success, "Static call failed");

        uint256 result = abi.decode(returnData, (uint256));

        return result;
    }

    /* Callback function */
    function callbackMethod(int x) external returns (int) {
        emit CallbackEvent(x);
        return x;
    }

    /* ERC677 token transfer */

    event CallDataLength(uint256 length);

    function onTokenTransfer(address sender, uint256 amount, bytes calldata data) external {
        emit CallDataLength(data.length);

        if (data.length == 0) {
            revert CustomErr({
                available: 99,
                required: 101
            });
        }
        (bool success, bytes memory returnData) = address(this).delegatecall(data);

        if (!success) {
            if (returnData.length > 0) {
                assembly {
                    let returndata_size := mload(returnData)
                    revert(add(32, returnData), returndata_size)
                }
            } else {
                revert CustomErrWithMessage("delegatecall failed with no reason");
            }
        }
        this.performStaticCall();

        bytes4 selector = bytes4(data[:4]);
        if (selector == this.traceYetDifferent.selector) {
            revert CustomErrWithMessage("oh oh oh it's magic!");
        }

        this.traceSubWithCallback(1, 2);
    }
}

// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

contract Counter {
    uint256 public number;
    uint256[] public values;
    mapping(address => uint256) public scores;

    // private array and map of structs
    struct Signer {
        address addr;
        uint8 index;
        uint8 group;
    }
    mapping(address => Signer) s_signers;
    Signer[] a_signers;

    constructor() public {
        number = 1;
        values = [1, 2, 3];
        scores[address(0x5FbDB2315678afecb367f032d93F642f64180aa3)] = 1;

        for (uint8 i = 0; i < 5; i++) {
            address signerAddr = address(uint160(uint256(keccak256(abi.encodePacked(i)))));
            Signer memory signer = Signer({
                addr: signerAddr,
                index: i,
                group: i + 10
            });

            a_signers.push(signer);
            s_signers[signerAddr] = signer;
        }
    }

    function setNumber(uint256 newNumber) public {
        number = newNumber;
    }

    function increment() public {
        number++;
    }

    function pushValue(uint256 value) public {
        values.push(value);
    }

    function setScore(address who, uint256 score) public {
        scores[who] = score;
    }

    // this function may not be present in real contracts but we need them for tests to verify mutated data
    function getASigner(uint i) external view returns (address, uint8, uint8) {
        Signer memory s = a_signers[i];
        return (s.addr, s.index, s.group);
    }

    function getSSigner(address who) external view returns (address, uint8, uint8) {
        Signer memory s = s_signers[who];
        return (s.addr, s.index, s.group);
    }
}

// SPDX-License-Identifier: MIT
pragma solidity 0.8.19;

abstract contract AbstractContractWithEvent {
    event NonUniqueEvent(int256 indexed a, int256 indexed b);
}

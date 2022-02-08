pragma solidity ^0.7.6;

contract KeeperConsumerPerformance {
    event PerformingUpkeep (
        bool eligible,
        address from,
        uint256 initialCall,
        uint256 nextEligible,
        uint256 blockNumber
    );

  uint256 public initialCall = 0;
  uint256 public nextEligible = 0;
  uint256 public testRange;
  uint256 public averageEligibilityCadence;
  uint256 count = 0;
  
  constructor(uint256 _testRange, uint256 _averageEligibilityCadence) {
    testRange = _testRange;
    averageEligibilityCadence = _averageEligibilityCadence;
  }

  function checkUpkeep(bytes calldata data) external returns (bool, bytes memory) {
    return (eligible(), bytes(""));
  }

  function performUpkeep(bytes calldata data) external {
    bool eligible = eligible();
    uint256 blockNum = block.number;
    emit PerformingUpkeep(eligible, tx.origin, initialCall, nextEligible, blockNum);
    require(eligible);
    if (initialCall == 0) {
      initialCall = blockNum;
    }
    nextEligible = blockNum + averageEligibilityCadence; // Switched this from a rand check to more predictable behavior 
    count++;
  }

  function getCountPerforms() view public returns(uint256) {
    return count;
  }

  function eligible() view internal returns(bool) {
    return initialCall == 0 ||
      (
        block.number - initialCall < testRange &&
        block.number >= nextEligible
      );
  }

  function checkEligible() view public returns(bool) {
    return eligible();
  }

  function reset() external {
      initialCall = 0;
      count = 0;
  }

  function setSpread(uint _newTestRange, uint _newAverageEligibilityCadence) external {
    testRange = _newTestRange;
    averageEligibilityCadence = _newAverageEligibilityCadence;
  }
}
import "../../v0.6/src/interfaces/AggregatorV3Interface.sol";

contract MockETHLINKAggregator is AggregatorV3Interface {
    int256 public answer;
    constructor (int256 _answer) {
        answer = _answer;
    }
    function decimals() external override view returns (uint8) {
        return 18;
    }
    function description() external override view returns (string memory) {
        return "MockETHLINKAggregator";
    }
    function version() external override view returns (uint256) {
        return 1;
    }
    function getRoundData(uint80 _roundId) external override view returns (
        uint80 roundId,
        int256 answer,
        uint256 startedAt,
        uint256 updatedAt,
        uint80 answeredInRound
    ) {
       return (1, answer, block.timestamp, block.timestamp, 1);
    }
    function latestRoundData() external override view returns (
        uint80 roundId,
        int256 answer,
        uint256 startedAt,
        uint256 updatedAt,
        uint80 answeredInRound
    ) {
        return (1, answer, block.timestamp, block.timestamp, 1);
    }
}
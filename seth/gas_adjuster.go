package seth

import (
	"context"
	"fmt"
	"math"
	"math/big"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/avast/retry-go"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
)

const (
	Priority_Degen    = "degen" //this is undocumented option, which we left for cases, when we need to set the highest gas price
	Priority_Fast     = "fast"
	Priority_Standard = "standard"
	Priority_Slow     = "slow"

	Congestion_Low      = "low"
	Congestion_Medium   = "medium"
	Congestion_High     = "high"
	Congestion_VeryHigh = "extreme"
)

const (
	// each block has the same weight in the computation
	CongestionStrategy_Simple = "simple"
	// newer blocks have more weight in the computation
	CongestionStrategy_NewestFirst = "newest_first"
)

var (
	ZeroGasSuggestedErr = "either base fee or suggested tip is 0"
	BlockFetchingErr    = "failed to fetch enough block headers for congestion calculation"
)

// CalculateNetworkCongestionMetric calculates a simple congestion metric based on the last N blocks
// according to selected strategy.
func (m *Client) CalculateNetworkCongestionMetric(blocksNumber uint64, strategy string) (float64, error) {
	if m.HeaderCache == nil {
		return 0, fmt.Errorf("header cache is nil")
	}
	var getHeaderData = func(bn *big.Int) (*types.Header, error) {
		if bn == nil {
			return nil, fmt.Errorf("block number is nil")
		}
		cachedHeader, ok := m.HeaderCache.Get(bn.Int64())
		if ok {
			return cachedHeader, nil
		}

		timeout := blocksNumber / 100
		if timeout < 3 {
			timeout = 3
		} else if timeout > 6 {
			timeout = 6
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(mustSafeInt64(timeout))*time.Second)
		defer cancel()
		header, err := m.Client.HeaderByNumber(ctx, bn)
		if err != nil {
			return nil, err
		}
		// ignore the error here as at this point it is very improbable that block is nil and there's no error
		_ = m.HeaderCache.Set(header)
		return header, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(2*time.Second))
	defer cancel()
	lastBlockNumber, err := m.Client.BlockNumber(ctx)
	if err != nil {
		return 0, err
	}

	L.Trace().Msgf("Block range for gas calculation: %d - %d", lastBlockNumber-blocksNumber, lastBlockNumber)

	lastBlock, err := getHeaderData(big.NewInt(mustSafeInt64(lastBlockNumber)))
	if err != nil {
		return 0, err
	}

	var headers []*types.Header
	headers = append(headers, lastBlock)

	var wg sync.WaitGroup
	dataCh := make(chan *types.Header)

	go func() {
		for header := range dataCh {
			headers = append(headers, header)
			// placed here, because we want to wait for all headers to be received and added to slice before continuing
			wg.Done()
		}
	}()

	startTime := time.Now()
	for i := lastBlockNumber; i > lastBlockNumber-blocksNumber; i-- {
		// better safe than sorry (might happen for brand-new chains)
		if i <= 1 {
			break
		}

		wg.Add(1)
		go func(bn *big.Int) {
			header, err := getHeaderData(bn)
			if err != nil {
				L.Debug().Msgf("Failed to get block %d header due to: %s", bn.Int64(), err.Error())
				wg.Done()
				return
			}
			dataCh <- header
		}(big.NewInt(mustSafeInt64(i)))
	}

	wg.Wait()
	close(dataCh)

	endTime := time.Now()
	L.Debug().Msgf("Time to fetch %d block headers: %v", blocksNumber, endTime.Sub(startTime))

	minBlockCount := int(float64(blocksNumber) * 0.8)
	if len(headers) < minBlockCount {
		return 0, fmt.Errorf("%s. Wanted at least %d, got %d", BlockFetchingErr, minBlockCount, len(headers))
	}

	switch strategy {
	case CongestionStrategy_Simple:
		return calculateSimpleNetworkCongestionMetric(headers), nil
	case CongestionStrategy_NewestFirst:
		return calculateNewestFirstNetworkCongestionMetric(headers), nil
	default:
		return 0, fmt.Errorf("unknown congestion strategy: %s", strategy)
	}
}

// average gas used ratio for a basic congestion metric
func calculateSimpleNetworkCongestionMetric(headers []*types.Header) float64 {
	return calculateGasUsedRatio(headers)
}

// calculates a congestion metric using a logarithmic function that gives more weight to most recent block headers
func calculateNewestFirstNetworkCongestionMetric(headers []*types.Header) float64 {
	// sort blocks so that we are sure they are in ascending order
	slices.SortFunc(headers, func(i, j *types.Header) int {
		if i.Number.Uint64() < j.Number.Uint64() {
			return -1
		} else if i.Number.Uint64() > j.Number.Uint64() {
			return 1
		}
		return 0
	})

	var weightedSum, totalWeight float64
	// Determines how quickly the weight decreases. The lower the number, the higher the weight of newer blocks.
	scaleFactor := 10.0

	// Calculate weights starting from the older to most recent block header.
	for i, header := range headers {
		congestion := float64(header.GasUsed) / float64(header.GasLimit)

		distance := float64(len(headers) - 1 - i)
		weight := 1.0 / math.Log10(distance+scaleFactor)

		weightedSum += congestion * weight
		totalWeight += weight
	}

	if totalWeight == 0 {
		return 0
	}
	return weightedSum / totalWeight
}

// GetSuggestedEIP1559Fees returns suggested tip/fee cap calculated based on historical data, current congestion, and priority.
func (m *Client) GetSuggestedEIP1559Fees(ctx context.Context, priority string) (maxFeeCap *big.Int, adjustedTipCap *big.Int, err error) {
	L.Info().Msg("Calculating suggested EIP-1559 fees")
	var suggestedGasTip *big.Int
	var baseFee64, historicalSuggestedTip64 float64
	attempts := getSafeGasEstimationsAttemptCount(m.Cfg)

	retryErr := retry.Do(func() error {
		var tipErr error
		suggestedGasTip, tipErr = m.Client.SuggestGasTipCap(ctx)
		if tipErr != nil {
			return tipErr
		}

		L.Debug().
			Str("CurrentGasTip", fmt.Sprintf("%s wei / %s ether", suggestedGasTip.String(), WeiToEther(suggestedGasTip).Text('f', -1))).
			Msg("Current suggested gas tip")

		// Fetch the baseline historical base fee and tip for the selected priority
		var historyErr error
		//nolint
		baseFee64, historicalSuggestedTip64, historyErr = m.HistoricalFeeData(priority)
		return historyErr
	},
		retry.Attempts(attempts),
		retry.Delay(1*time.Second),
		retry.LastErrorOnly(true),
		retry.DelayType(retry.FixedDelay),
		retry.OnRetry(func(i uint, retryErr error) {
			L.Debug().
				Msgf("Retrying fetching of EIP1559 suggested fees due to: %s. Attempt %d/%d", retryErr.Error(), (i + 1), attempts)
		}))

	if retryErr != nil {
		err = retryErr
		return
	}

	L.Debug().
		Str("HistoricalBaseFee", fmt.Sprintf("%.0f wei / %s ether", baseFee64, WeiToEther(big.NewInt(int64(baseFee64))).Text('f', -1))).
		Str("HistoricalSuggestedTip", fmt.Sprintf("%.0f wei / %s ether", historicalSuggestedTip64, WeiToEther(big.NewInt(int64(historicalSuggestedTip64))).Text('f', -1))).
		Str("Priority", priority).
		Msg("Historical fee data")

	_, tipMagnitudeDiffText := calculateMagnitudeDifference(big.NewFloat(historicalSuggestedTip64), new(big.Float).SetInt(suggestedGasTip))

	L.Debug().
		Msgf("Historical tip is %s than suggested tip", tipMagnitudeDiffText)

	currentGasTip := suggestedGasTip
	if big.NewInt(int64(historicalSuggestedTip64)).Cmp(currentGasTip) > 0 {
		L.Debug().Msg("Historical suggested tip is higher than current suggested tip. Will use it instead.")
		currentGasTip = big.NewInt(int64(historicalSuggestedTip64))
	} else {
		L.Debug().Msg("Suggested tip is higher than historical tip. Will use suggested tip.")
	}

	if m.Cfg.IsExperimentEnabled(Experiment_Eip1559FeeEqualier) {
		L.Debug().Msg("FeeEqualier experiment is enabled. Will adjust base fee and tip to be of the same order of magnitude.")
		baseFeeTipMagnitudeDiff, _ := calculateMagnitudeDifference(big.NewFloat(baseFee64), new(big.Float).SetInt(currentGasTip))

		//one of values is 0, infinite order of magnitude smaller or larger
		if baseFeeTipMagnitudeDiff == -0 {
			if baseFee64 == 0.0 {
				L.Debug().Msg("Historical base fee is 0.0. Will use suggested tip as base fee.")
				baseFee64 = float64(currentGasTip.Int64())
			} else {
				L.Debug().Msg("Suggested tip is 0.0. Will use historical base fee as tip.")
				currentGasTip = big.NewInt(int64(baseFee64))
			}
		} else if baseFeeTipMagnitudeDiff < 3 {
			L.Debug().Msg("Historical base fee is 3 orders of magnitude lower than suggested tip. Will use suggested tip as base fee.")
			baseFee64 = float64(currentGasTip.Int64())
		} else if baseFeeTipMagnitudeDiff > 3 {
			L.Debug().Msg("Suggested tip is 3 orders of magnitude lower than historical base fee. Will use historical base fee as tip.")
			currentGasTip = big.NewInt(int64(baseFee64))
		}
	}

	if baseFee64 == 0.0 {
		L.Debug().
			Float64("BaseFee", baseFee64).
			Int64("SuggestedTip", currentGasTip.Int64()).
			Msgf("Incorrect gas data received from node: historical base fee was 0. Skipping automation gas estimation")
		return
	}

	if currentGasTip.Int64() == 0 {
		L.Debug().
			Int64("SuggestedTip", currentGasTip.Int64()).
			Str("Fallback gas tip", fmt.Sprintf("%d wei / %s ether", m.Cfg.Network.GasTipCap, WeiToEther(big.NewInt(m.Cfg.Network.GasTipCap)).Text('f', -1))).
			Msg("Suggested tip is 0.0. Although not strictly incorrect, it is unusual. Will use fallback value instead")

		currentGasTip = big.NewInt(m.Cfg.Network.GasTipCap)
	}

	// between 0.8 and 1.5
	var adjustmentFactor float64
	adjustmentFactor, err = getAdjustmentFactor(priority)
	if err != nil {
		return
	}

	// Calculate adjusted tip based on priority
	adjustedTipCapFloat := new(big.Float).Mul(big.NewFloat(adjustmentFactor), new(big.Float).SetFloat64(float64(currentGasTip.Int64())))
	adjustedTipCap, _ = adjustedTipCapFloat.Int(nil)

	adjustedBaseFeeFloat := new(big.Float).Mul(big.NewFloat(adjustmentFactor), new(big.Float).SetFloat64(baseFee64))
	adjustedBaseFee, _ := adjustedBaseFeeFloat.Int(nil)

	initialFeeCap := new(big.Int).Add(big.NewInt(int64(baseFee64)), currentGasTip)

	// skip if we do not want to calculate congestion metrics
	if m.Cfg.Network.GasPriceEstimationBlocks > 0 {
		// between 0 and 1 (empty blocks - full blocks)
		var congestionMetric float64
		//nolint
		congestionMetric, err = m.CalculateNetworkCongestionMetric(m.Cfg.Network.GasPriceEstimationBlocks, CongestionStrategy_NewestFirst)
		if err == nil {
			congestionClassification := classifyCongestion(congestionMetric)

			L.Debug().
				Str("CongestionMetric", fmt.Sprintf("%.4f", congestionMetric)).
				Str("CongestionClassification", congestionClassification).
				Float64("AdjustmentFactor", adjustmentFactor).
				Str("Priority", priority).
				Msg("Adjustment factors")

			// between 1.1 and 1.4
			var bufferAdjustment float64
			bufferAdjustment, err = getCongestionFactor(congestionClassification)
			if err != nil {
				return
			}

			// Calculate base fee buffer
			bufferedBaseFeeFloat := new(big.Float).Mul(new(big.Float).SetInt(adjustedBaseFee), big.NewFloat(bufferAdjustment))
			adjustedBaseFee, _ = bufferedBaseFeeFloat.Int(nil)

			// Apply buffer also to the tip
			bufferedTipCapFloat := new(big.Float).Mul(new(big.Float).SetInt(adjustedTipCap), big.NewFloat(bufferAdjustment))
			adjustedTipCap, _ = bufferedTipCapFloat.Int(nil)
		} else if !strings.Contains(err.Error(), BlockFetchingErr) {
			return
		} else {
			L.Debug().
				Msgf("Failed to calculate congestion metric due to: %s. Skipping congestion buffer adjustment", err.Error())

			// set error to nil, as we can still calculate the fees, but without congestion buffer
			// we don't want to return an error in this case
			err = nil
		}
	}

	maxFeeCap = new(big.Int).Add(adjustedBaseFee, adjustedTipCap)

	baseFeeDiff := big.NewInt(0).Sub(adjustedBaseFee, big.NewInt(int64(baseFee64)))
	gasTipDiff := big.NewInt(0).Sub(adjustedTipCap, currentGasTip)
	gasCapDiff := big.NewInt(0).Sub(maxFeeCap, initialFeeCap)

	L.Debug().
		Str("Diff (Wei/Ether)", fmt.Sprintf("%s wei / %s ether", gasTipDiff.String(), WeiToEther(gasTipDiff).Text('f', -1))).
		Str("Initial Tip", fmt.Sprintf("%s wei / %s ether", currentGasTip.String(), WeiToEther(currentGasTip).Text('f', -1))).
		Str("Final Tip", fmt.Sprintf("%s wei / %s ether", adjustedTipCap.String(), WeiToEther(adjustedTipCap).Text('f', -1))).
		Msg("Tip adjustment")

	L.Debug().
		Str("Diff (Wei/Ether)", fmt.Sprintf("%s wei / %s ether", baseFeeDiff.String(), WeiToEther(baseFeeDiff).Text('f', -1))).
		Str("Initial Base Fee", fmt.Sprintf("%s wei / %s ether", big.NewInt(int64(baseFee64)).String(), WeiToEther(big.NewInt(int64(baseFee64))).Text('f', -1))).
		Str("Final Base Fee", fmt.Sprintf("%s wei / %s ether", adjustedBaseFee.String(), WeiToEther(adjustedBaseFee).Text('f', -1))).
		Msg("Base Fee adjustment")

	L.Debug().
		Str("Diff (Wei/Ether)", fmt.Sprintf("%s wei / %s ether", gasCapDiff.String(), WeiToEther(gasCapDiff).Text('f', -1))).
		Str("Initial Fee Cap", fmt.Sprintf("%s wei / %s ether", initialFeeCap.String(), WeiToEther(initialFeeCap).Text('f', -1))).
		Str("Final Fee Cap", fmt.Sprintf("%s wei / %s ether", maxFeeCap.String(), WeiToEther(maxFeeCap).Text('f', -1))).
		Msg("Fee Cap adjustment")

	L.Info().
		Str("GasTipCap", fmt.Sprintf("%s wei / %s ether", adjustedTipCap.String(), WeiToEther(adjustedTipCap).Text('f', -1))).
		Str("GasFeeCap", fmt.Sprintf("%s wei / %s ether", maxFeeCap.String(), WeiToEther(maxFeeCap).Text('f', -1))).
		Msg("Calculated suggested EIP-1559 fees")

	return
}

// GetSuggestedLegacyFees calculates the suggested gas price based on historical data, current congestion, and priority.
func (m *Client) GetSuggestedLegacyFees(ctx context.Context, priority string) (adjustedGasPrice *big.Int, err error) {
	L.Info().
		Msg("Calculating suggested Legacy fees")

	var suggestedGasPrice *big.Int
	attempts := getSafeGasEstimationsAttemptCount(m.Cfg)

	retryErr := retry.Do(func() error {
		var priceErr error
		suggestedGasPrice, priceErr = m.Client.SuggestGasPrice(ctx)
		if priceErr != nil {
			return priceErr
		}

		if suggestedGasPrice.Int64() == 0 {
			return errors.New("suggested gas price is 0")
		}

		return nil
	},
		retry.Attempts(attempts),
		retry.Delay(1*time.Second),
		retry.LastErrorOnly(true),
		retry.DelayType(retry.FixedDelay),
		retry.OnRetry(func(i uint, retryErr error) {
			L.Debug().
				Msgf("Retrying fetching of legacy suggested gas price due to: %s. Attempt %d/%d", retryErr.Error(), (i + 1), attempts)
		}))

	if retryErr != nil {
		err = retryErr
		return
	}

	var adjustmentFactor float64
	adjustmentFactor, err = getAdjustmentFactor(priority)
	if err != nil {
		return
	}

	// Calculate adjusted tip based on congestion and priority
	adjustedGasPriceFloat := new(big.Float).Mul(big.NewFloat(adjustmentFactor), new(big.Float).SetFloat64(float64(suggestedGasPrice.Int64())))
	adjustedGasPrice, _ = adjustedGasPriceFloat.Int(nil)

	// skip if we do not want to calculate congestion metrics
	if m.Cfg.Network.GasPriceEstimationBlocks > 0 {
		// between 0 and 1 (empty blocks - full blocks)
		var congestionMetric float64
		//nolint
		congestionMetric, err = m.CalculateNetworkCongestionMetric(m.Cfg.Network.GasPriceEstimationBlocks, CongestionStrategy_NewestFirst)
		if err == nil {
			congestionClassification := classifyCongestion(congestionMetric)

			L.Debug().
				Str("CongestionMetric", fmt.Sprintf("%.4f", congestionMetric)).
				Str("CongestionClassification", congestionClassification).
				Float64("AdjustmentFactor", adjustmentFactor).
				Str("Priority", priority).
				Msg("Suggested Legacy fees")

			// between 1.1 and 1.4
			var bufferAdjustment float64
			bufferAdjustment, err = getCongestionFactor(congestionClassification)
			if err != nil {
				return
			}

			// Calculate and apply the buffer.
			bufferedGasPriceFloat := new(big.Float).Mul(new(big.Float).SetInt(adjustedGasPrice), big.NewFloat(bufferAdjustment))
			adjustedGasPrice, _ = bufferedGasPriceFloat.Int(nil)
		} else if !strings.Contains(err.Error(), BlockFetchingErr) {
			return
		} else {
			L.Debug().
				Msgf("Failed to calculate congestion metric due to: %s. Skipping congestion buffer adjustment", err.Error())

			// set error to nil, as we can still calculate the fees, but without congestion buffer
			// we don't want to return an error in this case
			err = nil
		}
	}

	L.Debug().
		Str("Diff (Wei/Ether)", fmt.Sprintf("%s/%s", big.NewInt(0).Sub(adjustedGasPrice, suggestedGasPrice).String(), WeiToEther(big.NewInt(0).Sub(adjustedGasPrice, suggestedGasPrice)).Text('f', -1))).
		Str("Initial GasPrice (Wei/Ether)", fmt.Sprintf("%s/%s", suggestedGasPrice.String(), WeiToEther(suggestedGasPrice).Text('f', -1))).
		Str("Final GasPrice (Wei/Ether)", fmt.Sprintf("%s/%s", adjustedGasPrice.String(), WeiToEther(adjustedGasPrice).Text('f', -1))).
		Msg("Suggested Legacy fees")

	L.Info().
		Str("GasPrice", fmt.Sprintf("%s wei / %s ether", adjustedGasPrice.String(), WeiToEther(adjustedGasPrice).Text('f', -1))).
		Msg("Calculated suggested Legacy fees")

	return
}

func getAdjustmentFactor(priority string) (float64, error) {
	switch priority {
	case Priority_Degen:
		return 1.5, nil
	case Priority_Fast:
		return 1.2, nil
	case Priority_Standard:
		return 1.0, nil
	case Priority_Slow:
		return 0.8, nil
	default:
		return 0, fmt.Errorf("unknown priority: %s", priority)
	}
}

func getCongestionFactor(congestionClassification string) (float64, error) {
	switch congestionClassification {
	case Congestion_Low:
		return 1.10, nil
	case Congestion_Medium:
		return 1.20, nil
	case Congestion_High:
		return 1.30, nil
	case Congestion_VeryHigh:
		return 1.40, nil
	default:
		return 0, fmt.Errorf("unknown congestion classification: %s", congestionClassification)
	}
}

func classifyCongestion(congestionMetric float64) string {
	switch {
	case congestionMetric < 0.33:
		return Congestion_Low
	case congestionMetric <= 0.66:
		return Congestion_Medium
	case congestionMetric <= 0.75:
		return Congestion_High
	default:
		return Congestion_VeryHigh
	}
}

func (m *Client) HistoricalFeeData(priority string) (baseFee float64, historicalGasTipCap float64, err error) {
	var percentileTip float64

	// based on priority decide, which percentile to use to get historical tip values, when calling FeeHistory
	switch priority {
	case Priority_Degen:
		percentileTip = 100
	case Priority_Fast:
		percentileTip = 99
	case Priority_Standard:
		percentileTip = 50
	case Priority_Slow:
		percentileTip = 25
	default:
		err = fmt.Errorf("unknown priority: %s", priority)
		L.Debug().
			Str("Priority", priority).
			Msgf("Unknown priority: %s", err.Error())

		return
	}

	estimator := NewGasEstimator(m)
	stats, err := estimator.Stats(m.Cfg.Network.GasPriceEstimationBlocks, percentileTip)
	if err != nil {
		L.Debug().
			Msgf("Failed to get fee history due to: %s", err.Error())

		return
	}

	// base fee should still be based on priority, because FeeHistory returns whole base fee history, not just the requested percentile
	switch priority {
	case Priority_Degen:
		baseFee = stats.GasPrice.Max
	case Priority_Fast:
		baseFee = stats.GasPrice.Perc99
	case Priority_Standard:
		baseFee = stats.GasPrice.Perc50
	case Priority_Slow:
		baseFee = stats.GasPrice.Perc25
	default:
		err = fmt.Errorf("unknown priority: %s", priority)
		L.Debug().
			Str("Priority", priority).
			Msgf("Unknown priority: %s", err.Error())

		return
	}

	// since we have already requested reward percentiles based on priority, let's now use the median, i.e. most common tip
	historicalGasTipCap = stats.TipCap.Perc50

	return
}

// calculateGasUsedRatio averages the gas used ratio for a sense of how full blocks are
func calculateGasUsedRatio(headers []*types.Header) float64 {
	if len(headers) == 0 {
		return 0
	}

	var totalRatio float64
	for _, header := range headers {
		if header.GasLimit == 0 {
			continue
		}
		ratio := float64(header.GasUsed) / float64(header.GasLimit)
		totalRatio += ratio
	}
	averageRatio := totalRatio / float64(len(headers))
	return averageRatio
}

func calculateMagnitudeDifference(first, second *big.Float) (int, string) {
	firstFloat, _ := first.Float64()
	secondFloat, _ := second.Float64()

	if firstFloat == 0.0 {
		return -0, "infinite orders of magnitude smaller"
	}

	if secondFloat == 0.0 {
		return -0, "infinite orders of magnitude larger"
	}

	firstOrderOfMagnitude := math.Log10(firstFloat)
	secondOrderOfMagnitude := math.Log10(secondFloat)

	diff := firstOrderOfMagnitude - secondOrderOfMagnitude

	if diff < 0 {
		intDiff := math.Floor(diff)
		return int(intDiff), fmt.Sprintf("%d orders of magnitude smaller", int(math.Abs(intDiff)))
	} else if diff > 0 && diff <= 1 {
		return 0, "the same order of magnitude"
	}

	intDiff := int(math.Ceil(diff))
	return intDiff, fmt.Sprintf("%d orders of magnitude larger", intDiff)
}

func getSafeGasEstimationsAttemptCount(cfg *Config) uint {
	if cfg.Network.GasPriceEstimationAttemptCount == 0 {
		return DefaultGasPriceEstimationsAttemptCount
	}
	return cfg.Network.GasPriceEstimationAttemptCount
}

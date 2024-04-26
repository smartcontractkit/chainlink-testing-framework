package client

import (
	"fmt"
	"strconv"
	"time"

	"github.com/rs/zerolog"
)

// ModulateBaseFeeOverDuration will cause the gas price to rise or drop to a certain percentage of the starting gas price
// over the duration specified.
// Minimum duration is 1 s
// if spike is true, the gas price will rise to the target price
// if spike is false, the gas price will drop to the target price
func (m *RPCClient) ModulateBaseFeeOverDuration(lggr zerolog.Logger, startingBaseFee int64, percentage float64, duration time.Duration, spike bool) error {
	if duration < time.Second {
		return fmt.Errorf("duration must be at least 1s")
	}
	// Calculate the target gas price
	targetBaseFee := float64(startingBaseFee) * (1 + percentage)
	if !spike {
		targetBaseFee = float64(startingBaseFee) * (1 - percentage)
	}
	lggr.Info().
		Int64("Starting Base Fee", startingBaseFee).
		Float64("Percentage", percentage).
		Dur("Duration", duration).
		Int64("Target Base Fee", int64(targetBaseFee)).
		Msg("Modulating base fee per gas over duration")

	// Divide the duration into 10 parts and update the gas price every part
	intTargetBaseFee := int64(targetBaseFee)
	partUpdate := (intTargetBaseFee - startingBaseFee) / 10
	partDuration := duration / 10
	ticker := time.NewTicker(partDuration)
	defer ticker.Stop()
	baseFeeToUpdate := startingBaseFee
	for range ticker.C {
		lggr.Info().
			Int64("Base Fee", baseFeeToUpdate).
			Int64("Updating By", partUpdate).
			Msg("Updating base fee per gas")
		baseFeeToUpdate = baseFeeToUpdate + partUpdate
		if spike {
			if baseFeeToUpdate > intTargetBaseFee {
				baseFeeToUpdate = intTargetBaseFee
			}
		} else {
			if baseFeeToUpdate < intTargetBaseFee {
				baseFeeToUpdate = intTargetBaseFee
			}
		}
		err := m.AnvilSetNextBlockBaseFeePerGas([]interface{}{strconv.FormatInt(baseFeeToUpdate, 10)})
		if err != nil {
			return fmt.Errorf("failed to set base fee %d: %w", baseFeeToUpdate, err)
		}
		lggr.Info().Int64("NextBlockBaseFeePerGas", baseFeeToUpdate).Msg("Updated base fee per gas")
		if baseFeeToUpdate == intTargetBaseFee {
			lggr.Info().
				Int64("Base Fee", baseFeeToUpdate).
				Msg("Reached target base fee")
			return nil
		}
	}
	return fmt.Errorf("failed to reach target base fee %d", intTargetBaseFee)
}

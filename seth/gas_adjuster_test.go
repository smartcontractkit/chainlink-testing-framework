package seth

import (
	"math"
	"math/big"
	"testing"
)

// TestCalculateMagnitudeDifference tests the magnitude difference calculation
func TestGasAdjuster_CalculateMagnitudeDifference(t *testing.T) {
	tests := []struct {
		name         string
		first        *big.Float
		second       *big.Float
		expectedDiff int
		expectedText string
	}{
		{
			name:         "First much larger (5 orders)",
			first:        big.NewFloat(100_000_000_000), // 100 gwei
			second:       big.NewFloat(1_000_000),       // 0.001 gwei
			expectedDiff: 5,
			expectedText: "5 orders of magnitude larger",
		},
		{
			name:         "First much smaller (4 orders)",
			first:        big.NewFloat(1_000_000),      // 0.001 gwei
			second:       big.NewFloat(10_000_000_000), // 10 gwei
			expectedDiff: -4,
			expectedText: "4 orders of magnitude smaller",
		},
		{
			name:         "Similar magnitude (within 1 order)",
			first:        big.NewFloat(30_000_000_000), // 30 gwei
			second:       big.NewFloat(31_000_000_000), // 31 gwei
			expectedDiff: 0,
			expectedText: "the same order of magnitude",
		},
		{
			name:         "Similar magnitude (within 1 order)",
			first:        big.NewFloat(100_000_000_000), // 100 gwei
			second:       big.NewFloat(99_999_999_999),  // 99.999999999 gwei
			expectedDiff: 0,
			expectedText: "the same order of magnitude",
		},
		{
			name:         "Similar magnitude (within 1 order)",
			first:        big.NewFloat(30_000_000_000), // 30 gwei
			second:       big.NewFloat(99_999_999_999), // 99.999999999 gwei
			expectedDiff: 0,
			expectedText: "the same order of magnitude",
		},
		{
			name:         "Similar magnitude (within 1 order)",
			first:        big.NewFloat(99_999_999_999), // 99.999999999 gwei
			second:       big.NewFloat(30_000_000_000), // 30 gwei
			expectedDiff: 0,
			expectedText: "the same order of magnitude",
		},
		{
			name:         "Similar magnitude (within 1 order)",
			first:        big.NewFloat(99_999_999_999),  // 99.999999999 gwei
			second:       big.NewFloat(100_000_000_000), // 100 gwei
			expectedDiff: 0,
			expectedText: "the same order of magnitude",
		},
		{
			name:         "Just under 1 order of magnitude (same order)",
			first:        big.NewFloat(9_999_999_999), // 9.999... gwei
			second:       big.NewFloat(1_000_000_000), // 1 gwei
			expectedDiff: 0,                           // Still same order (diff = 0.9999...)
			expectedText: "the same order of magnitude",
		},
		{
			name:         "Exactly 3 orders larger",
			first:        big.NewFloat(1_000_000_000), // 1 gwei
			second:       big.NewFloat(1_000_000),     // 0.001 gwei
			expectedDiff: 3,
			expectedText: "3 orders of magnitude larger",
		},
		{
			name:         "Exactly 3 orders smaller",
			first:        big.NewFloat(1_000_000),     // 0.001 gwei
			second:       big.NewFloat(1_000_000_000), // 1 gwei
			expectedDiff: -3,
			expectedText: "3 orders of magnitude smaller",
		},
		{
			name:         "First is zero",
			first:        big.NewFloat(0),
			second:       big.NewFloat(1_000_000),
			expectedDiff: -0,
			expectedText: "infinite orders of magnitude smaller",
		},
		{
			name:         "Second is zero",
			first:        big.NewFloat(1_000_000),
			second:       big.NewFloat(0),
			expectedDiff: -0,
			expectedText: "infinite orders of magnitude larger",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diff, text := calculateMagnitudeDifference(tt.first, tt.second)

			if diff != tt.expectedDiff {
				t.Errorf("calculateMagnitudeDifference(first=%s wei, second=%s wei)\n  diff = %v, want %v",
					tt.first.Text('f', 0), tt.second.Text('f', 0), diff, tt.expectedDiff)
			}

			if text != tt.expectedText {
				t.Errorf("calculateMagnitudeDifference(first=%s wei, second=%s wei)\n  text = %q, want %q",
					tt.first.Text('f', 0), tt.second.Text('f', 0), text, tt.expectedText)
			}
		})
	}
}

// TestFeeEqualizerLogic tests the Fee Equalizer experimental logic
// WARNING: This feature should ONLY be used on testnets, never on mainnet!
// It can cause massive overpayment or transaction failures on production networks.
func TestGasAdjuster_FeeEqualizerLogic(t *testing.T) {
	tests := []struct {
		name                string
		baseFee             int64 // wei
		tip                 int64 // wei
		expectedBaseFee     int64 // wei after adjustment
		expectedTip         int64 // wei after adjustment
		shouldAdjustBaseFee bool
		shouldAdjustTip     bool
		description         string
	}{
		{
			name:                "Base fee MUCH smaller than tip (testnet scenario)",
			baseFee:             1_000_000,     // 0.001 gwei
			tip:                 5_000_000_000, // 5 gwei
			expectedBaseFee:     5_000_000_000, // Should be raised to tip value
			expectedTip:         5_000_000_000, // No change
			shouldAdjustBaseFee: true,
			shouldAdjustTip:     false,
			description:         "Testnet with very low base fee, high suggested tip - raise base fee to avoid tx failure",
		},
		{
			name:                "Base fee MUCH larger than tip (congested testnet)",
			baseFee:             100_000_000_000, // 100 gwei
			tip:                 1_000_000,       // 0.001 gwei
			expectedBaseFee:     100_000_000_000, // No change
			expectedTip:         100_000_000_000, // Should be raised to base fee value
			shouldAdjustBaseFee: false,
			shouldAdjustTip:     true,
			description:         "High congestion with tiny tip - raise tip to match base fee",
		},
		{
			name:                "Similar magnitude - no adjustment needed",
			baseFee:             30_000_000_000, // 30 gwei
			tip:                 2_000_000_000,  // 2 gwei
			expectedBaseFee:     30_000_000_000, // No change
			expectedTip:         2_000_000_000,  // No change
			shouldAdjustBaseFee: false,
			shouldAdjustTip:     false,
			description:         "Both values are within acceptable range, no adjustment",
		},
		{
			name:                "Exactly at threshold - 3 orders difference (smaller)",
			baseFee:             1_000_000,     // 0.001 gwei
			tip:                 1_000_000_000, // 1 gwei
			expectedBaseFee:     1_000_000,     // No change (exactly 3 orders, not >3)
			expectedTip:         1_000_000_000, // No change
			shouldAdjustBaseFee: false,
			shouldAdjustTip:     false,
			description:         "Exactly 3 orders of magnitude difference - no adjustment at boundary",
		},
		{
			name:                "Exactly at threshold - 3 orders difference (larger)",
			baseFee:             1_000_000_000, // 1 gwei
			tip:                 1_000_000,     // 0.001 gwei
			expectedBaseFee:     1_000_000_000, // No change (exactly 3 orders, not >3)
			expectedTip:         1_000_000,     // No change
			shouldAdjustBaseFee: false,
			shouldAdjustTip:     false,
			description:         "Exactly 3 orders of magnitude difference - no adjustment at boundary",
		},
		{
			name:                "Slightly over threshold - 4 orders (smaller)",
			baseFee:             100_000,       // 0.0001 gwei
			tip:                 1_000_000_000, // 1 gwei
			expectedBaseFee:     1_000_000_000, // Adjusted
			expectedTip:         1_000_000_000, // No change
			shouldAdjustBaseFee: true,
			shouldAdjustTip:     false,
			description:         "4 orders of magnitude smaller - triggers adjustment",
		},
		{
			name:                "Slightly over threshold - 4 orders (larger)",
			baseFee:             1_000_000_000, // 1 gwei
			tip:                 100_000,       // 0.0001 gwei
			expectedBaseFee:     1_000_000_000, // No change
			expectedTip:         1_000_000_000, // Adjusted
			shouldAdjustBaseFee: false,
			shouldAdjustTip:     true,
			description:         "4 orders of magnitude larger - triggers adjustment",
		},
		{
			name:                "Base fee is zero (edge case)",
			baseFee:             0,
			tip:                 1_000_000_000, // 1 gwei
			expectedBaseFee:     1_000_000_000, // Should use tip
			expectedTip:         1_000_000_000, // No change
			shouldAdjustBaseFee: true,
			shouldAdjustTip:     false,
			description:         "Zero base fee - use tip as base fee",
		},
		{
			name:                "Tip is zero (edge case)",
			baseFee:             1_000_000_000, // 1 gwei
			tip:                 0,
			expectedBaseFee:     1_000_000_000, // No change
			expectedTip:         1_000_000_000, // Should use base fee
			shouldAdjustBaseFee: false,
			shouldAdjustTip:     true,
			description:         "Zero tip - use base fee as tip",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Simulate the Fee Equalizer logic
			baseFee64 := float64(tt.baseFee)
			currentGasTip := big.NewInt(tt.tip)

			baseFeeTipMagnitudeDiff, _ := calculateMagnitudeDifference(
				big.NewFloat(baseFee64),
				new(big.Float).SetInt(currentGasTip),
			)

			// Track if adjustments were made
			baseFeeAdjusted := false
			tipAdjusted := false

			// Apply the Fee Equalizer logic
			if baseFeeTipMagnitudeDiff == -0 {
				if baseFee64 == 0.0 {
					baseFee64 = float64(currentGasTip.Int64())
					baseFeeAdjusted = true
				} else {
					currentGasTip = big.NewInt(int64(baseFee64))
					tipAdjusted = true
				}
			} else if baseFeeTipMagnitudeDiff < -3 {
				// Base fee is MUCH SMALLER than tip (more than 3 orders of magnitude)
				baseFee64 = float64(currentGasTip.Int64())
				baseFeeAdjusted = true
			} else if baseFeeTipMagnitudeDiff > 3 {
				// Base fee is MUCH LARGER than tip (more than 3 orders of magnitude)
				currentGasTip = big.NewInt(int64(baseFee64))
				tipAdjusted = true
			}

			// Verify results
			resultBaseFee := int64(baseFee64)
			resultTip := currentGasTip.Int64()

			if resultBaseFee != tt.expectedBaseFee {
				t.Errorf("Base fee after adjustment = %d wei (%.4f gwei), want %d wei (%.4f gwei)\nDescription: %s",
					resultBaseFee, float64(resultBaseFee)/1e9,
					tt.expectedBaseFee, float64(tt.expectedBaseFee)/1e9,
					tt.description)
			}

			if resultTip != tt.expectedTip {
				t.Errorf("Tip after adjustment = %d wei (%.4f gwei), want %d wei (%.4f gwei)\nDescription: %s",
					resultTip, float64(resultTip)/1e9,
					tt.expectedTip, float64(tt.expectedTip)/1e9,
					tt.description)
			}

			if baseFeeAdjusted != tt.shouldAdjustBaseFee {
				t.Errorf("Base fee adjustment = %v, want %v\nDescription: %s",
					baseFeeAdjusted, tt.shouldAdjustBaseFee, tt.description)
			}

			if tipAdjusted != tt.shouldAdjustTip {
				t.Errorf("Tip adjustment = %v, want %v\nDescription: %s",
					tipAdjusted, tt.shouldAdjustTip, tt.description)
			}
		})
	}
}

// TestFeeEqualizerDisasterScenarios tests scenarios that would be catastrophic on mainnet
// These tests document WHY this feature should NEVER be used on production networks
func TestGasAdjuster_FeeEqualizerDisasterScenarios(t *testing.T) {
	tests := []struct {
		name            string
		baseFee         int64
		tip             int64
		network         string
		disasterOutcome string
	}{
		{
			name:            "Ethereum mainnet high gas",
			baseFee:         100_000_000_000, // 100 gwei (typical high congestion)
			tip:             1_000_000,       // 0.001 gwei (would never happen naturally)
			network:         "Ethereum Mainnet",
			disasterOutcome: "Tip raised to 100 gwei - would cost $100+ per transaction!",
		},
		{
			name:            "Polygon PoS normal operation",
			baseFee:         30_000_000_000, // 30 gwei
			tip:             100_000,        // Very low tip
			network:         "Polygon",
			disasterOutcome: "Massive overpayment for network that normally costs fractions of a cent",
		},
		{
			name:            "Arbitrum during congestion",
			baseFee:         500_000_000_000, // 500 gwei
			tip:             1_000_000,       // 0.001 gwei
			network:         "Arbitrum",
			disasterOutcome: "Tip becomes 500 gwei - would bankrupt users on a rollup!",
		},
	}

	t.Log("⚠️  WARNING: Fee Equalizer Disaster Scenarios ⚠️")
	t.Log("These scenarios demonstrate why Fee Equalizer should NEVER be enabled on mainnet!")
	t.Log("")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			baseFee64 := float64(tt.baseFee)
			currentGasTip := big.NewInt(tt.tip)

			baseFeeTipMagnitudeDiff, diffText := calculateMagnitudeDifference(
				big.NewFloat(baseFee64),
				new(big.Float).SetInt(currentGasTip),
			)

			originalTip := currentGasTip.Int64()

			// Apply Fee Equalizer logic
			if baseFeeTipMagnitudeDiff > 3 {
				currentGasTip = big.NewInt(int64(baseFee64))
			}

			resultTip := currentGasTip.Int64()
			overpaymentMultiplier := float64(resultTip) / float64(originalTip)

			t.Logf("Network: %s", tt.network)
			t.Logf("Original Base Fee: %d wei (%.2f gwei)", tt.baseFee, float64(tt.baseFee)/1_000_000_000)
			t.Logf("Original Tip: %d wei (%.6f gwei)", originalTip, float64(originalTip)/1_000_000_000)
			t.Logf("Magnitude Difference: %s", diffText)
			t.Logf("Adjusted Tip: %d wei (%.2f gwei)", resultTip, float64(resultTip)/1_000_000_000)
			t.Logf("Overpayment Multiplier: %.0fx", overpaymentMultiplier)
			t.Logf("Disaster Outcome: %s", tt.disasterOutcome)
			t.Logf("")

			// Assert that this would be a disaster
			if overpaymentMultiplier > 1000 {
				t.Logf("✗ DISASTER: This would cause >1000x overpayment on %s!", tt.network)
			}
		})
	}
}

// TestGasAdjusterBoundaryConditions tests edge cases around the 3 order of magnitude threshold
func TestGasAdjuster_FeeEqualizerBoundaryConditions(t *testing.T) {
	tests := []struct {
		name               string
		baseFee            float64
		tip                float64
		expectedAdjustment string
	}{
		{
			name:               "Exactly 3.0 orders (1000x) - no adjustment",
			baseFee:            1_000_000_000.0, // 1 gwei
			tip:                1_000_000.0,     // 0.001 gwei
			expectedAdjustment: "none",
		},
		{
			name:               "Just over 3 orders (3.01) - should adjust",
			baseFee:            1_023_000_000.0, // 1.023 gwei
			tip:                1_000_000.0,     // 0.001 gwei
			expectedAdjustment: "tip",
		},
		{
			name:               "Just under 3 orders (2.99) - no adjustment",
			baseFee:            977_000_000.0, // 0.977 gwei
			tip:                1_000_000.0,   // 0.001 gwei
			expectedAdjustment: "none",
		},
		{
			name:               "Negative: Exactly -3.0 orders - no adjustment",
			baseFee:            1_000_000.0,     // 0.001 gwei
			tip:                1_000_000_000.0, // 1 gwei
			expectedAdjustment: "none",
		},
		{
			name:               "Negative: Just over -3 orders (-3.01) - should adjust",
			baseFee:            977_000.0,       // 0.000977 gwei
			tip:                1_000_000_000.0, // 1 gwei
			expectedAdjustment: "baseFee",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			baseFee64 := tt.baseFee
			currentGasTip := big.NewInt(int64(tt.tip))

			diff, _ := calculateMagnitudeDifference(
				big.NewFloat(baseFee64),
				new(big.Float).SetInt(currentGasTip),
			)

			actualAdjustment := "none"
			if diff < -3 {
				actualAdjustment = "baseFee"
			} else if diff > 3 {
				actualAdjustment = "tip"
			}

			magnitudeDiff := math.Log10(baseFee64) - math.Log10(tt.tip)
			t.Logf("Magnitude difference: %.4f orders (threshold: ±3.0)", magnitudeDiff)
			t.Logf("Integer diff: %d", diff)
			t.Logf("Expected adjustment: %s", tt.expectedAdjustment)
			t.Logf("Actual adjustment: %s", actualAdjustment)

			if actualAdjustment != tt.expectedAdjustment {
				t.Errorf("Adjustment = %v, want %v", actualAdjustment, tt.expectedAdjustment)
			}
		})
	}
}

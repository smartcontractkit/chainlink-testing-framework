package benchspy

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculatePercentile(t *testing.T) {
	t.Run("basic percentile calculations", func(t *testing.T) {
		numbers := []float64{1, 2, 3, 4, 5}

		// Test median (50th percentile)
		assert.Equal(t, 3.0, CalculatePercentile(numbers, 0.5))

		// Test minimum (0th percentile)
		assert.Equal(t, 1.0, CalculatePercentile(numbers, 0))

		// Test maximum (100th percentile)
		assert.Equal(t, 5.0, CalculatePercentile(numbers, 1))
	})

	t.Run("unsorted input", func(t *testing.T) {
		numbers := []float64{5, 2, 1, 4, 3}
		assert.Equal(t, 3.0, CalculatePercentile(numbers, 0.5))
	})

	t.Run("interpolation cases", func(t *testing.T) {
		numbers := []float64{1, 2, 3, 4}

		// Test 25th percentile (should interpolate between 1 and 2)
		expected25 := 1.75
		assert.InDelta(t, expected25, CalculatePercentile(numbers, 0.25), 0.000001)

		// Test 75th percentile (should interpolate between 3 and 4)
		expected75 := 3.25
		assert.InDelta(t, expected75, CalculatePercentile(numbers, 0.75), 0.000001)
	})

	t.Run("single element", func(t *testing.T) {
		numbers := []float64{42}
		assert.Equal(t, 42.0, CalculatePercentile(numbers, 0.5))
	})

	t.Run("duplicate values", func(t *testing.T) {
		numbers := []float64{1, 2, 2, 3, 3, 3, 4}
		assert.Equal(t, 3.0, CalculatePercentile(numbers, 0.5))
	})

	t.Run("panic on empty slice", func(t *testing.T) {
		assert.Panics(t, func() {
			CalculatePercentile([]float64{}, 0.5)
		})
	})

	t.Run("panic on invalid percentile - negative", func(t *testing.T) {
		assert.Panics(t, func() {
			CalculatePercentile([]float64{1, 2, 3}, -0.1)
		})
	})

	t.Run("panic on invalid percentile - greater than 1", func(t *testing.T) {
		assert.Panics(t, func() {
			CalculatePercentile([]float64{1, 2, 3}, 1.1)
		})
	})

	t.Run("large dataset", func(t *testing.T) {
		numbers := make([]float64, 1000)
		for i := 0; i < 1000; i++ {
			numbers[i] = float64(i)
		}
		// 90th percentile of 0-999 should be 899
		assert.InDelta(t, 899.1, CalculatePercentile(numbers, 0.9), 0.000001)
	})
}

func TestStringSliceToFloat64Slice(t *testing.T) {
	t.Run("valid conversion", func(t *testing.T) {
		input := []string{"1.0", "2.5", "3.14"}
		expected := []float64{1.0, 2.5, 3.14}

		result, err := StringSliceToFloat64Slice(input)
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("invalid number format", func(t *testing.T) {
		input := []string{"1.0", "invalid", "3.14"}
		result, err := StringSliceToFloat64Slice(input)
		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("empty slice", func(t *testing.T) {
		input := []string{}
		result, err := StringSliceToFloat64Slice(input)
		assert.NoError(t, err)
		assert.Empty(t, result)
	})

	t.Run("scientific notation", func(t *testing.T) {
		input := []string{"1e-10", "2e5"}
		expected := []float64{1e-10, 2e5}

		result, err := StringSliceToFloat64Slice(input)
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})
}

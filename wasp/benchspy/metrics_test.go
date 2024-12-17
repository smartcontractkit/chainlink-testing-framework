package benchspy

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBenchSpy_StringSliceToFloat64Slice(t *testing.T) {
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

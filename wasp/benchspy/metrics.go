package benchspy

import (
	"math"
	"sort"
	"strconv"
)

func CalculatePercentile(numbers []float64, percentile float64) float64 {
	// Sort the slice
	sort.Float64s(numbers)

	n := len(numbers)
	if n == 0 {
		panic("cannot calculate percentile of an empty slice")
	}
	if percentile < 0 || percentile > 1 {
		panic("percentile must be between 0 and 1")
	}

	// Calculate the rank (index)
	rank := percentile * float64(n-1) // Use n-1 for zero-based indexing
	lowerIndex := int(math.Floor(rank))
	upperIndex := int(math.Ceil(rank))

	// Interpolate if necessary
	if lowerIndex == upperIndex {
		return numbers[lowerIndex]
	}
	weight := rank - float64(lowerIndex)
	return numbers[lowerIndex]*(1-weight) + numbers[upperIndex]*weight
}

func StringSliceToFloat64Slice(s []string) ([]float64, error) {
	numbers := make([]float64, len(s))
	for i, str := range s {
		var err error
		numbers[i], err = strconv.ParseFloat(str, 64)
		if err != nil {
			return nil, err
		}
	}
	return numbers, nil
}

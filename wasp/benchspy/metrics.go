package benchspy

import (
	"strconv"
)

// StringSliceToFloat64Slice converts a slice of strings to a slice of float64 values.
// It returns an error if any string cannot be parsed as a float64, making it useful for data conversion tasks.
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

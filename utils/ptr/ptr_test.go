package ptr

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// write test for Ptr
// TestPtr tests the Ptr function.
func TestPtr(t *testing.T) {
	type sampleStruct struct {
		Name string
		Age  int
	}
	nonZero := sampleStruct{Name: "John", Age: 30}
	zero := sampleStruct{}
	var nilPointer *sampleStruct
	pointerToNonZero := Ptr(nonZero)
	require.NotNil(t, pointerToNonZero, "Ptr returned nil")
	require.Equal(t, nonZero, Value(pointerToNonZero), "Ptr did not return the expected value")
	pointerToZero := Ptr(zero)
	require.NotNil(t, pointerToZero, "Ptr returned nil")
	require.Equal(t, zero, Value(pointerToZero), "Ptr did not return the expected value")
	require.Equal(t, zero, Value(nilPointer), "Value did not return the zero value of T")
}

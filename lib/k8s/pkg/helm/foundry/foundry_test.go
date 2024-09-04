package foundry

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestChartNameUniqueness(t *testing.T) {
	props := defaultProps()

	// Create two Chart instances
	chart1 := NewVersioned("", props)
	chart2 := NewVersioned("", props)

	// Assert that the names of the charts are unique
	require.NotEqual(t, chart1.GetName(), chart2.GetName(), "Chart names should be unique")
}

func TestChartNameWithFullNameOverrideUniqueness(t *testing.T) {
	// Generate unique names for testing
	name1 := fmt.Sprintf("custom-%s", uuid.New().String()[0:5])
	name2 := fmt.Sprintf("custom-%s", uuid.New().String()[0:5])

	// Create two Chart instances with different fullnameOverride
	chart1 := NewVersioned("", &Props{
		Values: map[string]interface{}{
			"fullnameOverride": name1,
		},
	})
	chart2 := NewVersioned("", &Props{
		Values: map[string]interface{}{
			"fullnameOverride": name2,
		},
	})

	// Assert that the names of the charts are unique and correct
	require.Equal(t, name1, chart1.GetName(), "Chart name should match the fullnameOverride")
	require.Equal(t, name2, chart2.GetName(), "Chart name should match the fullnameOverride")
	require.NotEqual(t, chart1.GetName(), chart2.GetName(), "Chart names should be unique when fullnameOverrides are different")
}

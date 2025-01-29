package parrot

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCage(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		route       *Route
		expectedErr error
	}{
		{
			name: "basic route",
			route: &Route{
				Method:             http.MethodGet,
				Path:               "/test",
				RawResponseBody:    "Squawk",
				ResponseStatusCode: http.StatusOK,
			},
		},
		{
			name: "wildcard route",
			route: &Route{
				Method:             http.MethodGet,
				Path:               "/*",
				RawResponseBody:    "Squawk",
				ResponseStatusCode: http.StatusOK,
			},
		},
		{
			name: "nested route",
			route: &Route{
				Method:             http.MethodGet,
				Path:               "/test/nested",
				RawResponseBody:    "Squawk",
				ResponseStatusCode: http.StatusOK,
			},
		},
		{
			name: "nested wild ard route",
			route: &Route{
				Method:             http.MethodGet,
				Path:               "/test/nested/*",
				RawResponseBody:    "Squawk",
				ResponseStatusCode: http.StatusOK,
			},
		},
		{
			name: "multi-nested wild card route",
			route: &Route{
				Method:             http.MethodGet,
				Path:               "/test/*/nested/*",
				RawResponseBody:    "Squawk",
				ResponseStatusCode: http.StatusOK,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			c := newCage()
			require.NotNil(t, c, "Cage should not be nil")

			// Create the new route
			err := c.newRoute(tc.route)
			require.NoError(t, err, "newRoute should not return an error")

			// Check that the proper cage level got created
			cageLevel, err := c.getCageLevel(tc.route.Path, false)
			require.NoError(t, err, "cageLevel should not return an error")
			require.NotNil(t, cageLevel, "cageLevel should not return nil")
			// Check that the route was created and can be found from the cage level
			route, found, err := cageLevel.route(tc.route.Segment(), tc.route.Method)
			require.NoError(t, err, "route should not return an error")
			require.True(t, found, "route should be found in the found cage level")
			require.Equal(t, tc.route, route, "route should be equal to the route in the cage")

			// Check that the route was created and can be found from the base cage
			route, err = c.getRoute(tc.route.Path, tc.route.Method)
			require.NoError(t, err, "getRoute should not return an error")
			require.NotNil(t, route, "route should not be nil")
			require.Equal(t, tc.route, route, "route should be equal to the route in the cage")
			allRoutes := c.routes()
			require.Len(t, allRoutes, 1, "cage should have 1 route")

			// Check that we can properly delete the route
			err = c.deleteRoute(tc.route)
			require.NoError(t, err, "deleteRoute should not return an error")
			_, err = c.getRoute(tc.route.Path, tc.route.Method)
			require.ErrorIs(t, err, ErrRouteNotFound, "should error getting route after deleting it")
			allRoutes = c.routes()
			require.Empty(t, allRoutes, "cage should no longer have any routes")

		})
	}
}

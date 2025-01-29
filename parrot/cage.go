package parrot

import (
	"fmt"
	"path/filepath"
	"strings"
	"sync"
)

// Cage is a container for all routes and sub-cages to handle routing and wildcard matching for a parrot
// Note: Should only be used internally by the parrot.
type cage struct {
	*cageLevel
}

// MethodAny will match to any other method
const MethodAny = "ANY"

// CageLevel holds a single level of routes and further sub cages
// Note: Should only be used internally by the parrot.
type cageLevel struct {
	// cagePath is the path to this cage level
	cagePath string
	// rwMu is a read write mutex for the cage level
	rwMu sync.RWMutex
	// TODO: Make all lowercase
	// routes contains all of the plain routes at this current cage level
	// route.path -> route.method -> route
	routes map[string]map[string]*Route
	// wildCardRoutes contains all the wildcard routes at this current cage level
	// route.path -> route.method -> route
	wildCardRoutes map[string]map[string]*Route
	// subCages contains sub cages at this current cage level
	// cage name -> cage level
	subCages map[string]*cageLevel
	// wildCardSubCages contains wildcard sub cages at this current cage level
	// cage name -> cage level
	wildCardSubCages map[string]*cageLevel
}

// newCage creates a new cage with an empty cage level for a new parrot instance
func newCage() *cage {
	return &cage{
		cageLevel: newCageLevel("/"),
	}
}

// newCageLevel creates a new cageLevel with empty maps
func newCageLevel(cagePath string) *cageLevel {
	return &cageLevel{
		cagePath:         cagePath,
		routes:           make(map[string]map[string]*Route),
		wildCardRoutes:   make(map[string]map[string]*Route),
		subCages:         make(map[string]*cageLevel),
		wildCardSubCages: make(map[string]*cageLevel),
	}
}

// cageLevel searches for a cage level based on the path provided
// If createMode is true, it will create any cage levels if they don't exist
func (c *cage) getCageLevel(path string, createMode bool) (cageLevel *cageLevel, err error) {
	splitPath := strings.Split(path, "/")
	splitPath = splitPath[:len(splitPath)-1] // Only looking for the cage level, exclude the route
	if splitPath[0] == "" {
		splitPath = splitPath[1:] // Remove the empty string at the beginning of the split
	}
	currentCageLevel := c.cageLevel

	for _, pathSegment := range splitPath { // Iterate through each path segment to look for matches
		cageLevel, found, err := currentCageLevel.subCageLevel(pathSegment, createMode)
		if err != nil {
			return nil, err
		}
		if found {
			currentCageLevel = cageLevel
			continue
		}

		if !found {
			return nil, newDynamicError(ErrCageNotFound, fmt.Sprintf("path: '%s'", path))
		}
	}

	return currentCageLevel, nil
}

// getRoute searches for a route based on the path provided
func (c *cage) getRoute(routePath, routeMethod string) (*Route, error) {
	cageLevel, err := c.getCageLevel(routePath, false)
	if err != nil {
		return nil, err
	}
	routeSegments := strings.Split(routePath, "/")
	if len(routeSegments) == 0 {
		return nil, ErrRouteNotFound
	}
	routeSegment := routeSegments[len(routeSegments)-1]

	route, found, err := cageLevel.route(routeSegment, routeMethod)
	if err != nil {
		return nil, err
	}
	if found {
		return route, nil
	}

	return nil, ErrRouteNotFound
}

// newRoute creates a new route, creating new cages if necessary
func (c *cage) newRoute(route *Route) error {
	cageLevel, err := c.getCageLevel(route.Path, true)
	if err != nil {
		return err
	}

	cageLevel.newRoute(route)

	return nil
}

// deleteRoute deletes a route
func (c *cage) deleteRoute(route *Route) error {
	cageLevel, err := c.getCageLevel(route.Path, false)
	if err != nil {
		return err
	}

	if strings.Contains(route.Segment(), "*") {
		cageLevel.rwMu.RLock()
		if _, found := cageLevel.wildCardRoutes[route.Segment()][route.Method]; !found {
			cageLevel.rwMu.RUnlock()
			return ErrRouteNotFound
		}
		cageLevel.rwMu.RUnlock()

		cageLevel.rwMu.Lock()
		delete(cageLevel.wildCardRoutes[route.Segment()], route.Method)
		cageLevel.rwMu.Unlock()
	} else {
		cageLevel.rwMu.RLock()
		if _, found := cageLevel.routes[route.Segment()][route.Method]; !found {
			cageLevel.rwMu.RUnlock()
			return ErrRouteNotFound
		}
		cageLevel.rwMu.RUnlock()

		cageLevel.rwMu.Lock()
		delete(cageLevel.routes[route.Segment()], route.Method)
		cageLevel.rwMu.Unlock()
	}

	return nil
}

// routes returns all the routes in the cage
func (c *cage) routes() []*Route {
	return c.routesRecursive()
}

// routesRecursive returns all the routes in the cage recursively.
// Should only be used internally by the cage. Use routes() instead.
func (cl *cageLevel) routesRecursive() (routes []*Route) {
	// Add all the routes at this level
	cl.rwMu.RLock()
	for _, routePath := range cl.routes {
		for _, route := range routePath {
			routes = append(routes, route)
		}
	}

	// Add all the wildcard routes at this level
	for _, routePath := range cl.wildCardRoutes {
		for _, route := range routePath {
			routes = append(routes, route)
		}
	}
	cl.rwMu.RUnlock()

	for _, subCage := range cl.subCages {
		routes = append(routes, subCage.routesRecursive()...)
	}
	for _, subCage := range cl.wildCardSubCages {
		routes = append(routes, subCage.routesRecursive()...)
	}

	return routes
}

// route searches for a route based on the route segment provided
func (cl *cageLevel) route(routeSegment, routeMethod string) (route *Route, found bool, err error) {
	// First check for an exact match
	cl.rwMu.RLock()
	defer cl.rwMu.RUnlock()

	if _, ok := cl.routes[routeSegment]; ok {
		if route, found = cl.routes[routeSegment][routeMethod]; found {
			return route, true, nil
		}
		if route, found = cl.routes[routeSegment][MethodAny]; found { // Fallthrough to any method if it's designed
			return route, true, nil
		}
	}

	// if not, look for wildcard routes
	for wildCardPattern, routePath := range cl.wildCardRoutes {
		pathMatch, err := filepath.Match(wildCardPattern, routeSegment)
		if err != nil {
			return nil, false, newDynamicError(ErrInvalidPath, err.Error())
		}
		if pathMatch {
			// Found a path match, now check for the method
			if route, found = routePath[routeMethod]; found {
				return route, true, nil
			}
			if route, found = routePath[MethodAny]; found {
				return route, true, nil
			}
		}
	}

	return nil, false, nil
}

// subCageLevel searches for a sub cage level based on the segment provided
// if createMode is true, it will create the cage level if it doesn't exist
func (cl *cageLevel) subCageLevel(subCageSegment string, createMode bool) (cageLevel *cageLevel, found bool, err error) {
	// First check for an exact match
	cl.rwMu.RLock()
	if cageLevel, exists := cl.subCages[subCageSegment]; exists {
		cl.rwMu.RUnlock()
		return cageLevel, true, nil
	}

	// if not, look for wildcard cages
	for wildCardPattern, cageLevel := range cl.wildCardSubCages {
		match, err := filepath.Match(wildCardPattern, subCageSegment)
		if err != nil {
			cl.rwMu.RUnlock()
			return nil, false, newDynamicError(ErrInvalidPath, err.Error())
		}
		if match {
			cl.rwMu.RUnlock()
			return cageLevel, true, nil
		}
	}
	cl.rwMu.RUnlock()

	// We didn't find a match, so we'll create a new cage level if we're in create mode
	if createMode {
		newCage := newCageLevel(filepath.Join(cl.cagePath, subCageSegment))
		cl.rwMu.Lock()
		defer cl.rwMu.Unlock()
		if strings.Contains(subCageSegment, "*") {
			cl.wildCardSubCages[subCageSegment] = newCage
		} else {
			cl.subCages[subCageSegment] = newCage
		}
		return newCage, true, nil
	}

	return nil, false, nil
}

// newRoute creates a new route in the cage level
func (cl *cageLevel) newRoute(route *Route) {
	cl.rwMu.Lock()
	defer cl.rwMu.Unlock()
	if strings.Contains(route.Segment(), "*") {
		if _, found := cl.wildCardRoutes[route.Segment()]; !found {
			cl.wildCardRoutes[route.Segment()] = make(map[string]*Route)
		}
		cl.wildCardRoutes[route.Segment()][route.Method] = route
	} else {
		if _, found := cl.routes[route.Segment()]; !found {
			cl.routes[route.Segment()] = make(map[string]*Route)
		}
		cl.routes[route.Segment()][route.Method] = route
	}
}

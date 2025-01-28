package parrot

import (
	"fmt"
	"path/filepath"
	"strings"
	"sync"
)

// Cage is a container for all routes and sub-cages to handle routing and wildcard matching for a parrot
// Note: Should only be used internally by the parrot.
type Cage struct {
	*CageLevel
}

// CageLevel holds a single level of routes and further sub cages
// Note: Should only be used internally by the parrot.
type CageLevel struct {
	// Routes contains all of the plain routes at this current cage level
	// path -> route
	Routes     map[string]*Route `json:"routes"`
	routesRWMu sync.RWMutex      // sync.Map might be better here, but eh
	// WildCardRoutes contains all the wildcard routes at this current cage level
	// path -> route
	WildCardRoutes     map[string]*Route `json:"wild_card_routes"`
	wildCardRoutesRWMu sync.RWMutex
	// SubCages contains sub cages at this current cage level
	// cage name -> cage level
	SubCages     map[string]*CageLevel `json:"sub_cages"`
	subCagesRWMu sync.RWMutex
	// WildCardSubCages contains wildcard sub cages at this current cage level
	// cage name -> cage level
	WildCardSubCages     map[string]*CageLevel `json:"wild_card_sub_cages"`
	wildCardSubCagesRWMu sync.RWMutex
}

// newCage creates a new cage with an empty cage level for a new parrot instance
func newCage() *Cage {
	return &Cage{
		CageLevel: newCageLevel(),
	}
}

// newCageLevel creates a new cageLevel with empty maps
func newCageLevel() *CageLevel {
	return &CageLevel{
		Routes:           make(map[string]*Route),
		WildCardRoutes:   make(map[string]*Route),
		SubCages:         make(map[string]*CageLevel),
		WildCardSubCages: make(map[string]*CageLevel),
	}
}

// cageLevel searches for a cage level based on the path provided
// If createMode is true, it will create any cage levels if they don't exist
func (c *Cage) cageLevel(path string, createMode bool) (cageLevel *CageLevel, routeSegment string, err error) {
	splitPath := strings.Split(path, "/")
	routeSegment = splitPath[len(splitPath)-1] // The final path segment is the route
	splitPath = splitPath[:len(splitPath)-1]   // Only looking for the cage level, exclude the route
	if splitPath[0] == "" {
		splitPath = splitPath[1:] // Remove the empty string at the beginning of the split
	}
	currentCageLevel := c.CageLevel

	for _, pathSegment := range splitPath { // Iterate through each path segment to look for matches
		cageLevel, found, err := currentCageLevel.subCageLevel(pathSegment, createMode)
		if err != nil {
			return nil, routeSegment, err
		}
		if found {
			currentCageLevel = cageLevel
			continue
		}

		if !found {
			return nil, routeSegment, newDynamicError(ErrCageNotFound, fmt.Sprintf("path: '%s'", path))
		}
	}

	return currentCageLevel, routeSegment, nil
}

// getRoute searches for a route based on the path provided
func (c *Cage) getRoute(path string) (*Route, error) {
	cageLevel, routeSegment, err := c.cageLevel(path, false)
	if err != nil {
		return nil, err
	}

	route, found, err := cageLevel.route(routeSegment)
	if err != nil {
		return nil, err
	}
	if found {
		return route, nil
	}

	return nil, ErrRouteNotFound
}

// newRoute creates a new route, creating new cages if necessary
func (c *Cage) newRoute(route *Route) error {
	cageLevel, routeSegment, err := c.cageLevel(route.Path, true)
	if err != nil {
		return err
	}

	if strings.Contains(routeSegment, "*") {
		cageLevel.wildCardRoutesRWMu.Lock()
		defer cageLevel.wildCardRoutesRWMu.Unlock()
		cageLevel.WildCardRoutes[routeSegment] = route
	} else {
		cageLevel.routesRWMu.Lock()
		defer cageLevel.routesRWMu.Unlock()
		cageLevel.Routes[routeSegment] = route
	}

	return nil
}

// deleteRoute deletes a route based on the path provided
func (c *Cage) deleteRoute(route *Route) error {
	cageLevel, routeSegment, err := c.cageLevel(route.Path, true)
	if err != nil {
		return err
	}

	if strings.Contains(routeSegment, "*") {
		cageLevel.wildCardRoutesRWMu.Lock()
		defer cageLevel.wildCardRoutesRWMu.Unlock()
		delete(cageLevel.WildCardRoutes, routeSegment)
	} else {
		cageLevel.routesRWMu.Lock()
		defer cageLevel.routesRWMu.Unlock()
		delete(cageLevel.Routes, routeSegment)
	}

	return nil
}

// route searches for a route based on the route segment provided
func (cl *CageLevel) route(routeSegment string) (route *Route, found bool, err error) {
	// First check for an exact match
	cl.routesRWMu.Lock()
	if route, found = cl.Routes[routeSegment]; found {
		defer cl.routesRWMu.Unlock()
		return route, true, nil
	}
	cl.routesRWMu.Unlock()

	// if not, look for wildcard routes
	cl.wildCardRoutesRWMu.Lock()
	defer cl.wildCardRoutesRWMu.Unlock()
	for wildCardPattern, route := range cl.WildCardRoutes {
		match, err := filepath.Match(wildCardPattern, routeSegment)
		if err != nil {
			return nil, false, newDynamicError(ErrInvalidPath, err.Error())
		}
		if match {
			return route, true, nil
		}
	}

	return nil, false, nil
}

// subCageLevel searches for a sub cage level based on the segment provided
// if createMode is true, it will create the cage level if it doesn't exist
func (cl *CageLevel) subCageLevel(subCageSegment string, createMode bool) (cageLevel *CageLevel, found bool, err error) {
	// First check for an exact match
	cl.subCagesRWMu.RLock()
	if cageLevel, exists := cl.SubCages[subCageSegment]; exists {
		defer cl.subCagesRWMu.RUnlock()
		return cageLevel, true, nil
	}
	cl.subCagesRWMu.RUnlock()

	// if not, look for wildcard cages
	cl.wildCardSubCagesRWMu.RLock()
	for wildCardPattern, cageLevel := range cl.WildCardSubCages {
		match, err := filepath.Match(wildCardPattern, subCageSegment)
		if err != nil {
			cl.wildCardSubCagesRWMu.RUnlock()
			return nil, false, newDynamicError(ErrInvalidPath, err.Error())
		}
		if match {
			cl.wildCardSubCagesRWMu.RUnlock()
			return cageLevel, true, nil
		}
	}
	cl.wildCardSubCagesRWMu.RUnlock()

	// We didn't find a match, so we'll create a new cage level if we're in create mode
	if createMode {
		newCage := newCageLevel()
		if strings.Contains(subCageSegment, "*") {
			cl.wildCardSubCagesRWMu.Lock()
			defer cl.wildCardSubCagesRWMu.Unlock()
			cl.WildCardSubCages[subCageSegment] = newCage
		} else {
			cl.subCagesRWMu.Lock()
			defer cl.subCagesRWMu.Unlock()
			cl.SubCages[subCageSegment] = newCage
		}
		return newCage, true, nil
	}

	return nil, false, nil
}

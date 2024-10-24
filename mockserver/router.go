package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

// Route holds information about the mock route configuration
type Route struct {
	// Method is the HTTP method to match
	Method string `json:"method"`
	// Path is the URL path to match
	Path string `json:"path"`
	// Response is the JSON response to return when called
	Response string `json:"response"`
	// StatusCode is the HTTP status code to return when called
	StatusCode int `json:"status_code"`
	// ContentType is the Content-Type header to use when called
	ContentType string `json:"content_type"`
}

var (
	routes   = make(map[string]Route) // Store routes based on "Method:Path" keys
	routesMu sync.RWMutex             // Protects access to the routes map
)

// RegisterRouteHandler handles the dynamic route registration.
func RegisterRouteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var route Route
	if err := json.NewDecoder(r.Body).Decode(&route); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if route.Method == "" || route.Path == "" {
		http.Error(w, "Method and Path are required", http.StatusBadRequest)
		return
	}

	routesMu.Lock()
	routes[route.Method+":"+route.Path] = route
	routesMu.Unlock()

	w.WriteHeader(http.StatusCreated)
	log.Info().Str("Path", route.Path).Str("Method", route.Method).Msg("Route registered")
}

// DynamicHandler handles all incoming requests and responds based on the registered routes.
func DynamicHandler(w http.ResponseWriter, r *http.Request) {
	routesMu.RLock()
	route, exists := routes[r.Method+":"+r.URL.Path]
	routesMu.RUnlock()

	if !exists {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", route.ContentType)
	w.WriteHeader(route.StatusCode)
	_, err := io.WriteString(w, route.Response)
	if err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
	}
}

// Load loads all registered routes from a file.
func Load() error {
	if _, err := os.Stat(config.SaveFile); os.IsNotExist(err) {
		log.Debug().Str("Save File", config.SaveFile).Msg("No routes to load")
		return nil
	}

	start := time.Now()
	log.Debug().Str("Save File", config.SaveFile).Msg("Loading routes")

	routesMu.Lock()
	defer routesMu.Unlock()

	data, err := os.ReadFile(config.SaveFile)
	if err != nil {
		return fmt.Errorf("failed to read routes from file: %w", err)
	}

	if err = json.Unmarshal(data, &routes); err != nil {
		return fmt.Errorf("failed to unmarshal routes: %w", err)
	}

	log.Debug().Str("Save File", config.SaveFile).Int("Number", len(routes)).Str("Duration", time.Since(start).String()).Msg("Routes loaded")
	return nil
}

// Save saves all registered routes to a file.
func Save() error {
	start := time.Now()
	log.Debug().Str("Save File", config.SaveFile).Msg("Saving routes")

	routesMu.Lock()
	defer routesMu.Unlock()

	jsonData, err := json.MarshalIndent(routes, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal routes: %w", err)
	}

	if err = os.WriteFile(config.SaveFile, jsonData, 0644); err != nil { //nolint:gosec
		return fmt.Errorf("failed to write routes to file: %w", err)
	}

	log.Debug().Str("Save File", config.SaveFile).Str("Duration", time.Since(start).String()).Msg("Routes saved")
	return nil
}

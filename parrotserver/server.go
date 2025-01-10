package parrotserver

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Route holds information about the mock route configuration
type Route struct {
	// Method is the HTTP method to match
	Method string `json:"method"`
	// Path is the URL path to match
	Path string `json:"path"`
	// Response is the static JSON response to return when called
	Response any `json:"response"`
	// Handler is the dynamic handler function to use when called
	Handler http.HandlerFunc `json:"-"`
	// StatusCode is the HTTP status code to return when called
	StatusCode int `json:"status_code"`
	// ContentType is the Content-Type header to return the response with
	ContentType string `json:"content_type"`
}

// ParrotServer is a mock HTTP server that can register and respond to dynamic routes
type ParrotServer struct {
	port     int
	saveFile string
	l        zerolog.Logger

	server   *http.Server
	routes   map[string]Route // Store routes based on "Method:Path" keys
	routesMu sync.RWMutex
}

// ParrotServerOption defines functional options for configuring the ParrotServer
type ParrotServerOption func(*ParrotServer) error

// WithPort sets the port for the ParrotServer to run on
func WithPort(port int) ParrotServerOption {
	return func(s *ParrotServer) error {
		if port <= 0 {
			return fmt.Errorf("invalid port: %d", port)
		}
		s.port = port
		return nil
	}
}

func WithLogLevel(logLevel string) ParrotServerOption {
	return func(s *ParrotServer) error {
		level, err := zerolog.ParseLevel(logLevel)
		if err != nil {
			return fmt.Errorf("invalid log level: %s", logLevel)
		}
		s.l = s.l.Level(level)
		return nil
	}
}

// WithSaveFile sets the file to save the routes to
func WithSaveFile(saveFile string) ParrotServerOption {
	return func(s *ParrotServer) error {
		if saveFile == "" {
			return fmt.Errorf("invalid save file name: %s", saveFile)
		}
		s.saveFile = saveFile
		return nil
	}
}

// WithRoutes sets the initial routes for the ParrotServer
func WithRoutes(routes map[string]Route) ParrotServerOption {
	return func(s *ParrotServer) error {
		for k, v := range routes {
			if v.Path == "" || v.Path == "/" || v.Path == "/register" {
				return fmt.Errorf("invalid route path: %s", v.Path)
			}
			if v.Method == "" {
				return fmt.Errorf("invalid route method: %s", v.Method)
			}
			s.routes[k] = v
		}
		return nil
	}
}

// New creates a new HTTP server with the dynamic route handling.
func New(options ...ParrotServerOption) (*ParrotServer, error) {
	p := &ParrotServer{
		port:     8080,
		saveFile: "routes.json",
		l: zerolog.New(os.Stderr).With().
			Timestamp().Logger().Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339Nano}),

		routes:   make(map[string]Route),
		routesMu: sync.RWMutex{},
	}

	for _, option := range options {
		if err := option(p); err != nil {
			return nil, err
		}
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/register", p.registerRouteHandler)
	mux.HandleFunc("/", p.dynamicHandler)

	p.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", p.port),
		Handler: mux,
	}

	log.Info().Int("port", p.port).Str("saveFile", p.saveFile).Msg("Parrot server started")
	return p, nil
}

// registerRouteHandler handles the dynamic route registration.
func (p *ParrotServer) registerRouteHandler(w http.ResponseWriter, r *http.Request) {
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

	p.routesMu.Lock()
	p.routes[route.Method+":"+route.Path] = route
	p.routesMu.Unlock()

	w.WriteHeader(http.StatusCreated)
	log.Info().Str("Path", route.Path).Str("Method", route.Method).Msg("Route registered")
}

// dynamicHandler handles all incoming requests and responds based on the registered routes.
func (p *ParrotServer) dynamicHandler(w http.ResponseWriter, r *http.Request) {
	p.routesMu.RLock()
	route, exists := p.routes[r.Method+":"+r.URL.Path]
	p.routesMu.RUnlock()

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

// load loads all registered routes from a file.
func load() error {
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

// save saves all registered routes to a file.
func save() error {
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

func (p *ParrotServer) RegisterRoute(route Route) {
	routesMu.Lock()
	routes[route.Method+":"+route.Path] = route
	routesMu.Unlock()
}

func CallRoute(method, path string) (int, string, string) {
	routesMu.RLock()
	route, exists := routes[method+":"+path]
	routesMu.RUnlock()

	if !exists {
		return http.StatusNotFound, "", ""
	}

	return route.StatusCode, route.ContentType, route.Response
}

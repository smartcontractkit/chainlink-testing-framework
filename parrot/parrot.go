package parrot

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/rs/zerolog"
)

// Route holds information about the mock route configuration
type Route struct {
	// Method is the HTTP method to match
	Method string `json:"Method"`
	// Path is the URL path to match
	Path string `json:"Path"`
	// Handler is the dynamic handler function to use when called
	// Can only be set upon creation of the server
	Handler http.HandlerFunc `json:"-"`
	// RawResponseBody is the static, raw string response to return when called
	RawResponseBody string `json:"raw_response_body"`
	// ResponseBody will be marshalled to JSON and returned when called
	ResponseBody any `json:"response_body"`
	// ResponseStatusCode is the HTTP status code to return when called
	ResponseStatusCode int `json:"response_status_code"`
	// ResponseContentType is the Content-Type header to return the response with
	ResponseContentType string `json:"response_content_type"`
}

// Server is a mock HTTP server that can register and respond to dynamic routes
type Server struct {
	port     int
	host     string
	address  string
	saveFile string
	log      zerolog.Logger

	server   *http.Server
	routes   map[string]*Route // Store routes based on "Method:Path" keys
	routesMu sync.RWMutex
}

// ServerOption defines functional options for configuring the ParrotServer
type ServerOption func(*Server) error

// WithPort sets the port for the ParrotServer to run on
func WithPort(port int) ServerOption {
	return func(s *Server) error {
		if port == 0 {
			s.log.Debug().Msg("Configuring Parrot: No port specified, using random port")
		} else if port < 0 || port > 65535 {
			return fmt.Errorf("invalid port: %d", port)
		}
		s.port = port
		s.log.Debug().Int("port", port).Msg("Configuring Parrot: Setting port")
		return nil
	}
}

func WithLogLevel(level zerolog.Level) ServerOption {
	return func(s *Server) error {
		s.log = s.log.Level(level)
		s.log.Debug().Str("log level", level.String()).Msg("Configuring Parrot: Setting log level")
		return nil
	}
}

// WithLogger sets the logger for the ParrotServer
func WithLogger(l zerolog.Logger) ServerOption {
	return func(s *Server) error {
		s.log = l
		s.log.Debug().Msg("Configuring Parrot: Setting custom logger")
		return nil
	}
}

// WithJSONLogs sets the logger to output JSON logs
func WithJSONLogs() ServerOption {
	return func(s *Server) error {
		s.log = s.log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339Nano})
		s.log.Debug().Msg("Configuring Parrot: Setting log output to JSON")
		return nil
	}
}

// WithSaveFile sets the file to save the routes to
func WithSaveFile(saveFile string) ServerOption {
	return func(s *Server) error {
		if saveFile == "" {
			return fmt.Errorf("invalid save file name: %s", saveFile)
		}
		s.saveFile = saveFile
		s.log.Debug().Str("file", saveFile).Msg("Configuring Parrot: Setting save file")
		return nil
	}
}

// WithRoutes sets the initial routes for the Parrot
func WithRoutes(routes []*Route) ServerOption {
	return func(s *Server) error {
		for _, route := range routes {
			if err := s.Register(route); err != nil {
				return fmt.Errorf("failed to register route: %w", err)
			}
			s.log.Debug().Str("Path", route.Path).Str("Method", route.Method).Msg("Configuring Parrot: Pre-registered route")
		}
		return nil
	}
}

// Wake creates a new Parrot server with dynamic route handling
func Wake(options ...ServerOption) (*Server, error) {
	p := &Server{
		port:     0,
		saveFile: "save.json",
		log: zerolog.New(os.Stderr).Level(zerolog.InfoLevel).With().
			Timestamp().Logger().Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339Nano}),

		routes:   make(map[string]*Route),
		routesMu: sync.RWMutex{},
	}

	for _, option := range options {
		if err := option(p); err != nil {
			return nil, err
		}
	}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", p.port))
	if err != nil {
		return nil, fmt.Errorf("failed to start listener: %w", err)
	}
	host, port, err := net.SplitHostPort(listener.Addr().String())
	if err != nil {
		return nil, fmt.Errorf("failed to split host and port: %w", err)
	}
	p.host = host
	p.address = listener.Addr().String()
	p.port, err = strconv.Atoi(port)
	if err != nil {
		return nil, fmt.Errorf("failed to parse port: %w", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/register", p.registerRouteHandler)
	mux.HandleFunc("/", p.dynamicHandler)

	p.server = &http.Server{
		ReadHeaderTimeout: 5 * time.Second,
		Addr:              fmt.Sprintf(":%d", p.port),
		Handler:           mux,
	}

	if err = p.load(); err != nil {
		return nil, fmt.Errorf("failed to load saved routes: %w", err)
	}

	go func() {
		defer func() {
			if err = p.save(); err != nil {
				p.log.Error().Err(err).Msg("Failed to save routes")
			}
		}()

		p.log.Info().Int("port", p.Port()).Str("address", p.address).Str("save file", p.saveFile).Msg("Parrot awake and ready to squawk")
		if err = p.server.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
			p.log.Fatal().Err(err).Msg("Error while running server")
		}
	}()

	return p, nil
}

// Shutdown gracefully shuts down the parrot server
func (p *Server) Shutdown(ctx context.Context) error {
	p.log.Info().Msg("Putting cloth over the parrot's cage...")
	return p.server.Shutdown(ctx)
}

// Host returns the host the parrot is running on
func (p *Server) Host() string {
	return p.host
}

// Port returns the port the parrot is running on
func (p *Server) Port() int {
	return p.port
}

// Address returns the address the parrot is running on
func (p *Server) Address() string {
	return p.address
}

// Register adds a new route to the parrot
func (p *Server) Register(route *Route) error {
	if route.Path == "" || route.Path == "/" || route.Path == "/register" {
		return fmt.Errorf("invalid route path: %s", route.Path)
	}
	if route.Method == "" {
		return fmt.Errorf("invalid route method: %s", route.Method)
	}
	if route.Handler == nil && route.ResponseBody == nil && route.RawResponseBody == "" {
		return fmt.Errorf("route must have a handler or response body")
	}
	if route.Handler != nil && (route.ResponseBody != nil || route.RawResponseBody != "") {
		return fmt.Errorf("route cannot have both a handler and response body")
	}
	if route.ResponseBody != nil && route.RawResponseBody != "" {
		return fmt.Errorf("route cannot have both a response body and raw response body")
	}
	if route.ResponseBody != nil {
		if _, err := json.Marshal(route.ResponseBody); err != nil {
			return fmt.Errorf("response body is unable to be marshalled into JSON: %w", err)
		}
	}

	p.routesMu.Lock()
	defer p.routesMu.Unlock()
	p.routes[route.Method+":"+route.Path] = route

	return nil
}

// Routes returns all registered routes
func (p *Server) Routes() map[string]*Route {
	p.routesMu.RLock()
	defer p.routesMu.RUnlock()
	return p.routes
}

// Unregister removes a route from the parrot
func (p *Server) Unregister(method, path string) {
	p.routesMu.Lock()
	defer p.routesMu.Unlock()
	delete(p.routes, method+":"+path)
}

// Call makes a request to the parrot server
func (p *Server) Call(method, path string) (*http.Response, error) {
	req, err := http.NewRequest(method, "http://"+filepath.Join(p.Address(), path), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	client := &http.Client{}
	return client.Do(req)
}

// registerRouteHandler handles the dynamic route registration.
func (p *Server) registerRouteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		p.log.Trace().Str("Method", r.Method).Msg("Invalid method")
		return
	}

	var route *Route
	if err := json.NewDecoder(r.Body).Decode(&route); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if route.Method == "" || route.Path == "" {
		err := errors.New("Method and path are required")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := p.Register(route)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		p.log.Trace().Err(err).Msg("Failed to register route")
		return
	}

	w.WriteHeader(http.StatusCreated)
	p.log.Info().Str("Path", route.Path).Str("Method", route.Method).Msg("Route registered")
}

// dynamicHandler handles all incoming requests and responds based on the registered routes.
func (p *Server) dynamicHandler(w http.ResponseWriter, r *http.Request) {
	p.routesMu.RLock()
	route, exists := p.routes[r.Method+":"+r.URL.Path]
	p.routesMu.RUnlock()

	if !exists {
		http.NotFound(w, r)
		p.log.Trace().Str("Remote Addr", r.RemoteAddr).Str("Path", r.URL.Path).Str("Method", r.Method).Msg("Route not found")
		return
	}

	if route.ResponseContentType != "" {
		w.Header().Set("Content-Type", route.ResponseContentType)
	}
	w.WriteHeader(route.ResponseStatusCode)

	if route.Handler != nil {
		p.log.Trace().Str("Remote Addr", r.RemoteAddr).Str("Path", r.URL.Path).Str("Method", r.Method).Msg("Calling route handler")
		route.Handler(w, r)
	} else if route.RawResponseBody != "" {
		if route.ResponseContentType == "" {
			w.Header().Set("Content-Type", "text/plain")
		}
		if _, err := w.Write([]byte(route.RawResponseBody)); err != nil {
			p.log.Trace().Err(err).Str("Remote Addr", r.RemoteAddr).Str("Path", r.URL.Path).Str("Method", r.Method).Msg("Failed to write response")
			http.Error(w, "Failed to write response", http.StatusInternalServerError)
			return
		}
		p.log.Trace().
			Str("Remote Addr", r.RemoteAddr).
			Str("Response", route.RawResponseBody).
			Str("Path", r.URL.Path).
			Str("Method", r.Method).
			Msg("Returned raw response")
	} else if route.ResponseBody != nil {
		if route.ResponseContentType == "" {
			w.Header().Set("Content-Type", "application/json")
		}
		rawJSON, err := json.Marshal(route.ResponseBody)
		if err != nil {
			p.log.Trace().Err(err).
				Str("Remote Addr", r.RemoteAddr).
				Str("Path", r.URL.Path).
				Str("Method", r.Method).
				Msg("Failed to marshal JSON response")
			http.Error(w, "Failed to marshal response into json", http.StatusInternalServerError)
			return
		}
		if _, err = w.Write(rawJSON); err != nil {
			p.log.Trace().Err(err).
				RawJSON("Response", rawJSON).
				Str("Remote Addr", r.RemoteAddr).
				Str("Path", r.URL.Path).
				Str("Method", r.Method).
				Msg("Failed to write response")
			http.Error(w, "Failed to write JSON response", http.StatusInternalServerError)
			return
		}
		p.log.Trace().
			Str("Remote Addr", r.RemoteAddr).
			RawJSON("Response", rawJSON).
			Str("Path", r.URL.Path).
			Str("Method", r.Method).
			Msg("Returned JSON response")
	} else {
		p.log.Trace().Str("Remote Addr", r.RemoteAddr).Str("Path", r.URL.Path).Str("Method", r.Method).Msg("Route has no response")
	}
}

// load loads all registered routes from a file.
func (p *Server) load() error {
	if _, err := os.Stat(p.saveFile); os.IsNotExist(err) {
		p.log.Trace().Str("file", p.saveFile).Msg("No routes to load")
		return nil
	}

	p.log.Debug().Str("file", p.saveFile).Msg("Loading routes")

	data, err := os.ReadFile(p.saveFile)
	if err != nil {
		return fmt.Errorf("failed to read routes from file: %w", err)
	}
	if len(data) == 0 {
		p.log.Trace().Str("file", p.saveFile).Msg("No routes to load")
		return nil
	}

	p.routesMu.Lock()
	defer p.routesMu.Unlock()

	if err = json.Unmarshal(data, &p.routes); err != nil {
		return fmt.Errorf("failed to unmarshal routes: %w", err)
	}

	p.log.Info().Str("file", p.saveFile).Int("number", len(p.routes)).Msg("Loaded routes")
	return nil
}

// save saves all registered routes to a file.
func (p *Server) save() error {
	if len(p.routes) == 0 {
		p.log.Debug().Msg("No routes to save")
		return nil
	}
	p.log.Trace().Str("file", p.saveFile).Msg("Saving routes")

	p.routesMu.RLock()
	defer p.routesMu.RUnlock()

	jsonData, err := json.Marshal(p.routes)
	if err != nil {
		return fmt.Errorf("failed to marshal routes: %w", err)
	}

	if err = os.WriteFile(p.saveFile, jsonData, 0644); err != nil { //nolint:gosec
		return fmt.Errorf("failed to write routes to file: %w", err)
	}

	p.log.Trace().Str("file", p.saveFile).Msg("Saved routes")
	return nil
}

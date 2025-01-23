package parrot

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
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
}

// ID returns the unique identifier for the route
func (r *Route) ID() string {
	return r.Method + ":" + r.Path
}

// RouteRequest is the request body for querying the server on a specific route
type RouteRequest struct {
	ID string `json:"id"`
}

// Server is a mock HTTP server that can register and respond to dynamic routes
type Server struct {
	port    int
	host    string
	address string

	saveFileName    string
	useCustomLogger bool
	logFileName     string
	logFile         *os.File
	logLevel        zerolog.Level
	jsonLogs        bool
	log             zerolog.Logger

	server   *http.Server
	routes   map[string]*Route // Store routes based on "Method:Path" keys
	routesMu sync.RWMutex

	recorderHooks []string
	recordersMu   sync.RWMutex
}

// ServerOption defines functional options for configuring the ParrotServer
type ServerOption func(*Server) error

// WithPort sets the port for the ParrotServer to run on
func WithPort(port int) ServerOption {
	return func(s *Server) error {
		if port < 0 || port > 65535 {
			return fmt.Errorf("invalid port: %d", port)
		}
		s.port = port
		return nil
	}
}

// WithLogLevel sets the visible log level of the default logger
func WithLogLevel(level zerolog.Level) ServerOption {
	return func(s *Server) error {
		s.logLevel = level
		return nil
	}
}

// WithLogger sets the logger for the ParrotServer
func WithLogger(l zerolog.Logger) ServerOption {
	return func(s *Server) error {
		s.log = l
		s.useCustomLogger = true
		return nil
	}
}

// WithJSONLogs sets the logger to output JSON logs
func WithJSONLogs() ServerOption {
	return func(s *Server) error {
		s.jsonLogs = true
		return nil
	}
}

// WithSaveFile sets the file to save the routes to
func WithSaveFile(saveFile string) ServerOption {
	return func(s *Server) error {
		if saveFile == "" {
			return fmt.Errorf("invalid save file name: %s", saveFile)
		}
		s.saveFileName = saveFile
		return nil
	}
}

// WithLogFile sets the file to save the logs to
func WithLogFile(logFile string) ServerOption {
	return func(s *Server) error {
		if logFile == "" {
			return fmt.Errorf("invalid log file name: %s", logFile)
		}
		s.logFileName = logFile
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
		}
		return nil
	}
}

// Wake creates a new Parrot server with dynamic route handling
func Wake(options ...ServerOption) (*Server, error) {
	p := &Server{
		port:         0,
		saveFileName: "parrot_save.json",
		logLevel:     zerolog.InfoLevel,
		logFileName:  "parrot.log",

		routes:   make(map[string]*Route),
		routesMu: sync.RWMutex{},
	}

	for _, option := range options {
		if err := option(p); err != nil {
			return nil, err
		}
	}

	var err error
	p.logFile, err = os.Create(p.logFileName)
	if err != nil {
		return nil, fmt.Errorf("failed to create log file: %w", err)
	}

	if !p.useCustomLogger { // Build default logger
		var writers []io.Writer

		if p.jsonLogs {
			writers = append(writers, os.Stderr)
		} else {
			consoleOut := zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "2006-01-02T15:04:05.000"}
			writers = append(writers, consoleOut)
		}

		if p.logFile != nil {
			writers = append(writers, p.logFile)
		}

		multiWriter := zerolog.MultiLevelWriter(writers...)
		p.log = zerolog.New(multiWriter).Level(p.logLevel).With().Timestamp().Logger()
	}

	listener, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", p.port))
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
	mux.HandleFunc("/record", p.recordHandler)
	mux.HandleFunc("/", p.dynamicHandler)

	p.server = &http.Server{
		ReadHeaderTimeout: 5 * time.Second,
		Addr:              listener.Addr().String(),
		Handler:           p.loggingMiddleware(mux),
	}

	if err = p.load(); err != nil {
		return nil, fmt.Errorf("failed to load saved routes: %w", err)
	}

	go p.run(listener)

	return p, nil
}

func (p *Server) run(listener net.Listener) {
	defer func() {
		if err := p.save(); err != nil {
			p.log.Error().Err(err).Msg("Failed to save routes")
		}
		if err := p.logFile.Close(); err != nil {
			p.log.Error().Err(err).Msg("Failed to close log file")
		}
	}()

	p.log.Info().Int("Port", p.Port()).Str("Address", p.address).Msg("Parrot awake and ready to squawk")
	p.log.Debug().Str("Save File", p.saveFileName).Str("Log File", p.logFileName).Msg("Configuration")
	if err := p.server.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
		p.log.Fatal().Err(err).Msg("Error while running server")
	}
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
	if route == nil {
		return ErrNilRoute
	}
	if !isValidPath(route.Path) {
		return newDynamicError(ErrInvalidPath, fmt.Sprintf("'%s'", route.Path))
	}
	if _, err := url.Parse(route.Path); err != nil {
		return newDynamicError(ErrInvalidPath, fmt.Sprintf("%s: '%s'", err.Error(), route.Path))
	}
	if route.Method == "" {
		return ErrNoMethod
	}
	if route.Handler == nil && route.ResponseBody == nil && route.RawResponseBody == "" {
		return ErrNoResponse
	}
	if route.Handler != nil && (route.ResponseBody != nil || route.RawResponseBody != "") {
		return newDynamicError(ErrOnlyOneResponse, "handler and another response type provided")
	}
	if route.ResponseBody != nil && route.RawResponseBody != "" {
		return ErrOnlyOneResponse
	}
	if route.ResponseBody != nil {
		if _, err := json.Marshal(route.ResponseBody); err != nil {
			return newDynamicError(ErrResponseMarshal, err.Error())
		}
	}

	p.routesMu.Lock()
	defer p.routesMu.Unlock()
	p.routes[route.ID()] = route
	p.log.Info().
		Str("Route ID", route.ID()).
		Str("Path", route.Path).
		Str("Method", route.Method).
		Msg("Route registered")

	return nil
}

// registerRouteHandler handles the dynamic route registration.
func (p *Server) registerRouteHandler(w http.ResponseWriter, r *http.Request) {
	registerLogger := zerolog.Ctx(r.Context())
	if r.Method == http.MethodDelete {
		var routeRequest *RouteRequest
		if err := json.NewDecoder(r.Body).Decode(&routeRequest); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			registerLogger.Debug().Err(err).Msg("Failed to decode request body")
			return
		}
		defer r.Body.Close()

		if routeRequest.ID == "" {
			http.Error(w, "Route ID required", http.StatusBadRequest)
			registerLogger.Debug().Msg("No Route ID provided")
			return
		}

		err := p.Unregister(routeRequest.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			registerLogger.Debug().Err(err).Msg("Failed to unregister route")
			return
		}

		w.WriteHeader(http.StatusNoContent)
		registerLogger.Info().
			Str("Route ID", routeRequest.ID).
			Msg("Route unregistered")
	} else if r.Method == http.MethodPost {
		var route *Route
		if err := json.NewDecoder(r.Body).Decode(&route); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			registerLogger.Debug().Err(err).Msg("Failed to decode request body")
			return
		}
		defer r.Body.Close()

		if route.Method == "" || route.Path == "" {
			err := errors.New("Method and path are required")
			http.Error(w, err.Error(), http.StatusBadRequest)
			registerLogger.Debug().Err(err).Msg("Method and path are required")
			return
		}

		err := p.Register(route)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			registerLogger.Debug().Err(err).Msg("Failed to register route")
			return
		}

		w.WriteHeader(http.StatusCreated)
	} else {
		http.Error(w, "Invalid method, only use POST or DELETE", http.StatusMethodNotAllowed)
		registerLogger.Debug().Msg("Invalid method")
		return
	}
}

// Record registers a new recorder with the parrot. All incoming requests to the parrot will be sent to the recorder.
func (p *Server) Record(recorder *Recorder) error {
	p.recordersMu.Lock()
	defer p.recordersMu.Unlock()
	if recorder == nil {
		return ErrNilRecorder
	}
	if recorder.URL == "" {
		return ErrNoRecorderURL
	}
	_, err := url.Parse(recorder.URL)
	if err != nil {
		return fmt.Errorf("failed to parse recorder URL: %w", err)
	}
	p.recorderHooks = append(p.recorderHooks, recorder.URL)
	return nil
}

func (p *Server) recordHandler(w http.ResponseWriter, r *http.Request) {
	recordLogger := zerolog.Ctx(r.Context())
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method, only use POST or DELETE", http.StatusMethodNotAllowed)
		recordLogger.Debug().Msg("Invalid method")
		return
	}

	var recorder *Recorder
	if err := json.NewDecoder(r.Body).Decode(&recorder); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		recordLogger.Err(err).Msg("Failed to decode request body")
		return
	}

	err := p.Record(recorder)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		recordLogger.Debug().Err(err).Msg("Failed to add recorder")
		return
	}

	w.WriteHeader(http.StatusCreated)
	recordLogger.Info().Str("Recorder URL", recorder.URL).Msg("Recorder added")
}

// Unregister removes a route from the parrot
func (p *Server) Unregister(routeID string) error {
	p.routesMu.RLock()
	_, exists := p.routes[routeID]
	p.routesMu.RUnlock()

	if !exists {
		return newDynamicError(ErrRouteNotFound, routeID)
	}
	p.routesMu.Lock()
	defer p.routesMu.Unlock()
	delete(p.routes, routeID)
	return nil
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

// dynamicHandler handles all incoming requests and responds based on the registered routes.
func (p *Server) dynamicHandler(w http.ResponseWriter, r *http.Request) {
	p.routesMu.RLock()
	route, exists := p.routes[r.Method+":"+r.URL.Path]
	p.routesMu.RUnlock()

	dynamicLogger := zerolog.Ctx(r.Context())
	if !exists {
		http.NotFound(w, r)
		dynamicLogger.Debug().Msg("Route not found")
		return
	}

	requestID := uuid.New().String()[0:8]
	dynamicLogger.UpdateContext(func(c zerolog.Context) zerolog.Context {
		return c.Str("Request ID", requestID).Str("Route ID", route.ID())
	})

	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		dynamicLogger.Debug().
			Err(err).
			Msg("Failed to read request body")
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	routeCall := &RouteCall{
		RouteID: r.Method + ":" + r.URL.Path,
		Request: &RouteCallRequest{
			Method: r.Method,
			URL:    r.URL,
			Header: r.Header,
			Body:   requestBody,
		},
	}
	recordingWriter := newResponseWriterRecorder(w)

	defer func() {
		res := recordingWriter.Result()
		resBody, err := io.ReadAll(res.Body)
		if err != nil {
			dynamicLogger.Debug().Err(err).Msg("Failed to read response body")
			http.Error(w, "Failed to read response body", http.StatusInternalServerError)
			return
		}

		routeCall.Response = &RouteCallResponse{
			StatusCode: res.StatusCode,
			Header:     res.Header,
			Body:       resBody,
		}
		p.sendToRecorders(routeCall)
	}()

	// Let the custom handler take over if it exists
	if route.Handler != nil {
		dynamicLogger.Debug().Msg("Calling route handler")
		route.Handler(recordingWriter, r)
		return
	}

	if route.RawResponseBody != "" {
		if _, err := w.Write([]byte(route.RawResponseBody)); err != nil {
			dynamicLogger.Debug().Err(err).Msg("Failed to write response")
			http.Error(recordingWriter, "Failed to write response", http.StatusInternalServerError)
			return
		}
		dynamicLogger.Debug().
			Str("Response", route.RawResponseBody).
			Msg("Returned raw response")
		recordingWriter.WriteHeader(route.ResponseStatusCode)
		return
	}

	if route.ResponseBody != nil {
		rawJSON, err := json.Marshal(route.ResponseBody)
		if err != nil {
			dynamicLogger.Debug().Err(err).Msg("Failed to marshal JSON response")
			http.Error(recordingWriter, "Failed to marshal response into json", http.StatusInternalServerError)
			return
		}
		if _, err = w.Write(rawJSON); err != nil {
			dynamicLogger.Debug().Err(err).
				RawJSON("Response", rawJSON).
				Msg("Failed to write response")
			http.Error(recordingWriter, "Failed to write JSON response", http.StatusInternalServerError)
			return
		}
		dynamicLogger.Debug().
			RawJSON("Response", rawJSON).
			Msg("Returned JSON response")
		recordingWriter.WriteHeader(route.ResponseStatusCode)
		return
	}

	dynamicLogger.Error().Msg("Route has no response")
	http.Error(recordingWriter, "Route has no response", http.StatusInternalServerError)
}

// load loads all registered routes from a file.
func (p *Server) load() error {
	if _, err := os.Stat(p.saveFileName); os.IsNotExist(err) {
		p.log.Trace().Str("file", p.saveFileName).Msg("No routes to load")
		return nil
	}

	p.log.Debug().Str("file", p.saveFileName).Msg("Loading routes")

	data, err := os.ReadFile(p.saveFileName)
	if err != nil {
		return fmt.Errorf("failed to read routes from file: %w", err)
	}
	if len(data) == 0 {
		p.log.Trace().Str("file", p.saveFileName).Msg("No routes to load")
		return nil
	}

	p.routesMu.Lock()
	defer p.routesMu.Unlock()

	if err = json.Unmarshal(data, &p.routes); err != nil {
		return fmt.Errorf("failed to unmarshal routes: %w", err)
	}

	p.log.Info().Str("file", p.saveFileName).Int("number", len(p.routes)).Msg("Loaded routes")
	return nil
}

// save saves all registered routes to a file.
func (p *Server) save() error {
	if len(p.routes) == 0 {
		p.log.Debug().Msg("No routes to save")
		return nil
	}
	p.log.Trace().Str("file", p.saveFileName).Msg("Saving routes")

	p.routesMu.RLock()
	defer p.routesMu.RUnlock()

	jsonData, err := json.Marshal(p.routes)
	if err != nil {
		return fmt.Errorf("failed to marshal routes: %w", err)
	}

	if err = os.WriteFile(p.saveFileName, jsonData, 0644); err != nil { //nolint:gosec
		return fmt.Errorf("failed to write routes to file: %w", err)
	}

	p.log.Trace().Str("file", p.saveFileName).Msg("Saved routes")
	return nil
}

// sendToRecorders sends the route call to all registered recorders
func (p *Server) sendToRecorders(routeCall *RouteCall) {
	p.recordersMu.RLock()
	defer p.recordersMu.RUnlock()
	if len(p.recorderHooks) == 0 {
		return
	}

	client := resty.New()
	p.log.Trace().Strs("Recorders", p.recorderHooks).Str("Route ID", routeCall.RouteID).Msg("Sending route call to recorders")

	for _, hook := range p.recorderHooks {
		go func(hook string) {
			resp, err := client.R().SetBody(routeCall).Post(hook)
			if err != nil {
				p.log.Error().Err(err).Str("Recorder Hook", hook).Msg("Failed to send route call to recorder")
				return
			}
			if resp.IsError() {
				p.log.Error().
					Str("Recorder Hook", hook).
					Int("Code", resp.StatusCode()).
					Str("Response", resp.String()).
					Msg("Failed to send route call to recorder")
				return
			}
			p.log.Trace().Str("Route ID", routeCall.RouteID).Str("Recorder Hook", hook).Msg("Route call sent to recorder")
		}(hook)
	}
}

func (p *Server) loggingMiddleware(next http.Handler) http.Handler {
	h := hlog.NewHandler(p.log)

	accessHandler := hlog.AccessHandler(
		func(r *http.Request, status, size int, duration time.Duration) {
			hlog.FromRequest(r).Trace().
				Str("Method", r.Method).
				Stringer("URL", r.URL).
				Int("Status Code", status).
				Int("Response Size Bytes", size).
				Str("Duration", duration.String()).
				Str("Remote Addr", r.RemoteAddr).
				Msg("Handled request")
		},
	)

	return h(accessHandler(next))
}

var pathRegex = regexp.MustCompile(`^\/[a-zA-Z0-9\-._~%!$&'()*+,;=:@\/]*$`)

func isValidPath(path string) bool {
	switch path {
	case "", "/", "//", "/register", "/.", "/..":
		return false
	}
	if !strings.HasPrefix(path, "/") {
		return false
	}
	if strings.HasPrefix(path, "/register") {
		return false
	}
	if strings.HasPrefix(path, "/unregister") {
		return false
	}
	u, err := url.Parse(path)
	if err != nil || u.Path != path {
		return false
	}
	return pathRegex.MatchString(u.Path)
}

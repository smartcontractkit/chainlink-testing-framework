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

const (
	healthRoute = "/health"
	routesRoute = "/routes"
	recordRoute = "/record"
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

// Segment returns the last segment of the route path
func (r *Route) Segment() string {
	segments := strings.Split(r.Path, "/")
	if len(segments) == 0 {
		return ""
	}

	return segments[len(segments)-1]
}

// Server is a mock HTTP server that can register and respond to dynamic routes
type Server struct {
	port    int
	host    string
	address string

	client             *resty.Client
	shutDown           bool
	shutDownChan       chan struct{}
	shutDownOnce       sync.Once
	saveFileName       string
	useCustomLogger    bool
	logFileName        string
	logFile            *os.File
	logLevel           zerolog.Level
	jsonLogs           bool
	disableConsoleLogs bool
	log                zerolog.Logger

	server *http.Server
	cage   *cage // The root cage for the parrot that manages all dynamic routes

	recorderHooks map[string]struct{} // Store recorders based on URL keys to avoid duplicates
	recordersMu   sync.RWMutex
}

// SaveFile is the structure of the file to save and load parrot data from
type SaveFile struct {
	Routes    []*Route `json:"routes"`
	Recorders []string `json:"recorders"`
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

// WithRecorders pre-registers recorders with the ParrotServer
func WithRecorders(recorderURLs ...string) ServerOption {
	return func(s *Server) error {
		for _, url := range recorderURLs {
			if err := s.Record(url); err != nil {
				return fmt.Errorf("failed to register recorder: %w", err)
			}
		}
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

// DisableConsoleLogs disables logging to the console
func DisableConsoleLogs() ServerOption {
	return func(s *Server) error {
		s.disableConsoleLogs = true
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

		client:       resty.New(),
		shutDownChan: make(chan struct{}),

		cage: newCage(),

		recorderHooks: make(map[string]struct{}),
		recordersMu:   sync.RWMutex{},
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

		zerolog.TimeFieldFormat = "2006-01-02T15:04:05.000"
		if !p.disableConsoleLogs {
			if p.jsonLogs {
				writers = append(writers, os.Stderr)
			} else {
				consoleOut := zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "2006-01-02T15:04:05.000"}
				writers = append(writers, consoleOut)
			}
		}

		if p.logFile != nil {
			writers = append(writers, p.logFile)
		}

		multiWriter := zerolog.MultiLevelWriter(writers...)
		p.log = zerolog.New(multiWriter).Level(p.logLevel).With().Timestamp().Logger()
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
	// TODO: Add a route to enable registering recorders
	mux.HandleFunc(routesRoute, p.routeHandler)
	mux.HandleFunc(recordRoute, p.recordHandler)
	mux.HandleFunc(healthRoute, p.healthHandler)
	mux.HandleFunc("/", p.dynamicHandler)

	p.server = &http.Server{
		ReadHeaderTimeout: 5 * time.Second,
		Addr:              listener.Addr().String(),
		Handler:           p.loggingMiddleware(mux),
	}

	if err = p.load(); err != nil {
		return nil, fmt.Errorf("failed to load data from '%s': %w", p.saveFileName, err)
	}

	go p.run(listener)

	return p, nil
}

// run starts the parrot server
func (p *Server) run(listener net.Listener) {
	defer func() {
		p.log.Info().Msg("Putting cloth over the parrot's cage...")
		p.shutDown = true
		if err := p.save(); err != nil {
			p.log.Error().Err(err).Msg("Failed to save routes")
		}
		if err := p.logFile.Close(); err != nil {
			p.log.Error().Err(err).Msg("Failed to close log file")
		}
		p.shutDownOnce.Do(func() {
			close(p.shutDownChan)
		})
	}()

	p.log.Info().Str("Address", p.address).Msg("Parrot awake and ready to squawk")
	p.log.Debug().Str("Save File", p.saveFileName).Str("Log File", p.logFileName).Msg("Parrot configuration")
	if err := p.server.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
		p.log.Fatal().Err(err).Msg("Error while running server")
	}
}

// Shutdown gracefully shuts down the parrot server
func (p *Server) Shutdown(ctx context.Context) error {
	if p.shutDown {
		return ErrServerShutdown
	}
	return p.server.Shutdown(ctx)
}

// WaitShutdown blocks until the parrot server has shut down
func (p *Server) WaitShutdown() {
	<-p.shutDownChan
}

// Address returns the address the parrot is running on
func (p *Server) Address() string {
	return p.address
}

// Register adds a new route to the parrot
func (p *Server) Register(route *Route) error {
	if p.shutDown {
		return ErrServerShutdown
	}
	if route == nil {
		return ErrNilRoute
	}
	if !isValidPath(route.Path) {
		return newDynamicError(ErrInvalidPath, fmt.Sprintf("'%s'", route.Path))
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

	err := p.cage.newRoute(route)
	if err != nil {
		return err
	}
	p.log.Info().
		Str("Route ID", route.ID()).
		Str("Path", route.Path).
		Str("Method", route.Method).
		Msg("Route registered")

	return nil
}

// Record registers a new recorder with the parrot. All incoming requests to the parrot will be sent to the recorder.
func (p *Server) Record(recorderURL string) error {
	if p.shutDown {
		return ErrServerShutdown
	}

	p.recordersMu.Lock()
	defer p.recordersMu.Unlock()
	if recorderURL == "" {
		return ErrNoRecorderURL
	}
	_, err := url.ParseRequestURI(recorderURL)
	if err != nil {
		return ErrInvalidRecorderURL
	}
	p.recorderHooks[recorderURL] = struct{}{}
	return nil
}

// Recorders returns the URLs of all registered recorders
func (p *Server) Recorders() []string {
	if p.shutDown {
		return nil
	}

	p.recordersMu.RLock()
	defer p.recordersMu.RUnlock()
	recorders := make([]string, 0, len(p.recorderHooks))
	for recorder := range p.recorderHooks {
		recorders = append(recorders, recorder)
	}
	return recorders
}

// Delete removes a route from the parrot
func (p *Server) Delete(route *Route) error {
	if p.shutDown {
		return ErrServerShutdown
	}

	return p.cage.deleteRoute(route)
}

// Call makes a request to the parrot server
func (p *Server) Call(method, path string) (*resty.Response, error) {
	if p.shutDown {
		return nil, ErrServerShutdown
	}
	return p.client.R().Execute(method, "http://"+filepath.Join(p.Address(), path))
}

func (p *Server) Routes() []*Route {
	return p.cage.routes()
}

// routeHandler handles registering, unregistering, and querying routes
func (p *Server) routeHandler(w http.ResponseWriter, r *http.Request) {
	routesLogger := zerolog.Ctx(r.Context())
	if r.Method == http.MethodDelete {
		var route *Route
		if err := json.NewDecoder(r.Body).Decode(&route); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			routesLogger.Debug().Err(err).Msg("Failed to decode request body")
			return
		}
		defer r.Body.Close()

		err := p.Delete(route)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			routesLogger.Debug().Err(err).Msg("Failed to unregister route")
			return
		}

		w.WriteHeader(http.StatusNoContent)
		routesLogger.Info().
			Str("Route ID", route.ID()).
			Msg("Route deleted")
		return
	}

	if r.Method == http.MethodPost {
		var route *Route
		if err := json.NewDecoder(r.Body).Decode(&route); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			routesLogger.Debug().Err(err).Msg("Failed to decode request body")
			return
		}
		defer r.Body.Close()

		err := p.Register(route)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			routesLogger.Debug().Err(err).Msg("Failed to register route")
			return
		}

		w.WriteHeader(http.StatusCreated)
		return
	}

	if r.Method == http.MethodGet {
		routes := p.Routes()
		jsonRoutes, err := json.Marshal(routes)
		if err != nil {
			http.Error(w, "Failed to marshal routes", http.StatusInternalServerError)
			routesLogger.Debug().Err(err).Msg("Failed to marshal routes")
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if _, err = w.Write(jsonRoutes); err != nil {
			http.Error(w, "Failed to write response", http.StatusInternalServerError)
			routesLogger.Debug().Err(err).Msg("Failed to write response")
			return
		}

		routesLogger.Debug().Int("Count", len(routes)).Msg("Returned routes")
		return
	}

	http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
	routesLogger.Debug().Msg("Invalid method")
}

// dynamicHandler handles all incoming requests and responds based on the registered routes.
func (p *Server) dynamicHandler(w http.ResponseWriter, r *http.Request) {
	routeCallID := uuid.New().String()[0:8]
	dynamicLogger := zerolog.Ctx(r.Context())
	dynamicLogger.UpdateContext(func(c zerolog.Context) zerolog.Context {
		return c.Str("Route Call ID", routeCallID)
	})

	route, err := p.cage.getRoute(r.URL.Path, r.Method)
	if err != nil {
		if errors.Is(err, ErrRouteNotFound) {
			http.Error(w, "Route not found", http.StatusNotFound)
			dynamicLogger.Debug().Msg("Route not found")
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		dynamicLogger.Error().Err(err).Msg("Route called does not exist")
		return
	}

	dynamicLogger.UpdateContext(func(c zerolog.Context) zerolog.Context {
		return c.Str("Route ID", route.ID())
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
		RouteCallID: routeCallID,
		RouteID:     r.Method + ":" + r.URL.Path,
		Request: &RouteCallRequest{
			Method:     r.Method,
			URL:        r.URL,
			Header:     r.Header,
			Body:       requestBody,
			RemoteAddr: r.RemoteAddr,
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

	recordingWriter.WriteHeader(route.ResponseStatusCode)

	if route.RawResponseBody != "" {
		if _, err := recordingWriter.Write([]byte(route.RawResponseBody)); err != nil {
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
		if _, err = recordingWriter.Write(rawJSON); err != nil {
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

// recordHandler handles registering recorders with the parrot
func (p *Server) recordHandler(w http.ResponseWriter, r *http.Request) {
	recordingLogger := zerolog.Ctx(r.Context())
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		recordingLogger.Debug().Msg("Invalid method")
		return
	}

	var recorder *Recorder
	if err := json.NewDecoder(r.Body).Decode(&recorder); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		recordingLogger.Debug().Err(err).Msg("Failed to decode request body")
		return
	}
	defer r.Body.Close()

	if recorder == nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		recordingLogger.Debug().Msg("No recorder provided")
		return
	}

	if recorder.URL() == "" {
		http.Error(w, "Recorder URL required", http.StatusBadRequest)
		recordingLogger.Debug().Msg("No recorder URL provided")
		return
	}

	if err := p.Record(recorder.URL()); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		recordingLogger.Debug().Err(err).Msg("Failed to register recorder")
		return
	}

	w.WriteHeader(http.StatusCreated)
	recordingLogger.Debug().Str("URL", recorder.URL()).Msg("Recorder added")
}

func (p *Server) healthHandler(w http.ResponseWriter, _ *http.Request) {
	if p.shutDown {
		http.Error(w, "Server is shutting down", http.StatusServiceUnavailable)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// load loads all registered routes from a file.
func (p *Server) load() error {
	if _, err := os.Stat(p.saveFileName); os.IsNotExist(err) {
		p.log.Trace().Str("file", p.saveFileName).Msg("No data to load")
		return nil
	}

	p.log.Debug().Str("File", p.saveFileName).Msg("Loading data")

	fileData, err := os.ReadFile(p.saveFileName)
	if err != nil {
		return fmt.Errorf("failed to read routes from file: %w", err)
	}
	if len(fileData) == 0 {
		p.log.Trace().Str("File", p.saveFileName).Msg("No data to load")
		return nil
	}

	var saveData SaveFile
	err = json.Unmarshal(fileData, &saveData)
	if err != nil {
		return fmt.Errorf("failed to unmarshal save file: %w", err)
	}

	for _, route := range saveData.Routes {
		if err = p.Register(route); err != nil {
			return fmt.Errorf("failed to register route: %w", err)
		}
	}

	for _, recorder := range saveData.Recorders {
		if err = p.Record(recorder); err != nil {
			return fmt.Errorf("failed to register recorder: %w", err)
		}
	}

	p.log.Info().Str("file", p.saveFileName).Msg("Loaded routes")
	return nil
}

// save saves all registered routes to a file.
func (p *Server) save() error {
	saveFile := &SaveFile{
		Routes:    p.Routes(),
		Recorders: p.Recorders(),
	}
	if len(saveFile.Routes) == 0 && len(saveFile.Recorders) == 0 {
		p.log.Trace().Str("File", p.saveFileName).Msg("No data to save")
		return nil
	}

	jsonData, err := json.Marshal(saveFile)
	if err != nil {
		return fmt.Errorf("failed to marshal save file: %w", err)
	}

	if err = os.WriteFile(p.saveFileName, jsonData, 0644); err != nil { //nolint:gosec
		return fmt.Errorf("failed to write to save file: %w", err)
	}

	p.log.Debug().Str("File", p.saveFileName).Msg("Saved data")
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
	p.log.Trace().Int("Recorder Count", len(p.recorderHooks)).Str("Route ID", routeCall.RouteID).Msg("Sending route call to recorders")

	for hook := range p.recorderHooks {
		go func(hook string) {
			resp, err := client.R().SetBody(routeCall).Post(hook)
			if err != nil {
				p.log.Error().Err(err).Str("Recorder Hook", hook).Msg("Failed to send route call to recorder")
				return
			}
			defer resp.RawResponse.Body.Close()
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

var validPathRegex = regexp.MustCompile(`^\/[a-zA-Z0-9\-._~%!$&'()+,;=:@\/]`)

func isValidPath(path string) bool {
	switch path {
	case "", "/", "//", healthRoute, recordRoute, routesRoute, "/..":
		return false
	}
	if strings.Contains(path, "/..") {
		return false
	}
	if strings.Contains(path, "/.") {
		return false
	}
	if strings.Contains(path, "//") {
		return false
	}
	if !strings.HasPrefix(path, "/") {
		return false
	}
	if strings.HasSuffix(path, "/") {
		return false
	}
	if strings.HasPrefix(path, recordRoute) {
		return false
	}
	if strings.HasPrefix(path, healthRoute) {
		return false
	}
	if strings.HasPrefix(path, routesRoute) {
		return false
	}
	return validPathRegex.MatchString(path)
}

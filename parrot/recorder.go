package parrot

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"time"
)

type Recorder struct {
	URL    string `json:"url"`
	server *http.Server

	recordChan chan *RouteCall
	errChan    chan error
}

// RouteCall records that a route was called
type RouteCall struct {
	RouteID  string             `json:"route_id"`
	Request  *RouteCallRequest  `json:"request"`
	Response *RouteCallResponse `json:"response"`
}

// RouteCallRequest records the request made to a route
type RouteCallRequest struct {
	Method string      `json:"method"`
	URL    *url.URL    `json:"url"`
	Header http.Header `json:"header"`
	Body   []byte      `json:"body"`
}

// RouteCallResponse records the response from a route
type RouteCallResponse struct {
	StatusCode int         `json:"status_code"`
	Header     http.Header `json:"header"`
	Body       []byte      `json:"body"`
}

func NewRecorder() (*Recorder, error) {
	r := &Recorder{
		recordChan: make(chan *RouteCall),
		errChan:    make(chan error),
	}

	// TODO: Will need a way to send out the URL to an external service (e.g. Parrotserver running in a docker container)
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, fmt.Errorf("failed to start listener: %w", err)
	}
	r.URL = "http://" + listener.Addr().String()

	mux := http.NewServeMux()
	mux.Handle("/", r.defaultRecordHandler())
	r.server = &http.Server{
		ReadHeaderTimeout: 5 * time.Second,
		Addr:              listener.Addr().String(),
		Handler:           mux,
	}

	go func() {
		if err := r.server.Serve(listener); err != nil {
			if err != http.ErrServerClosed {
				r.errChan <- fmt.Errorf("error serving recorder: %w", err)
			}
		}
	}()
	return r, nil
}

// Record receives recorded calls
func (r *Recorder) Record() chan *RouteCall {
	return r.recordChan
}

func (r *Recorder) Close() error {
	close(r.recordChan)
	close(r.errChan)
	return r.server.Close()
}

func (r *Recorder) Err() chan error {
	return r.errChan
}

func (r *Recorder) defaultRecordHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var recordedCall *RouteCall
		if err := json.NewDecoder(req.Body).Decode(&recordedCall); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		defer req.Body.Close()

		r.recordChan <- recordedCall
	}
}

// httpResponseRecorder is a wrapper around http.ResponseWriter that records the response
// for later inspection while still writing to the original writer.
// WARNING: If you mutate after calling Header(), the changes will not be reflected in the recorded response.
type responseWriterRecorder struct {
	originalWriter http.ResponseWriter
	record         *httptest.ResponseRecorder
}

func newResponseWriterRecorder(w http.ResponseWriter) *responseWriterRecorder {
	return &responseWriterRecorder{
		originalWriter: w,
		record:         httptest.NewRecorder(),
	}
}

// SetWriter sets a new writer to record and write to, flushing any previous record
func (rr *responseWriterRecorder) SetWriter(w http.ResponseWriter) {
	rr.originalWriter = w
	rr.record = httptest.NewRecorder()
}

func (rr *responseWriterRecorder) WriteHeader(code int) {
	rr.originalWriter.WriteHeader(code)
	rr.record.WriteHeader(code)
}

func (rr *responseWriterRecorder) Write(data []byte) (int, error) {
	_, _ = rr.record.Write(data) // ignore error as we still want to write to the original writer
	return rr.originalWriter.Write(data)
}

func (rr *responseWriterRecorder) Header() http.Header {
	for k, v := range rr.originalWriter.Header() {
		rr.record.Header()[k] = v
	}
	return rr.originalWriter.Header()
}

func (rr *responseWriterRecorder) Result() *http.Response {
	return rr.record.Result()
}

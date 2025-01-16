package parrot

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
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
	RouteID  string         `json:"route_id"`
	Request  *http.Request  `json:"request"`
	Response *http.Response `json:"response"`
}

type RecorderOption func(*Recorder) error

func NewRecorder(opts ...RecorderOption) (*Recorder, error) {
	r := &Recorder{
		recordChan: make(chan *RouteCall),
		errChan:    make(chan error),
	}

	for _, opt := range opts {
		err := opt(r)
		if err != nil {
			return nil, err
		}
	}

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, fmt.Errorf("failed to start listener: %w", err)
	}
	r.URL = listener.Addr().String()

	mux := http.NewServeMux()
	mux.Handle("/", r.defaultRecordHandler())
	r.server = &http.Server{
		ReadHeaderTimeout: 5 * time.Second,
		Addr:              listener.Addr().String(),
		Handler:           mux,
	}

	go func() {
		if err := r.server.Serve(listener); err != nil {
			r.errChan <- err
		}
	}()
	return r, nil
}

// Record receives recorded calls
func (r *Recorder) Record() chan *RouteCall {
	return r.recordChan
}

func (r *Recorder) Close() error {
	return r.server.Close()
}

func (r *Recorder) Error() chan error {
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

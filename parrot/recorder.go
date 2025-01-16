package parrot

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
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
	RouteID string       `json:"route_id"`
	Request http.Request `json:"request"`
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
	mux.Handle("/record", r.defaultRecordHandler())
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

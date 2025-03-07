package fake

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

var (
	R  = NewRecords()
	mu sync.Mutex
)

// Record is a request and response data
type Record struct {
	Method  string      `json:"method"`
	Path    string      `json:"path"`
	Headers http.Header `json:"headers"`
	ReqBody string      `json:"req_body"`
	ResBody string      `json:"res_body"`
	Status  int         `json:"status"`
}

type Records struct {
	Data map[string][]*Record
}

func NewRecords() *Records {
	return &Records{make(map[string][]*Record)}
}

func RecordKey(method, path string) string {
	return fmt.Sprintf("%s:%s", method, path)
}

func (r *Records) Get(method, path string) ([]*Record, error) {
	rec, ok := r.Data[RecordKey(method, path)]
	if !ok {
		return nil, fmt.Errorf("no record was found for path: %s", path)
	}
	return rec, nil
}

// CustomResponseWriter wraps gin.ResponseWriter to capture response data
type CustomResponseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// Write captures response data
func (w *CustomResponseWriter) Write(data []byte) (int, error) {
	w.body.Write(data) // Capture response data
	return w.ResponseWriter.Write(data)
}

// Middleware to record both requests and responses
func recordMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Capture request data
		var reqBodyBytes []byte
		if c.Request.Body != nil {
			reqBodyBytes, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(reqBodyBytes))
		}
		reqBody := string(reqBodyBytes)

		// Create custom response writer
		customWriter := &CustomResponseWriter{
			ResponseWriter: c.Writer,
			body:           bytes.NewBufferString(""),
		}
		c.Writer = customWriter

		// Process request
		c.Next()

		// Capture response data
		resBody := customWriter.body.String()
		status := c.Writer.Status()
		mu.Lock()
		defer mu.Unlock()
		if R.Data[c.Request.URL.Path] == nil {
			R.Data[c.Request.URL.Path] = make([]*Record, 0)
		}
		recKey := RecordKey(c.Request.Method, c.Request.URL.Path)
		R.Data[recKey] = append(R.Data[recKey], &Record{
			Method:  c.Request.Method,
			Path:    c.Request.URL.Path,
			Headers: c.Request.Header,
			ReqBody: reqBody,
			ResBody: resBody,
			Status:  status,
		},
		)
	}
}

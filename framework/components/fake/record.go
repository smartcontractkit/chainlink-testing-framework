package fake

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)

var (
	R = NewRecords()
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
	r map[string][]*Record
}

func NewRecords() *Records {
	return &Records{make(map[string][]*Record)}
}

func (r *Records) Get(path string) ([]*Record, error) {
	rec, ok := r.r[path]
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
			reqBodyBytes, _ = ioutil.ReadAll(c.Request.Body)
			c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(reqBodyBytes))
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
		if R.r[c.Request.URL.Path] == nil {
			R.r[c.Request.URL.Path] = make([]*Record, 0)
		}
		R.r[c.Request.URL.Path] = append(R.r[c.Request.URL.Path], &Record{
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

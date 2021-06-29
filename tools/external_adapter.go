package tools

import (
	"encoding/json"
	"fmt"
	"github.com/phayes/freeport"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog/log"
)

type ExternalAdapter struct {
	LocalAddr        string
	InsideDockerAddr string
}

type ExternalAdapterResponse struct {
	JobRunId string              `json:"id"`
	Data     ExternalAdapterData `json:"data"`
	Error    error               `json:"error"`
}

type ExternalAdapterData struct {
	Result int `json:"result"`
}

type OkResult struct{}

var variableData int

func FreePort() string {
	port, err := freeport.GetFreePort()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to acquire free port")
	}
	return strconv.Itoa(port)
}

// NewExternalAdapter starts an external adapter on specified port
func NewExternalAdapter() ExternalAdapter {
	p := FreePort()
	// TODO: graceful is needed for some cases and nightly
	go func() {
		router := httprouter.New()
		router.GET("/", index)
		router.POST("/random", randomNumber)
		router.POST("/five", five)
		router.POST("/variable", variable)
		router.POST("/set_variable", setVariable)

		log.Info().Str("Port", p).Msg("Starting external adapter")
		log.Fatal().AnErr("Error", http.ListenAndServe(":"+p, router)).Msg("Error occured while running external adapter")
	}()
	time.Sleep(1 * time.Second)
	return ExternalAdapter{
		LocalAddr:        fmt.Sprintf("http://0.0.0.0:%s", p),
		InsideDockerAddr: fmt.Sprintf("http://host.docker.internal:%s", p),
	}
}

func SetVariableMockData(url string, data int) (*http.Response, error) {
	resp, err := http.Post(
		fmt.Sprintf("%s/set_variable?var=%d", url, data),
		"application/json",
		nil,
	)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func GetVariableMockData(url string) (*http.Response, error) {
	resp, err := http.Post(
		fmt.Sprintf("%s/variable", url),
		"application/json",
		nil,
	)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// index allows a status check on the adapter
func index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Adapter listening!")
}

// RandomNumber returns a random int from 0 to 100
func randomNumber(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	num := rand.Intn(100)
	result := &ExternalAdapterResponse{
		JobRunId: "",
		Data:     ExternalAdapterData{Result: num},
		Error:    nil,
	}
	_ = json.NewEncoder(w).Encode(result)
}

// five returns five
func five(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	result := &ExternalAdapterResponse{
		JobRunId: "",
		Data:     ExternalAdapterData{Result: 5},
		Error:    nil,
	}
	_ = json.NewEncoder(w).Encode(result)
}

func setVariable(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	q := r.URL.Query()
	log.Info().Interface("query", q).Msg("params received")
	v := q.Get("var")
	data, _ := strconv.Atoi(v)
	variableData = data
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	result := &OkResult{}
	_ = json.NewEncoder(w).Encode(result)
}

func variable(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	log.Info().Int("data", variableData).Msg("variable response")
	result := &ExternalAdapterResponse{
		JobRunId: "",
		Data:     ExternalAdapterData{Result: variableData},
		Error:    nil,
	}
	_ = json.NewEncoder(w).Encode(result)
}

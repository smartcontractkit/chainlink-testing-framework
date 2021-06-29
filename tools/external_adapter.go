package tools

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog/log"
)

type ExternalAdapterResponse struct {
	JobRunId string              `json:"id"`
	Data     ExternalAdapterData `json:"data"`
	Error    error               `json:"error"`
}

type ExternalAdapterFloatResponse struct {
	JobRunId string                   `json:"id"`
	Data     ExternalAdapterFloatData `json:"data"`
	Error    error                    `json:"error"`
}

type ExternalAdapterData struct {
	Result int `json:"result"`
}

type ExternalAdapterFloatData struct {
	Result float64 `json:"result"`
}

type OkResult struct{}

var variableData float64

// NewExternalAdapter starts an external adapter on specified port
func NewExternalAdapter(portNumber string) {
	router := httprouter.New()
	router.GET("/", index)
	router.POST("/random", randomNumber)
	router.POST("/five", five)
	router.POST("/variable", variable)
	router.POST("/set_variable", setVariable)

	log.Info().Str("Port", portNumber).Msg("Starting external adapter")
	log.Fatal().AnErr("Error", http.ListenAndServe(":"+portNumber, router)).Msg("Error occured while running external adapter")
}

func SetVariableMockData(data float64) (*http.Response, error) {
	resp, err := http.Post(
		fmt.Sprintf("%s/set_variable?var=%f", "http://0.0.0.0:6644", data),
		"application/json",
		nil,
	)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func GetVariableMockData() (*http.Response, error) {
	resp, err := http.Post(
		fmt.Sprintf("%s/variable", "http://0.0.0.0:6644"),
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
	data, _ := strconv.ParseFloat(v, 64)
	variableData = data
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	result := &OkResult{}
	_ = json.NewEncoder(w).Encode(result)
}

func variable(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	log.Info().Float64("data", variableData).Msg("variable response")
	result := &ExternalAdapterFloatResponse{
		JobRunId: "",
		Data:     ExternalAdapterFloatData{Result: variableData},
		Error:    nil,
	}
	_ = json.NewEncoder(w).Encode(result)
}

package tools

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog/log"
)

type ExternalAdapterResponse struct {
	JobRunId string              `json:"id"`
	Data     ExternalAdapterData `json:"data"`
	Error    error               `json:"error"`
}

type ExternalAdapterData struct {
	Result int `json:"result"`
}

// NewExternalAdapter starts an external adapter on specified port
func NewExternalAdapter(portNumber string) {
	router := httprouter.New()
	router.GET("/", index)
	router.POST("/random", randomNumber)
	router.POST("/five", five)

	log.Info().Str("Port", portNumber).Msg("Starting external adapter")
	log.Fatal().AnErr("Error", http.ListenAndServe(":"+portNumber, router)).Msg("Error occured while running external adapter")
}

func index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Adapter listening!")
}

// RandomNumber returns a random int from 0 to 10
func randomNumber(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	num := rand.Intn(100)
	result := &ExternalAdapterResponse{
		JobRunId: "",
		Data:     ExternalAdapterData{Result: num},
		Error:    nil,
	}
	log.Info().Int("Result", num).Str("Requester", r.RemoteAddr).Msg("Sending Result from Adapter")
	json.NewEncoder(w).Encode(result)
}

func five(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	result := &ExternalAdapterResponse{
		JobRunId: "",
		Data:     ExternalAdapterData{Result: 5},
		Error:    nil,
	}
	log.Info().Int("Result", 5).Str("Requester", r.RemoteAddr).Msg("Sending Result from Adapter")
	json.NewEncoder(w).Encode(result)

}

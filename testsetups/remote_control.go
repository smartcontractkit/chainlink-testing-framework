package testsetups

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog/log"
)

var (
	started         bool
	testName        string
	confirmCode     int
	stopChannel     chan struct{}
	shutDownStarted bool
)

// StartRemoteControlServer starts a server to control soak test shut downs. Shutdown signal is sent to the stopChan
func StartRemoteControlServer(name string, stopChan chan struct{}) {
	accessPort := os.Getenv("ACCESS_PORT")
	if accessPort == "" {
		accessPort = "8080"
	}
	if started {
		log.Warn().Str("Port", accessPort).Msg("Already started remote control server")
		return
	}
	started = true
	testName = name
	stopChannel = stopChan

	router := httprouter.New()
	router.GET("/", Index)
	router.GET("/shutdown", ShutDown)
	router.GET("/shutdown/:confirmCode", ShutDownConfirm)

	log.Info().Str("Port", accessPort).Msg("Remote control server listening")

	// Start listening for shutdown calls
	go func() {
		log.Error().Err(http.ListenAndServe(fmt.Sprintf(":%s", accessPort), router)).Msg("Error running remote control server")
	}()

	// TODO: Enable this shutdown to happen on SIGTERMs as well (I messed around with signal.Notify and couldn't get it to work for some reason)
}

// Index to see barbones running statistics
func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if shutDownStarted {
		fmt.Fprint(w, "Shut down is in progress")
		return
	}
	fmt.Fprint(w, "Remote Control Server is Running!\nUse /shutdown to initiate a graceful test shutdown")
}

// ShutDown attempts to gracefully shut down the running test
func ShutDown(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if shutDownStarted {
		fmt.Fprint(w, "Shut down already in progress")
		return
	}
	confirmCode = 1000 + rand.Intn(9999-1000) // #nosec G404 | 4 digit random number
	fmt.Fprintf(w, "Are you sure you want to shut down '%s'? Call /shutdown/%d to confirm.", testName, confirmCode)
}

// ShutDownConfirm requires the user to enter a randomly generated code to shut down the test
func ShutDownConfirm(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if shutDownStarted {
		fmt.Fprint(w, "Shut down already in progress")
		return
	}
	enteredCode := ps.ByName("confirmCode")
	if enteredCode == fmt.Sprint(confirmCode) {
		fmt.Fprintf(w, "Shutting down '%s'", testName)
		shutdown()
	} else {
		fmt.Fprintf(w, "Incorrect confirmation code '%s'", enteredCode)
	}
}

func shutdown() {
	shutDownStarted = true
	stopChannel <- struct{}{}
	close(stopChannel)
}

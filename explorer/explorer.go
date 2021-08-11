package explorer

import (
	"encoding/json"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"net/http"
	"sync"
)

const (
	WSReadBufferSizeDefault  = 1024
	WSWriteBufferSizeDefault = 1024
)

var (
	anyOriginWSUpgrader = websocket.Upgrader{
		ReadBufferSize:  WSReadBufferSizeDefault,
		WriteBufferSize: WSWriteBufferSizeDefault,
		CheckOrigin:     func(*http.Request) bool { return true },
	}
)

// MessagesCount counts total exolorer messages received
type MessagesCount struct {
	Received        int `json:"received"`
	Broadcast       int `json:"broadcast"`
	Sent            int `json:"sent"`
	AssertViolation int `json:"assert_violation"`
	RoundStarted    int `json:"round_started"`
	Unknown         int `json:"unknown"`
	Errors          int `json:"errors"`
}

// Messages explorer messages by type
type Messages struct {
	Received        []string `json:"received"`
	Broadcast       []string `json:"broadcast"`
	Sent            []string `json:"sent"`
	AssertViolation []string `json:"assert_violation"`
	RoundStarted    []string `json:"round_started"`
	Unknown         []string `json:"unknown"`
	Errors          []string `json:"errors"`
}

// Explorer explorer data
type Explorer struct {
	MessagesMu *sync.RWMutex
	Messages   *Messages `json:"messages"`
}

func NewExplorer() *Explorer {
	return &Explorer{
		MessagesMu: &sync.RWMutex{},
		Messages: &Messages{
			Received:        make([]string, 0),
			Broadcast:       make([]string, 0),
			Sent:            make([]string, 0),
			AssertViolation: make([]string, 0),
			RoundStarted:    make([]string, 0),
			Unknown:         make([]string, 0),
			Errors:          make([]string, 0),
		},
	}
}

func (e *Explorer) handleGetMessages(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	e.MessagesMu.RLock()
	defer e.MessagesMu.RUnlock()
	if err := json.NewEncoder(w).Encode(e.Messages); err != nil {
		log.Error().Err(err).Send()
	}
}

func (e *Explorer) handleCountMessages(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	e.MessagesMu.RLock()
	defer e.MessagesMu.RUnlock()
	mc := &MessagesCount{
		Received:        len(e.Messages.Received),
		Broadcast:       len(e.Messages.Broadcast),
		Sent:            len(e.Messages.Sent),
		AssertViolation: len(e.Messages.AssertViolation),
		RoundStarted:    len(e.Messages.RoundStarted),
		Errors:          len(e.Messages.Errors),
	}
	if err := json.NewEncoder(w).Encode(mc); err != nil {
		log.Error().Err(err).Send()
	}
}

// handleWSMessages receives and saves websocket messages to explorer
func (e *Explorer) handleWSMessages(w http.ResponseWriter, r *http.Request) {
	conn, err := anyOriginWSUpgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Error().Err(err).Send()
		e.MessagesMu.Lock()
		e.Messages.Errors = append(e.Messages.Errors, err.Error())
		e.MessagesMu.Unlock()
		return
	}
	for {
		_, bs, err := conn.ReadMessage()
		e.MessagesMu.Lock()
		if err != nil {
			log.Error().Err(err).Send()
			e.Messages.Errors = append(e.Messages.Errors, err.Error())
			e.MessagesMu.Unlock()
			return
		}

		var msg TelemetryWrapper
		err = proto.Unmarshal(bs, &msg)
		if err != nil {
			log.Error().Err(errors.Wrapf(err, "Telemetry pb parsing error: %s", string(bs))).Send()
			e.Messages.Errors = append(e.Messages.Errors, err.Error())
			e.MessagesMu.Unlock()
			return
		}
		switch {
		case msg.GetMessageReceived() != nil:
			log.Debug().Str("Received", msg.String()).Send()
			e.Messages.Received = append(e.Messages.Received, msg.String())
		case msg.GetMessageBroadcast() != nil:
			log.Debug().Str("Broadcast", msg.String()).Send()
			e.Messages.Broadcast = append(e.Messages.Received, msg.String())
		case msg.GetMessageSent() != nil:
			log.Debug().Str("Sent", msg.String()).Send()
			e.Messages.Sent = append(e.Messages.Received, msg.String())
		case msg.GetAssertionViolation() != nil:
			log.Debug().Str("AssertViolation", msg.String()).Send()
			e.Messages.AssertViolation = append(e.Messages.Received, msg.String())
		case msg.GetRoundStarted() != nil:
			log.Debug().Str("Round started", msg.String()).Send()
			e.Messages.RoundStarted = append(e.Messages.Received, msg.String())
		default:
			log.Error().Err(errors.New("unknown message type")).Send()
			e.Messages.Unknown = append(e.Messages.Received, msg.String())
		}
		e.MessagesMu.Unlock()
	}
}

func (e *Explorer) Run() error {
	port := ":4321"
	log.Info().Str("Port", port).Msg("Starting explorer stub")
	http.Handle("/", http.HandlerFunc(e.handleWSMessages))
	http.Handle("/messages", http.HandlerFunc(e.handleGetMessages))
	http.Handle("/count", http.HandlerFunc(e.handleCountMessages))
	if err := http.ListenAndServe(port, nil); err != nil {
		return err
	}
	return nil
}

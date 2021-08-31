package explorer

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/proto"
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
	// telemetry
	Received        int `json:"received"`
	Broadcast       int `json:"broadcast"`
	Sent            int `json:"sent"`
	AssertViolation int `json:"assert_violation"`
	RoundStarted    int `json:"round_started"`

	// messages
	NewEpoch   int `json:"new_epoch"`
	Final      int `json:"final"`
	FinalEcho  int `json:"final_echo"`
	Observe    int `json:"observe"`
	ObserveReq int `json:"observe_req"`
	Report     int `json:"report"`
	ReportReq  int `json:"report_req"`

	// DHT
	DHTAnnounce int `json:"dht_announce"`

	Unknown int `json:"unknown"`
	Errors  int `json:"errors"`
}

// Messages explorer messages by type
type Messages struct {
	// telemetry
	Received        []*TelemetryMessageReceived    `json:"received"`
	Broadcast       []*TelemetryMessageBroadcast   `json:"broadcast"`
	Sent            []*TelemetryMessageSent        `json:"sent"`
	AssertViolation []*TelemetryAssertionViolation `json:"assert_violation"`
	RoundStarted    []*TelemetryRoundStarted       `json:"round_started"`

	// messages
	NewEpoch   []*MessageNewEpoch   `json:"new_epoch"`
	Final      []*MessageFinal      `json:"final"`
	FinalEcho  []*MessageFinalEcho  `json:"final_echo"`
	Observe    []*MessageObserve    `json:"observe"`
	ObserveReq []*MessageObserveReq `json:"observe_req"`
	Report     []*MessageReport     `json:"report"`
	ReportReq  []*MessageReportReq  `json:"report_req"`

	// DHT
	DHTAnnounce []*SignedAnnouncement `json:"dht_announce"`

	Unknown []string `json:"unknown"`
	Errors  []string `json:"errors"`
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
			Received:        make([]*TelemetryMessageReceived, 0),
			Broadcast:       make([]*TelemetryMessageBroadcast, 0),
			Sent:            make([]*TelemetryMessageSent, 0),
			AssertViolation: make([]*TelemetryAssertionViolation, 0),
			RoundStarted:    make([]*TelemetryRoundStarted, 0),
			NewEpoch:        make([]*MessageNewEpoch, 0),
			Final:           make([]*MessageFinal, 0),
			FinalEcho:       make([]*MessageFinalEcho, 0),
			Observe:         make([]*MessageObserve, 0),
			ObserveReq:      make([]*MessageObserveReq, 0),
			Report:          make([]*MessageReport, 0),
			ReportReq:       make([]*MessageReportReq, 0),
			DHTAnnounce:     make([]*SignedAnnouncement, 0),
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
		NewEpoch:        len(e.Messages.NewEpoch),
		Final:           len(e.Messages.Final),
		FinalEcho:       len(e.Messages.FinalEcho),
		Observe:         len(e.Messages.Observe),
		ObserveReq:      len(e.Messages.ObserveReq),
		Report:          len(e.Messages.Report),
		ReportReq:       len(e.Messages.ReportReq),
		DHTAnnounce:     len(e.Messages.DHTAnnounce),
		Errors:          len(e.Messages.Errors),
	}
	if err := json.NewEncoder(w).Encode(mc); err != nil {
		log.Error().Err(err).Send()
	}
}

func (e *Explorer) tryUnmarshalDHT(bs []byte, foundFlag *bool) {
	var msg SignedAnnouncement
	_ = proto.Unmarshal(bs, &msg)
	switch {
	case msg.GetAddrs() != nil:
		*foundFlag = true
		log.Debug().Interface("DHT Announce", msg.ProtoReflect()).Send()
		e.Messages.DHTAnnounce = append(e.Messages.DHTAnnounce, &msg)
	}
}

func (e *Explorer) tryUnmarshalMessages(bs []byte, foundFlag *bool) {
	var msg MessageWrapper
	_ = proto.Unmarshal(bs, &msg)
	switch {
	case msg.GetMessageNewEpoch() != nil:
		*foundFlag = true
		log.Debug().Interface("New epoch", msg.GetMessageNewEpoch()).Send()
		e.Messages.NewEpoch = append(e.Messages.NewEpoch, msg.GetMessageNewEpoch())
	case msg.GetMessageFinal() != nil:
		*foundFlag = true
		log.Debug().Interface("Final", msg.GetMessageFinal()).Send()
		e.Messages.Final = append(e.Messages.Final, msg.GetMessageFinal())
	case msg.GetMessageFinalEcho() != nil:
		*foundFlag = true
		log.Debug().Interface("FinalEcho", msg.GetMessageFinalEcho()).Send()
		e.Messages.FinalEcho = append(e.Messages.FinalEcho, msg.GetMessageFinalEcho())
	case msg.GetMessageObserve() != nil:
		*foundFlag = true
		log.Debug().Interface("Observe", msg.GetMessageObserve()).Send()
		e.Messages.Observe = append(e.Messages.Observe, msg.GetMessageObserve())
	case msg.GetMessageObserveReq() != nil:
		*foundFlag = true
		log.Debug().Interface("ObserveReq", msg.GetMessageObserveReq()).Send()
		e.Messages.ObserveReq = append(e.Messages.ObserveReq, msg.GetMessageObserveReq())
	case msg.GetMessageReport() != nil:
		*foundFlag = true
		log.Debug().Interface("Report", msg.GetMessageReport()).Send()
		e.Messages.Report = append(e.Messages.Report, msg.GetMessageReport())
	case msg.GetMessageReportReq() != nil:
		*foundFlag = true
		log.Debug().Interface("ReportReq", msg.GetMessageReportReq()).Send()
		e.Messages.ReportReq = append(e.Messages.ReportReq, msg.GetMessageReportReq())
	}
}

func (e *Explorer) tryUnmarshalTelemetry(bs []byte, foundFlag *bool) {
	var msg TelemetryWrapper
	_ = proto.Unmarshal(bs, &msg)
	switch {
	case msg.GetMessageReceived() != nil:
		*foundFlag = true
		log.Debug().Interface("Received", msg.GetMessageReceived()).Send()
		e.Messages.Received = append(e.Messages.Received, msg.GetMessageReceived())
	case msg.GetMessageBroadcast() != nil:
		*foundFlag = true
		log.Debug().Interface("Broadcast", msg.GetMessageBroadcast()).Send()
		e.Messages.Broadcast = append(e.Messages.Broadcast, msg.GetMessageBroadcast())
	case msg.GetMessageSent() != nil:
		*foundFlag = true
		log.Debug().Interface("Sent", msg.GetMessageSent()).Send()
		e.Messages.Sent = append(e.Messages.Sent, msg.GetMessageSent())
	case msg.GetAssertionViolation() != nil:
		*foundFlag = true
		log.Debug().Interface("AssertViolation", msg.GetAssertionViolation()).Send()
		e.Messages.AssertViolation = append(e.Messages.AssertViolation, msg.GetAssertionViolation())
	case msg.GetRoundStarted() != nil:
		*foundFlag = true
		log.Debug().Interface("Round started", msg.GetRoundStarted()).Send()
		e.Messages.RoundStarted = append(e.Messages.RoundStarted, msg.GetRoundStarted())
	}
}

func (e *Explorer) unmarshalMessage(bs []byte) {
	foundFlag := false
	e.tryUnmarshalTelemetry(bs, &foundFlag)
	e.tryUnmarshalMessages(bs, &foundFlag)
	e.tryUnmarshalDHT(bs, &foundFlag)
	if !foundFlag {
		e.Messages.Unknown = append(e.Messages.Unknown, string(bs))
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
		e.unmarshalMessage(bs)
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

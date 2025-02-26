package havoc

import (
	"time"

	"github.com/rs/zerolog"
)

type ChaosLogger struct {
	logger zerolog.Logger
}

func NewChaosLogger(logger zerolog.Logger) *ChaosLogger {
	return &ChaosLogger{logger: logger}
}

func (l ChaosLogger) OnChaosCreated(chaos Chaos) {
	l.commonChaosLog("info", chaos).Msg("Chaos created")
}

func (l ChaosLogger) OnChaosCreationFailed(chaos Chaos, reason error) {
	l.commonChaosLog("error", chaos).
		Err(reason).
		Msg("Failed to create chaos object")
}

func (l ChaosLogger) OnChaosStarted(chaos Chaos) {
	experiment, _ := chaos.GetExperimentStatus()

	l.commonChaosLog("info", chaos).
		Interface("spec", chaos.GetChaosSpec()).
		Interface("records", experiment.Records).
		Msg("Chaos started")
}

func (l ChaosLogger) OnChaosPaused(chaos Chaos) {
	l.commonChaosLog("info", chaos).
		Msg("Chaos paused")
}

func (l ChaosLogger) OnChaosEnded(chaos Chaos) {
	l.commonChaosLog("info", chaos).
		Msg("Chaos ended")
}

func (l ChaosLogger) OnChaosDeleted(chaos Chaos) {
	l.commonChaosLog("info", chaos).
		Msg("Chaos deleted")
}

type SimplifiedEvent struct {
	LastTimestamp string
	Type          string
	Message       string
}

func (l ChaosLogger) OnChaosStatusUnknown(chaos Chaos) {
	status, _ := chaos.GetExperimentStatus()
	events, _ := chaos.GetChaosEvents()

	// Create a slice to hold the simplified events
	simplifiedEvents := make([]SimplifiedEvent, 0, len(events.Items))

	// Iterate over the events and extract the required information
	for _, event := range events.Items {
		simplifiedEvents = append(simplifiedEvents, SimplifiedEvent{
			LastTimestamp: event.LastTimestamp.Time.Format(time.RFC3339),
			Type:          event.Type,
			Message:       event.Message,
		})
	}

	l.commonChaosLog("error", chaos).
		Interface("status", status).
		Interface("events", simplifiedEvents).
		Msg("Chaos status unknown")
}

func (l ChaosLogger) commonChaosLog(logLevel string, chaos Chaos) *zerolog.Event {
	// Create a base event based on the dynamic log level
	var event *zerolog.Event
	switch logLevel {
	case "debug":
		event = l.logger.Debug()
	case "info":
		event = l.logger.Info()
	case "warn":
		event = l.logger.Warn()
	case "error":
		event = l.logger.Error()
	case "fatal":
		event = l.logger.Fatal()
	case "panic":
		event = l.logger.Panic()
	default:
		// Default to info level if an unknown level is provided
		event = l.logger.Info()
	}

	duration, _ := chaos.GetChaosDuration()

	return event.
		Str("logger", "chaos").
		Str("name", chaos.GetObject().GetName()).
		Str("namespace", chaos.GetObject().GetNamespace()).
		Str("description", chaos.GetChaosDescription()).
		Str("duration", duration.String()).
		Time("startTime", chaos.GetStartTime()).
		Time("endTime", chaos.GetEndTime())
}

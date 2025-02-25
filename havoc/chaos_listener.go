package havoc

// ChaosListener is an interface that can be implemented by clients to listen to and react to chaos events.
type ChaosListener interface {
	OnChaosCreated(chaos Chaos)
	OnChaosCreationFailed(chaos Chaos, reason error)
	OnChaosStarted(chaos Chaos)
	OnChaosPaused(chaos Chaos)
	OnChaosEnded(chaos Chaos)         // When the chaos is finished or deleted
	OnChaosStatusUnknown(chaos Chaos) // When the chaos status is unknown
}

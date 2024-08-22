package k8schaos

type ChaosListener interface {
	OnChaosCreated(chaos Chaos)
	OnChaosCreationFailed(chaos Chaos, reason error)
	OnChaosStarted(chaos Chaos)
	OnChaosPaused(chaos Chaos)
	OnChaosEnded(chaos Chaos)         // When the chaos is finished or deleted
	OnChaosStatusUnknown(chaos Chaos) // When the chaos status is unknown
	OnScheduleCreated(chaos Schedule)
	OnScheduleDeleted(chaos Schedule) // When the chaos is finished or deleted
}

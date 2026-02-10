package k8s


// SleepAction describes a "sleep" action.
type SleepAction struct {
	// Seconds is the number of seconds to sleep.
	Seconds *float64 `field:"required" json:"seconds" yaml:"seconds"`
}


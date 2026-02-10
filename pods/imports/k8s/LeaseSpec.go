package k8s

import (
	"time"
)

// LeaseSpec is a specification of a Lease.
type LeaseSpec struct {
	// acquireTime is a time when the current lease was acquired.
	AcquireTime *time.Time `field:"optional" json:"acquireTime" yaml:"acquireTime"`
	// holderIdentity contains the identity of the holder of a current lease.
	//
	// If Coordinated Leader Election is used, the holder identity must be equal to the elected LeaseCandidate.metadata.name field.
	HolderIdentity *string `field:"optional" json:"holderIdentity" yaml:"holderIdentity"`
	// leaseDurationSeconds is a duration that candidates for a lease need to wait to force acquire it.
	//
	// This is measured against the time of last observed renewTime.
	LeaseDurationSeconds *float64 `field:"optional" json:"leaseDurationSeconds" yaml:"leaseDurationSeconds"`
	// leaseTransitions is the number of transitions of a lease between holders.
	LeaseTransitions *float64 `field:"optional" json:"leaseTransitions" yaml:"leaseTransitions"`
	// PreferredHolder signals to a lease holder that the lease has a more optimal holder and should be given up.
	//
	// This field can only be set if Strategy is also set.
	PreferredHolder *string `field:"optional" json:"preferredHolder" yaml:"preferredHolder"`
	// renewTime is a time when the current holder of a lease has last updated the lease.
	RenewTime *time.Time `field:"optional" json:"renewTime" yaml:"renewTime"`
	// Strategy indicates the strategy for picking the leader for coordinated leader election.
	//
	// If the field is not specified, there is no active coordination for this lease. (Alpha) Using this field requires the CoordinatedLeaderElection feature gate to be enabled.
	Strategy *string `field:"optional" json:"strategy" yaml:"strategy"`
}


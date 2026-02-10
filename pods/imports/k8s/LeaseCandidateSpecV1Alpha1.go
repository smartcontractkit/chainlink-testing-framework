package k8s

import (
	"time"
)

// LeaseCandidateSpec is a specification of a Lease.
type LeaseCandidateSpecV1Alpha1 struct {
	// LeaseName is the name of the lease for which this candidate is contending.
	//
	// This field is immutable.
	LeaseName *string `field:"required" json:"leaseName" yaml:"leaseName"`
	// PreferredStrategies indicates the list of strategies for picking the leader for coordinated leader election.
	//
	// The list is ordered, and the first strategy supersedes all other strategies. The list is used by coordinated leader election to make a decision about the final election strategy. This follows as - If all clients have strategy X as the first element in this list, strategy X will be used. - If a candidate has strategy [X] and another candidate has strategy [Y, X], Y supersedes X and strategy Y
	// will be used.
	// - If a candidate has strategy [X, Y] and another candidate has strategy [Y, X], this is a user error and leader
	// election will not operate the Lease until resolved.
	// (Alpha) Using this field requires the CoordinatedLeaderElection feature gate to be enabled.
	PreferredStrategies *[]*string `field:"required" json:"preferredStrategies" yaml:"preferredStrategies"`
	// BinaryVersion is the binary version.
	//
	// It must be in a semver format without leading `v`. This field is required when strategy is "OldestEmulationVersion"
	BinaryVersion *string `field:"optional" json:"binaryVersion" yaml:"binaryVersion"`
	// EmulationVersion is the emulation version.
	//
	// It must be in a semver format without leading `v`. EmulationVersion must be less than or equal to BinaryVersion. This field is required when strategy is "OldestEmulationVersion"
	EmulationVersion *string `field:"optional" json:"emulationVersion" yaml:"emulationVersion"`
	// PingTime is the last time that the server has requested the LeaseCandidate to renew.
	//
	// It is only done during leader election to check if any LeaseCandidates have become ineligible. When PingTime is updated, the LeaseCandidate will respond by updating RenewTime.
	PingTime *time.Time `field:"optional" json:"pingTime" yaml:"pingTime"`
	// RenewTime is the time that the LeaseCandidate was last updated.
	//
	// Any time a Lease needs to do leader election, the PingTime field is updated to signal to the LeaseCandidate that they should update the RenewTime. Old LeaseCandidate objects are also garbage collected if it has been hours since the last renew. The PingTime field is updated regularly to prevent garbage collection for still active LeaseCandidates.
	RenewTime *time.Time `field:"optional" json:"renewTime" yaml:"renewTime"`
}


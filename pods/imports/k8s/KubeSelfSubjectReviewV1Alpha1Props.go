package k8s


// SelfSubjectReview contains the user information that the kube-apiserver has about the user making this request.
//
// When using impersonation, users will receive the user info of the user being impersonated.  If impersonation or request header authentication is used, any extra keys will have their case ignored and returned as lowercase.
type KubeSelfSubjectReviewV1Alpha1Props struct {
	// Standard object's metadata.
	//
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	Metadata *ObjectMeta `field:"optional" json:"metadata" yaml:"metadata"`
}


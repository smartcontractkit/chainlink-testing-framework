package client

import "net/http"

// KafkaRestClient kafka-rest client
type KafkaRestClient struct {
	*APIClient
	Config *KafkaRestConfig
}

// KafkaRestConfig holds config information for KafkaRestClient
type KafkaRestConfig struct {
	URL string
}

// NewKafkaRestClient creates a new KafkaRestClient
func NewKafkaRestClient(cfg *KafkaRestConfig) *KafkaRestClient {
	return &KafkaRestClient{
		Config:    cfg,
		APIClient: NewAPIClient(cfg.URL),
	}
}

// GetTopics Get a list of Kafka topics.
func (krc *KafkaRestClient) GetTopics() ([]string, error) {
	responseBody := []string{}
	_, err := krc.Request(http.MethodGet, "/topics", nil, &responseBody, http.StatusOK)
	return responseBody, err
}

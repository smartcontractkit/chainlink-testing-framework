# Kafka

This is a wrapper over HTTP client that can only return a list of topics from Kafka instance.

```go
client, err := NewKafkaRestClient(&NewKafkaRestClient{URL: "my-kafka-url"})
if err != nil {
    panic(err)
}

topis, err := client.GetTopics()
if err != nil {
    panic(err)
}

for _, topic := range topics {
    fmt.Println("topic: " + topic)
}
```
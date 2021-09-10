package sink

const (
	SINK_KAFKA SinkType = "kafka"
)

type KafkaSink struct {
	Brokers []string `json:"brokers"`
	Topic   string   `json:"topic"`
}

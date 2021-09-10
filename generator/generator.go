package generator

type StatusType string

const (
	STATUS_INIT    StatusType = "init"
	STATUS_RUNNING StatusType = "running"
	STATUS_STOPPED StatusType = "stopped"
	STATUS_FAILED  StatusType = "failed"
)

type SinkType string

const (
	SINK_KAFKA StatusType = "kafka"
)

type GeneratorSink struct {
	Name   string                 `json:"name"`
	Type   SinkType               `json:"type"`
	Config map[string]interface{} `json:"config"`
}

type KafkaSink struct {
	Brokers []string `json:"brokers"`
	Topic   string   `json:"topic"`
}

type GeneratorSource struct {
	Name string `json:"name"`
}

type GeneratorConfig struct {
	Name   string          `json:"name"`
	Sink   GeneratorSink   `json:"sink"`
	Source GeneratorSource `json:"source"`
}

type Generator struct {
	Name   string
	Status StatusType
	Config GeneratorConfig
}

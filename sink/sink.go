package sink

type SinkType string

type GeneratorSink struct {
	Name   string                 `json:"name"`
	Type   SinkType               `json:"type"`
	Config map[string]interface{} `json:"config"`
}

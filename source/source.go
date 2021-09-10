package source

import (
	"log"
)

type SourceFieldType string

const (
	FIELDTYPE_TIMESTAMP SourceFieldType = "timestamp"
	FIELDTYPE_STRING    SourceFieldType = "string"
	FIELDTYPE_INT32     SourceFieldType = "int"
	FIELDTYPE_FLOAT32   SourceFieldType = "float"
)

type SourceField struct {
	Name  string          `json:"name"`
	Type  SourceFieldType `json:"type"`
	Range []interface{}   `json:"range,omitempty"`
	Limit []interface{}   `json:"limit,omitempty"`
}

type SourceConfiguration struct {
	Name           string        `json:"name"`
	TimestampField string        `json:"timestamp_field"`
	BatchSize      int           `json:"batch_size"`
	Concurrency    int           `json:"concurrency"`
	Internval      []int         `json:"interval"`
	Fields         []SourceField `json:"fields"`
}

type Event struct {
	Key   string                 `json:"key"`
	Value map[string]interface{} `json:"value"`
}

type EventGenerator struct {
	Config       SourceConfiguration
	EventChannel chan Event
}

func NewEventGenerator(source *SourceConfiguration) *EventGenerator {
	eventChan := make(chan Event, 100)
	generator := EventGenerator{
		Config:       *source,
		EventChannel: eventChan,
	}

	return &generator
}

func (s *EventGenerator) run() {
	event := Event{
		Key:   "somekey",
		Value: map[string]interface{}{},
	}
	s.EventChannel <- event
}

func (s *EventGenerator) Stop() {
	close(s.EventChannel)
}

func (s *EventGenerator) Run() {
	config := s.Config
	log.Printf("start to run event generator : %s \n", config.Name)

	for i := 0; i < config.Concurrency; i++ {
		go s.run()
	}
}

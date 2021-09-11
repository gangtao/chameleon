package source

import (
	"log"
	"sync"
	"time"

	fake "github.com/brianvoe/gofakeit/v6"
)

type SourceFieldType string

const (
	FIELDTYPE_TIMESTAMP SourceFieldType = "timestamp"
	FIELDTYPE_STRING    SourceFieldType = "string"
	FIELDTYPE_INT       SourceFieldType = "int"
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
	Stopped      bool
	mu           sync.Mutex
}

func NewEventGenerator(source *SourceConfiguration) *EventGenerator {
	eventChan := make(chan Event, 100)
	generator := EventGenerator{
		Config:       *source,
		EventChannel: eventChan,
		Stopped:      false,
	}

	return &generator
}

func makeTimestamp() int64 {
	now := time.Now()
	nsec := now.UnixNano() / 1000
	return nsec
}

func makeInt(ranges *[]int, limits *[]int) int {
	range_length := len(*ranges)
	limit_length := len(*limits)

	faker := fake.New(0)

	if range_length > 0 {
		index := faker.Number(0, range_length-1)
		return (*ranges)[index]
	} else if limit_length > 1 {
		return faker.Number((*limits)[0], (*limits)[1])
	}

	return 0
}

func makeFloat(ranges *[]float32, limits *[]float32) float32 {
	range_length := len(*ranges)
	limit_length := len(*limits)
	faker := fake.New(0)

	if range_length > 0 {
		index := faker.Number(0, range_length-1)
		return (*ranges)[index]
	} else if limit_length > 1 {
		return faker.Float32Range((*limits)[0], (*limits)[1])
	}

	return 0.0
}

func makeString(ranges *[]string) string {
	range_length := len(*ranges)
	faker := fake.New(0)

	if range_length > 0 {
		return faker.RandomString(*ranges)
	}

	return faker.LetterN(8)
}

func makeValue(sourceType *SourceFieldType, sourceRange *[]interface{}, sourceLimit *[]interface{}) interface{} {
	switch s := *sourceType; s {
	case FIELDTYPE_TIMESTAMP:
		return makeTimestamp()
	case FIELDTYPE_STRING:
		ranges := make([]string, len(*sourceRange))
		for i := 0; i < len(*sourceRange); i++ {
			ranges[i] = (*sourceRange)[i].(string)
		}

		return makeString(&ranges)
	case FIELDTYPE_INT:
		ranges := make([]int, len(*sourceRange))
		for i := 0; i < len(*sourceRange); i++ {
			ranges[i] = (*sourceRange)[i].(int)
		}

		limits := make([]int, len(*sourceLimit))
		for i := 0; i < len(*sourceLimit); i++ {
			limits[i] = (*sourceLimit)[i].(int)
		}
		return makeInt(&ranges, &limits)
	case FIELDTYPE_FLOAT32:
		ranges := make([]float32, len(*sourceRange))
		for i := 0; i < len(*sourceRange); i++ {
			ranges[i] = (*sourceRange)[i].(float32)
		}

		limits := make([]float32, len(*sourceLimit))
		for i := 0; i < len(*sourceLimit); i++ {
			limits[i] = (*sourceLimit)[i].(float32)
		}
		return makeFloat(&ranges, &limits)
	default:
		return nil
	}

}

func (s *EventGenerator) generateEvent() *Event {
	value := make(map[string]interface{})
	fields := s.Config.Fields
	faker := fake.New(0)

	for _, f := range fields {
		value[f.Name] = makeValue(&f.Type, &f.Range, &f.Limit)
	}

	event := Event{
		Key:   faker.LetterN(8),
		Value: value,
	}

	return &event
}

func (s *EventGenerator) generateBatchEvent() []*Event {
	batchSize := s.Config.BatchSize
	events := make([]*Event, batchSize)

	for i := 0; i < batchSize; i++ {
		events[i] = s.generateEvent()
	}
	return events
}

func (s *EventGenerator) run() {
	for {
		if s.Stopped {
			break
		}

		for _, v := range s.generateBatchEvent() {
			s.mu.Lock()
			if !s.Stopped {
				s.EventChannel <- *v
			}
			s.mu.Unlock()
		}

		// TODO: use interval in configuration
		time.Sleep(time.Millisecond * 1)
	}
}

func (s *EventGenerator) Stop() {
	s.mu.Lock()
	s.Stopped = true
	close(s.EventChannel)
	s.mu.Unlock()
}

func (s *EventGenerator) Run() {
	config := s.Config
	log.Printf("start to run event generator : %s \n", config.Name)

	for i := 0; i < config.Concurrency; i++ {
		go s.run()
	}
}

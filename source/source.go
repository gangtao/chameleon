package source

import (
	"log"
	"sync"
	"time"

	fake "github.com/brianvoe/gofakeit/v6"
)

type SourceFieldType string
type TimestampFormatType string

const (
	FIELDTYPE_TIMESTAMP SourceFieldType = "timestamp"
	FIELDTYPE_STRING    SourceFieldType = "string"
	FIELDTYPE_INT       SourceFieldType = "int"
	FIELDTYPE_FLOAT     SourceFieldType = "float"
)

type SourceField struct {
	Name              string          `json:"name"`
	Type              SourceFieldType `json:"type"`
	Range             []interface{}   `json:"range,omitempty"`
	Limit             []interface{}   `json:"limit,omitempty"`
	TimestampFormat   string          `json:"timestamp_format,omitempty"`
	TimestampDelayMin int             `json:"timestamp_delay_min,omitempty"`
	TimestampDelayMax int             `json:"timestamp_delay_max,omitempty"`
}

type SourceConfiguration struct {
	Name           string        `json:"name,omitempty"`
	TimestampField string        `json:"timestamp_field,omitempty"`
	BatchSize      int           `json:"batch_size,omitempty"`
	Concurrency    int           `json:"concurrency,omitempty"`
	Internval      []int         `json:"interval,omitempty"`
	Fields         []SourceField `json:"fields,omitempty"`
}

type Event struct {
	Key   string                 `json:"key"`
	Value map[string]interface{} `json:"value"`
}

type EventGenerator struct {
	Config       SourceConfiguration
	EventChannel chan []*Event
	Stopped      bool
	Counter      int
	mu           sync.Mutex
}

func NewEventGenerator(source *SourceConfiguration) *EventGenerator {
	eventChan := make(chan []*Event, 100)
	generator := EventGenerator{
		Config:       *source,
		EventChannel: eventChan,
		Stopped:      false,
		Counter:      0,
	}

	return &generator
}

func makeTimestamp(timestampDeleyMin int, timestampDeleyMax int) int64 {
	now := time.Now()
	faker := fake.New(0)
	delay := faker.Number(timestampDeleyMin, timestampDeleyMax)
	nsec := now.UnixNano() / 1000 // using micro second as unit
	return nsec - int64(delay)
}

func makeTimestampString(format *string, timestampDeleyMin int, timestampDeleyMax int) string {
	now := time.Now()
	nsec := now.UnixNano() / 1000
	faker := fake.New(0)
	delay := faker.Number(timestampDeleyMin, timestampDeleyMax)

	log.Printf("timestamp with delay: %d \n", delay)

	t := nsec - int64(delay)
	timestamp := time.UnixMicro(t)

	return timestamp.Format(*format)
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

func makeValue(sourceType *SourceFieldType, sourceRange *[]interface{}, sourceLimit *[]interface{},
	timestampFormat *string, TimestampDelayMin int, TimestampDelayMax int) interface{} {
	switch s := *sourceType; s {
	case FIELDTYPE_TIMESTAMP:
		//log.Printf("timestamp with format : %s \n", *timestampFormat)
		//log.Printf("timestamp with delay min: %d \n", TimestampDelayMin)
		//log.Printf("timestamp with delay max: %d \n", TimestampDelayMax)

		if *timestampFormat == "int" || *timestampFormat == "" {
			return makeTimestamp(TimestampDelayMin, TimestampDelayMax)
		} else {
			return makeTimestampString(timestampFormat, TimestampDelayMin, TimestampDelayMax)
		}

	case FIELDTYPE_STRING:
		ranges := make([]string, len(*sourceRange))
		for i := 0; i < len(*sourceRange); i++ {
			ranges[i] = (*sourceRange)[i].(string)
		}

		return makeString(&ranges)
	case FIELDTYPE_INT:
		ranges := make([]int, len(*sourceRange))
		for i := 0; i < len(*sourceRange); i++ {
			ranges[i] = int((*sourceRange)[i].(float64))
		}

		limits := make([]int, len(*sourceLimit))
		for i := 0; i < len(*sourceLimit); i++ {
			limits[i] = int((*sourceLimit)[i].(float64))
		}
		return makeInt(&ranges, &limits)
	case FIELDTYPE_FLOAT:
		ranges := make([]float32, len(*sourceRange))
		for i := 0; i < len(*sourceRange); i++ {
			ranges[i] = float32((*sourceRange)[i].(float64))
		}

		limits := make([]float32, len(*sourceLimit))
		for i := 0; i < len(*sourceLimit); i++ {
			limits[i] = float32((*sourceLimit)[i].(float64))
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
		value[f.Name] = makeValue(&f.Type, &f.Range, &f.Limit, &f.TimestampFormat, f.TimestampDelayMin, f.TimestampDelayMax)
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
	intervalRange := s.Config.Internval
	interval := time.Duration(1) // default to 1 micro second
	useRandomInterval := false
	faker := fake.New(0)

	if len(intervalRange) == 1 {
		interval = time.Duration(intervalRange[0])
	} else if len(intervalRange) == 2 {
		useRandomInterval = true
	}

	for {
		if s.Stopped {
			break
		}

		events := s.generateBatchEvent()
		s.mu.Lock()
		if !s.Stopped {
			s.EventChannel <- events
			s.Counter += len(events)
		}
		s.mu.Unlock()

		if useRandomInterval {
			rInterval := faker.Number(int(intervalRange[0]), int(intervalRange[1]))
			time.Sleep(time.Microsecond * time.Duration(rInterval))
		} else {
			time.Sleep(time.Microsecond * interval)
		}
	}
}

func (s *EventGenerator) Stop() {
	s.mu.Lock()
	if !s.Stopped {
		s.Stopped = true
		close(s.EventChannel)
	}
	s.mu.Unlock()
}

func (s *EventGenerator) Run() {
	config := s.Config
	log.Printf("start to run event generator : %s \n", config.Name)

	for i := 0; i < config.Concurrency; i++ {
		go s.run()
	}
}

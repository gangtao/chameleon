package sink

import (
	"chameleon/source"
	"context"
	"encoding/json"
	"log"

	"github.com/mitchellh/mapstructure"
	kafka "github.com/segmentio/kafka-go"
)

const (
	SINK_KAFKA SinkType = "kafka"
)

type KafkaSinkConfiguration struct {
	Brokers []string `json:"brokers"`
	Topic   string   `json:"topic"`
}

type KafkaSink struct {
	SinkConfiguration
	KafkaConfig KafkaSinkConfiguration
	Writer      *kafka.Writer
	Counter     int
}

func NewKafkaSink(config *SinkConfiguration) *KafkaSink {

	var kafkaConfig KafkaSinkConfiguration
	mapstructure.Decode(config.Config, &kafkaConfig)

	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  kafkaConfig.Brokers,
		Topic:    kafkaConfig.Topic,
		Balancer: &kafka.Hash{},
	})

	result := KafkaSink{
		SinkConfiguration: SinkConfiguration{
			Name:   config.Name,
			Type:   config.Type,
			Config: config.Config,
		},
		KafkaConfig: kafkaConfig,
		Writer:      writer,
		Counter:     0,
	}

	return &result
}

func (s *KafkaSink) Write(events *[]*source.Event) error {
	messages := make([]kafka.Message, len(*events))

	for i := 0; i < len(*events); i++ {
		event := (*events)[i]
		key := event.Key
		value, _ := json.Marshal(event.Value)
		messages[i] = kafka.Message{
			Key:   []byte(key),
			Value: []byte(value),
		}
	}

	err := s.Writer.WriteMessages(context.Background(), messages...)

	if err != nil {
		log.Println("failed to write messages {}", err)

	} else {
		s.Counter += len(*events)
	}

	return err
}

func (s *KafkaSink) Count() int {
	return s.Counter
}

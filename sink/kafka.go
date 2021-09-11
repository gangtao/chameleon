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
}

func NewKafkaSink(config *SinkConfiguration) *KafkaSink {

	var kafkaConfig KafkaSinkConfiguration
	mapstructure.Decode(config.Config, &kafkaConfig)

	result := KafkaSink{
		SinkConfiguration: SinkConfiguration{
			Name:   config.Name,
			Type:   config.Type,
			Config: config.Config,
		},
		KafkaConfig: kafkaConfig,
	}

	return &result
}

func (s *KafkaSink) Write(events *[]*source.Event) error {
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  s.KafkaConfig.Brokers,
		Topic:    s.KafkaConfig.Topic,
		Balancer: &kafka.Hash{},
	})

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

	log.Println("write events in kafka")

	err := w.WriteMessages(context.Background(), messages...)

	if err != nil {
		log.Fatal("failed to write messages:", err)
	}

	return err
}

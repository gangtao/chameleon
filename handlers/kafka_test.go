package handlers

import (
	"context"
	"log"
	"time"

	kafka "github.com/segmentio/kafka-go"

	. "github.com/onsi/ginkgo"
)

var _ = Describe("Test", func() {

	Describe("Kafka Producer", func() {

		It("Test connection", func() {
			topic := "topic-A"
			partition := 0

			conn, err := kafka.DialLeader(context.Background(), "tcp", "localhost:9092", topic, partition)
			if err != nil {
				log.Fatal("failed to dial leader:", err)
			}

			conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			_, err = conn.WriteMessages(
				kafka.Message{Value: []byte("one!")},
				kafka.Message{Value: []byte("two!")},
				kafka.Message{Value: []byte("three!")},
			)
			if err != nil {
				log.Fatal("failed to write messages in connect:", err)
			}

			if err := conn.Close(); err != nil {
				log.Fatal("failed to close writer:", err)
			}
		})

		It("produce some message", func() {
			w := kafka.NewWriter(kafka.WriterConfig{
				Brokers:  []string{"localhost:9092"},
				Topic:    "topic-A",
				Balancer: &kafka.Hash{},
			})

			err := w.WriteMessages(context.Background(),
				kafka.Message{
					Key:   []byte("Key-A"),
					Value: []byte("Hello World!"),
				},
				kafka.Message{
					Key:   []byte("Key-B"),
					Value: []byte("One!"),
				},
				kafka.Message{
					Key:   []byte("Key-C"),
					Value: []byte("Two!"),
				},
			)

			if err != nil {
				log.Fatal("failed to write messages:", err)
			}

			if err := w.Close(); err != nil {
				log.Fatal("failed to close writer:", err)
			}

		})
	})
})

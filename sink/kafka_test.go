package sink_test

import (
	"context"
	"log"
	"time"

	kafka "github.com/segmentio/kafka-go"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
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

			Expect(err).ShouldNot(HaveOccurred())

			if err := w.Close(); err != nil {
				log.Fatal("failed to close writer:", err)
			}

		})

		It("read some message", func() {
			r := kafka.NewReader(kafka.ReaderConfig{
				Brokers:   []string{"localhost:9092", "localhost:9093", "localhost:9094"},
				Topic:     "topic-A",
				Partition: 0,
				MinBytes:  10e3, // 10KB
				MaxBytes:  10e6, // 10MB
			})
			r.SetOffset(0)

			for i := 0; i < 10; i++ {
				m, err := r.ReadMessage(context.Background())
				if err != nil {
					break
				}
				log.Printf("message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))
			}

			if err := r.Close(); err != nil {
				log.Fatal("failed to close reader:", err)
			}

		})
	})
})

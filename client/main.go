package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	kafka "github.com/segmentio/kafka-go"
)

func main() {
	broker := "localhost:9092"
	topic := "topic-A"

	if len(os.Args) > 1 {
		broker = os.Args[1]
	}
	log.Printf("Using broker %s", broker)

	if len(os.Args) > 2 {
		topic = os.Args[2]
	}
	log.Printf("Using topic broker %s", topic)

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{broker},
		Topic:     topic,
		Partition: 0,
		MinBytes:  10e3, // 10KB
		MaxBytes:  10e6, // 10MB
	})
	r.SetOffset(kafka.LastOffset)

	for {

		m, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Fatalf("failed to read messages with '%s'\n", err)
		}
		now := time.Now()
		//log.Printf("message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))
		data := make(map[string]interface{})
		json.Unmarshal(m.Value, &data)
		//fmt.Printf("Operation: %v", data)

		log.Printf("message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))
		log.Printf("t   is {}", int64(data["t"].(float64)))
		log.Printf("now is {}", now.UnixNano()/1000)

		ts := int64(data["t"].(float64))
		nsec := now.UnixNano() / 1000

		latency := nsec - ts
		log.Printf("Latency is : %d : %d \n", m.Offset, latency)

	}

	if err := r.Close(); err != nil {
		log.Fatal("failed to close reader:", err)
	}
}

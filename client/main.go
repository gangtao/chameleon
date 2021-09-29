package main

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	kafka "github.com/segmentio/kafka-go"
)

func writeResult(data [][]string) {
	newpath := filepath.Join(".", "data")
	err := os.MkdirAll(newpath, os.ModePerm)
	checkError("Cannot create dir", err)

	file, err := os.Create("./data/result.csv")
	checkError("Cannot create file", err)
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, value := range data {
		err := writer.Write(value)
		checkError("Cannot write to file", err)
	}
}

func checkError(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
}

func main() {
	broker := "localhost:9092"
	topic := "topic-A"
	size := 10000
	result := make([][]string, size+1)

	if len(os.Args) > 1 {
		broker = os.Args[1]
	}
	log.Printf("Using broker %s", broker)

	if len(os.Args) > 2 {
		topic = os.Args[2]
	}
	log.Printf("Using topic broker %s", topic)

	if len(os.Args) > 3 {
		size, _ = strconv.Atoi(os.Args[3])
	}
	log.Printf("Totoal size is %d", size)

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{broker},
		Topic:     topic,
		Partition: 0,
		MinBytes:  10e3, // 10KB
		MaxBytes:  10e6, // 10MB
	})
	r.SetOffset(kafka.LastOffset)

	result[0] = []string{"broker", "offset", "latency"}

	for i := 1; i <= size; i++ {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Fatalf("failed to read messages with '%s'\n", err)
		}
		now := time.Now()
		//log.Printf("message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))
		data := make(map[string]interface{})
		json.Unmarshal(m.Value, &data)
		//fmt.Printf("Operation: %v", data)

		//log.Printf("message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))
		//log.Printf("t   is {}", int64(data["t"].(float64)))
		//log.Printf("now is {}", now.UnixNano()/1000)

		ts := int64(data["t"].(float64))
		nsec := now.UnixNano() / 1000

		latency := nsec - ts
		//log.Printf("Latency is : %d : %d \n", m.Offset, latency)
		result[i] = []string{broker, strconv.FormatInt(m.Offset, 10), strconv.FormatInt(latency, 10)}
	}

	writeResult(result)

	if err := r.Close(); err != nil {
		log.Fatal("failed to close reader:", err)
	}
}

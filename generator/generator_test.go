package generator_test

import (
	"encoding/json"
	"log"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"chameleon/generator"
	"github.com/mitchellh/mapstructure"
)

var _ = Describe("Generator", func() {

	It("Test save generator configuration", func() {
		kafka_config := generator.KafkaSink{
			Brokers: []string{"localhost:9092"},
			Topic:   "topic_a",
		}

		// following code convert stuct to interface map
		var in map[string]interface{}
		inrec, err := json.Marshal(kafka_config)
		Expect(err).ShouldNot(HaveOccurred())
		json.Unmarshal(inrec, &in)

		config := generator.GeneratorConfig{
			Name: "testconfig",
			Sink: generator.GeneratorSink{
				Name:   "sinkname",
				Type:   "kafka",
				Config: in,
			},
			Source: generator.GeneratorSource{
				Name: "sourcename",
			},
		}

		log.Printf("%v", config)

		config_marshalled, err := json.Marshal(config)
		Expect(err).ShouldNot(HaveOccurred())

		log.Printf("%s", config_marshalled)

	})

	It("Test load generator configuration", func() {
		congig_string := `{
			"name":"testconfig",
			"sink":{
			   "name":"sinkname",
			   "type":"kafka",
			   "config":{
				  "brokers":[
					 "localhost:9092"
				  ],
				  "topic":"topic_a"
			   }
			},
			"source":{
			   "name":"sourcename"
			}
		 }`

		res := generator.GeneratorConfig{}
		err := json.Unmarshal([]byte(congig_string), &res)
		Expect(err).ShouldNot(HaveOccurred())

		log.Printf("%v", res)

		kafka_sink := res.Sink.Config
		log.Printf("%v", kafka_sink)

		// using mapstructure.Decode to convert interface map to structure
		var result generator.KafkaSink
		err = mapstructure.Decode(kafka_sink, &result)
		Expect(err).ShouldNot(HaveOccurred())

		log.Printf("%#v", result)

	})

})

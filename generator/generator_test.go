package generator_test

import (
	"encoding/json"
	"log"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"chameleon/generator"
	"chameleon/sink"
	"chameleon/source"
	"github.com/mitchellh/mapstructure"
)

var _ = Describe("Generator", func() {

	It("Test save generator configuration", func() {
		kafka_config := sink.KafkaSink{
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
			Sink: sink.GeneratorSink{
				Name:   "sinkname",
				Type:   "kafka",
				Config: in,
			},
			Source: source.GeneratorSource{
				Name:           "sourcename",
				TimestampField: "f1",
				BatchSize:      100,
				Internval:      []int{0, 100},
				Fields: []source.SourceField{
					source.SourceField{
						Name:  "f1",
						Type:  source.FIELDTYPE_TIMESTAMP,
						Range: []interface{}{1, 3},
					},
					source.SourceField{
						Name:  "f2",
						Type:  source.FIELDTYPE_STRING,
						Range: []interface{}{"a", "b", "c"},
					},
					source.SourceField{
						Name:  "f3",
						Type:  source.FIELDTYPE_INT32,
						Limit: []interface{}{0, 100},
					},
				},
			},
		}

		log.Printf("%v", config)

		config_marshalled, err := json.Marshal(config)
		Expect(err).ShouldNot(HaveOccurred())

		log.Printf("%s", config_marshalled)

	})

	It("Test load generator configuration", func() {
		congig_string := `
		{
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
			   "name":"sourcename",
			   "timestamp_field":"f1",
			   "batch_size":100,
			   "interval":[
				  0,
				  100
			   ],
			   "fields":[
				  {
					 "name":"f1",
					 "type":"timestamp",
					 "range":[
						1,
						3
					 ]
				  },
				  {
					 "name":"f2",
					 "type":"string",
					 "range":[
						"a",
						"b",
						"c"
					 ]
				  },
				  {
					 "name":"f3",
					 "type":"int",
					 "limit":[
						0,
						100
					 ]
				  }
			   ]
			}
		 }`

		res := generator.GeneratorConfig{}
		err := json.Unmarshal([]byte(congig_string), &res)
		Expect(err).ShouldNot(HaveOccurred())

		log.Printf("%v", res)

		kafka_sink := res.Sink.Config
		log.Printf("%v", kafka_sink)

		// using mapstructure.Decode to convert interface map to structure
		var result sink.KafkaSink
		err = mapstructure.Decode(kafka_sink, &result)
		Expect(err).ShouldNot(HaveOccurred())

		log.Printf("%#v", result)

	})

	It("Test create generator", func() {
		congig_string := `
		{
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
			   "name":"sourcename",
			   "timestamp_field":"f1",
			   "batch_size":100,
			   "interval":[
				  0,
				  100
			   ],
			   "fields":[
				  {
					 "name":"f1",
					 "type":"timestamp",
					 "range":[
						1,
						3
					 ]
				  },
				  {
					 "name":"f2",
					 "type":"string",
					 "range":[
						"a",
						"b",
						"c"
					 ]
				  },
				  {
					 "name":"f3",
					 "type":"int",
					 "limit":[
						0,
						100
					 ]
				  }
			   ]
			}
		 }`

		res := generator.GeneratorConfig{}
		err := json.Unmarshal([]byte(congig_string), &res)
		Expect(err).ShouldNot(HaveOccurred())

		g := generator.NewGenerator(&res)
		Expect(g).ShouldNot(BeNil())
		log.Printf("%v", g)

	})

})

package source_test

import (
	"log"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"chameleon/source"
)

var _ = Describe("Source", func() {
	It("Generator case 1", func() {
		config := source.SourceConfiguration{
			Name:           "testConfig",
			TimestampField: "t",
			BatchSize:      100,
			Concurrency:    3,
			Internval:      []int{100},
			Fields: []source.SourceField{
				source.SourceField{
					Name: "t",
					Type: source.FIELDTYPE_TIMESTAMP,
				},
				source.SourceField{
					Name:  "f1",
					Type:  source.FIELDTYPE_STRING,
					Range: []interface{}{"a", "b", "c"},
				},
				source.SourceField{
					Name:  "f2",
					Type:  source.FIELDTYPE_INT32,
					Limit: []interface{}{0, 100},
				},
			},
		}

		generator := source.NewEventGenerator(&config)
		Expect(generator).ShouldNot(BeNil())

		generator.Run()

		go func() {
			time.Sleep(time.Second * 1)
			generator.Stop()
		}()

		for {
			event, ok := <-generator.EventChannel
			log.Printf("%v, %v \n", event, ok)

			if !ok {
				break
			}

			time.Sleep(time.Second * 1)
		}
	})

})

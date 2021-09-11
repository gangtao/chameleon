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
			BatchSize:      2,
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
					Range: []interface{}{"aa", "bb", "cc"},
				},
				source.SourceField{
					Name:  "f2",
					Type:  source.FIELDTYPE_INT,
					Range: []interface{}{float64(0), float64(10), float64(100)},
				},
				source.SourceField{
					Name:  "f3",
					Type:  source.FIELDTYPE_FLOAT,
					Limit: []interface{}{float64(0.1), float64(100.0)},
				},
			},
		}

		generator := source.NewEventGenerator(&config)
		Expect(generator).ShouldNot(BeNil())

		generator.Run()

		go func() {
			time.Sleep(time.Millisecond * 300)
			generator.Stop()
		}()

		for {
			event, ok := <-generator.EventChannel
			log.Printf("%v, %v \n", event, ok)
			if !ok {
				break
			}
			time.Sleep(time.Millisecond * 1)
		}
	})

})

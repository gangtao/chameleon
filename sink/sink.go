package sink

import (
	"chameleon/source"
)

type SinkType string

type SinkConfiguration struct {
	Name   string                 `json:"name"`
	Type   SinkType               `json:"type"`
	Config map[string]interface{} `json:"config"`
}

type Writer interface {
	Write(events *[]*source.Event) error
	Count() int
}

func WriteEvents(writer *Writer, events *[]*source.Event) error {
	return (*writer).Write(events)
}

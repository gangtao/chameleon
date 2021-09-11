package generator

import (
	"chameleon/sink"
	"chameleon/source"
	"log"
	"time"

	uuid "github.com/satori/go.uuid"
)

type StatusType string

const (
	STATUS_INIT    StatusType = "init"
	STATUS_RUNNING StatusType = "running"
	STATUS_STOPPED StatusType = "stopped"
	STATUS_FAILED  StatusType = "failed"
)

type GeneratorConfig struct {
	Name   string                     `json:"name"`
	Sink   sink.SinkConfiguration     `json:"sink"`
	Source source.SourceConfiguration `json:"source"`
}

type Generator struct {
	Id     string
	Status StatusType
	Config GeneratorConfig
	Source *source.EventGenerator
	Sink   sink.Writer
}

func NewGenerator(config *GeneratorConfig) *Generator {
	kafkaSink := sink.NewKafkaSink(&config.Sink)
	generator := Generator{
		Id:     uuid.NewV4().String(),
		Status: STATUS_INIT,
		Config: *config,
		Source: source.NewEventGenerator(&config.Source),
		Sink:   kafkaSink,
	}

	return &generator
}

func (g *Generator) Run() error {
	g.Source.Run()

	// Test to run for a while
	go func() {
		time.Sleep(time.Millisecond * 300)
		g.Source.Stop()
	}()

	for {
		events, ok := <-g.Source.EventChannel
		log.Printf("%v, %v \n", events, ok)
		if !ok {
			break
		}
		//TODO: cache events to send
		g.Sink.Write(&events)
		// go func() {
		// 	g.Sink.Write(&events)
		// }()
		time.Sleep(time.Millisecond * 1)
	}

	return nil
}

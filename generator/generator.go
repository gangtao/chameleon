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
	//TODO: support other sink type as well
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

func (g *Generator) Run(timeout time.Duration) error {
	g.Status = STATUS_RUNNING
	g.Source.Run()

	// Test to run for a while
	if timeout != 0 {
		go func() {
			time.Sleep(time.Microsecond * timeout)
			g.Source.Stop()
			g.Status = STATUS_STOPPED
		}()
	}

	for {
		events, ok := <-g.Source.EventChannel
		log.Printf("%v, %v \n", events, ok)
		if !ok {
			log.Printf("generator stopped, with total generated event %d, sink event %d", g.Source.Counter, g.Sink.Count())
			break
		}

		// send data in coroutine, do not block the channel
		go func() {
			g.Sink.Write(&events)
		}()
		//time.Sleep(time.Millisecond * 1)
	}

	return nil
}

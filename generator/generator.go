package generator

import (
	"chameleon/sink"
	"chameleon/source"

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
	Sink   sink.GeneratorSink         `json:"sink"`
	Source source.SourceConfiguration `json:"source"`
}

type Generator struct {
	Id     string
	Status StatusType
	Config GeneratorConfig
}

func NewGenerator(config *GeneratorConfig) *Generator {
	generator := Generator{
		Id:     uuid.NewV4().String(),
		Status: STATUS_INIT,
		Config: *config,
	}

	return &generator
}

func (g *Generator) Run() error {
	return nil
}

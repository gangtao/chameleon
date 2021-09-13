package generator

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"
)

const CONFIG_DIR = "config"

//TODO : should store pointer to the generator instead of the object
type GeneratorManager struct {
	Manager map[string]*Generator
}

var lock = &sync.Mutex{}
var (
	instance *GeneratorManager
)

type GeneratorStatus struct {
	Status      StatusType `json:"status"`
	SourceCount int        `json:"source_count"`
	SinkCount   int        `json:"sink_count"`
}

func NewGeneratorManager() *GeneratorManager {
	lock.Lock()
	defer lock.Unlock()

	if instance == nil {
		instance = &GeneratorManager{
			Manager: make(map[string]*Generator), // <-- thread safe
		}

		// load all generator from file
		files, err := ioutil.ReadDir(CONFIG_DIR)
		if err == nil {
			for _, file := range files {
				if !file.IsDir() {
					jsonFile, err := os.Open(filepath.Join(CONFIG_DIR, file.Name()))
					defer jsonFile.Close()
					if err == nil {
						byteValue, _ := ioutil.ReadAll(jsonFile)
						var config GeneratorConfig
						json.Unmarshal(byteValue, &config)
						instance.CreateGenerator(&config)
					}
				}
			}
		}
	}
	return instance
}

func (gm *GeneratorManager) CreateGenerator(config *GeneratorConfig) error {
	g := NewGenerator(config)

	name := g.Config.Name
	_, exists := gm.Manager[name]
	if exists {
		return errors.New("Generator name already exists")
	}

	//TODO : validate the name should be a valid file name

	gm.Manager[name] = g

	//save config to config folder
	path := fmt.Sprintf("%s/%s.json", CONFIG_DIR, name)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		configJson, err := json.MarshalIndent(*config, "", " ")
		if err == nil {
			os.MkdirAll(CONFIG_DIR, 0700)
			_ = ioutil.WriteFile(path, configJson, 0644)
		}
	}

	return nil
}

func (gm *GeneratorManager) ListGenerator() []string {
	names := make([]string, len(gm.Manager))
	i := 0
	for k := range gm.Manager {
		names[i] = k
		i++
	}

	return names
}

func (gm *GeneratorManager) DeleteGenerator(name string) error {
	_, exists := gm.Manager[name]
	if exists {
		delete(gm.Manager, name)
		path := fmt.Sprintf("%s/%s.json", CONFIG_DIR, name)
		// delete the file as well
		os.Remove(path)
		return nil
	}

	return errors.New("Generator does not exist")
}

func (gm *GeneratorManager) StartGenerator(name string) error {
	g, exists := gm.Manager[name]
	if exists {
		switch status := g.Status; status {
		case STATUS_INIT:
			go func() {
				g.Run(1000 * 1000)
			}()
			return nil
		case STATUS_RUNNING:
			log.Println("the generator is in running state")
			return nil
		default:
			new_generator := NewGenerator(&g.Config)
			gm.Manager[name] = new_generator
			go func() {
				new_generator.Run(1000 * 1000)
			}()
			return nil
			return nil
		}

	}
	return errors.New("Generator does not exist")
}

func (gm *GeneratorManager) StopGenerator(name string) error {
	g, exists := gm.Manager[name]
	if exists {
		go func() {
			g.Stop()
		}()
		return nil
	}

	return errors.New("Generator does not exist")
}

func (gm *GeneratorManager) GetGeneratorStatus(name string) (*GeneratorStatus, error) {
	g, exists := gm.Manager[name]
	if exists {
		status := GeneratorStatus{
			Status:      g.Status,
			SourceCount: g.Source.Counter,
			SinkCount:   g.Sink.Count(),
		}
		return &status, nil
	}

	return nil, errors.New("Generator does not exist")
}

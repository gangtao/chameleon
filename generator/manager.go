package generator

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
)

const CONFIG_DIR = "config"

type GeneratorManager struct {
	Manager map[string]Generator
}

var lock = &sync.Mutex{}
var (
	instance *GeneratorManager
)

func NewGeneratorManager() *GeneratorManager {
	lock.Lock()
	defer lock.Unlock()

	if instance == nil {
		instance = &GeneratorManager{
			Manager: make(map[string]Generator), // <-- thread safe
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

	gm.Manager[name] = *g

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

func (gm *GeneratorManager) DeleteGenerator(name string) {
	delete(gm.Manager, name)
	path := fmt.Sprintf("%s/%s.json", CONFIG_DIR, name)
	// delete the file as well
	os.Remove(path)

}

package cmd

import (
	"flag"
	"io/ioutil"
	"log"
	"os"

	"github.com/s-rafaeldias/who-is-down/service"
	"gopkg.in/yaml.v2"
)

type Cli struct {
	configFilePath string
}

type Services map[string]map[string]string

func New() *Cli {
	return &Cli{}
}

func (c *Cli) Start() {
	flag.StringVar(&c.configFilePath, "configFile", "./config.yaml", "Path to configFile")
	flag.Parse()

	// parse file
	services := c.parseConfigFile()

	// create a supervisor and start watching the services
	supervisor := service.NewSupervisor(services)
	supervisor.Start()
}

func (c *Cli) parseConfigFile() []*service.Service {
	// open file
	file, err := os.Open(c.configFilePath)
	if err != nil {
		log.Panic(err)
	}
	defer file.Close()

	// read file
	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Panic(err)
	}

	// parse yaml
	var services Services
	err = yaml.Unmarshal([]byte(data), &services)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	// create a slice of service.Service
	servicesToWatch := make([]*service.Service, 0)
	for name, values := range services {
		servicesToWatch = append(servicesToWatch, service.NewService(name, values))
	}

	return servicesToWatch
}

package cmd

import (
	"flag"
	"io/ioutil"
	"log"
	"os"

	"github.com/s-rafaeldias/who-is-down/notification"
	"github.com/s-rafaeldias/who-is-down/service"
	"gopkg.in/yaml.v2"
)

type Cli struct {
	configFilePath string
}

type YamlFile map[string]map[string]string

// New creates a new Cli
func New() *Cli {
	return &Cli{}
}

// Start starts watching all services defined on `configFile`
func (c *Cli) Start() {
	flag.StringVar(&c.configFilePath, "configFile", "./config.yaml", "Path to configFile")
	flag.Parse()

	// parse file
	services := c.parseConfigFile()

	// TODO: add option to choose Notifier
	slack, err := notification.New()
	if err != nil {
		log.Panicln(err)
	}

	// create a supervisor and start watching the services
	supervisor := service.NewSupervisor(services, slack)
	supervisor.Start()
}

// parseConfigFile parses the configFile and return a slice of
// services to watch.
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
	var servicesFromConfig YamlFile
	err = yaml.Unmarshal([]byte(data), &servicesFromConfig)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	// create a slice of service.Service
	servicesToWatch := make([]*service.Service, 0)
	for name, values := range servicesFromConfig {
		servicesToWatch = append(servicesToWatch, service.NewService(name, values))
	}

	return servicesToWatch
}

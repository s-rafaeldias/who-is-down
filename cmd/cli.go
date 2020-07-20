package cmd

import (
	"flag"
	"io/ioutil"
	"log"
	"os"

	"github.com/s-rafaeldias/who-is-down/pkg"
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
	flag.StringVar(&c.configFilePath, "f", "config.yaml", "Path to config file")
	flag.Parse()

	// parse file
	services := c.parseConfigFile()

	// TODO: add option to choose Notifier
	slack, err := pkg.NewSlackClient()
	if err != nil {
		log.Panicln(err)
	}

	// create a supervisor and start watching the services
	supervisor := pkg.NewSupervisor(services, slack)
	supervisor.Start()
}

// parseConfigFile parses the configFile and return a slice of
// services to watch.
func (c *Cli) parseConfigFile() []*pkg.Service {
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
	servicesToWatch := make([]*pkg.Service, 0)
	for name, values := range servicesFromConfig {
		servicesToWatch = append(servicesToWatch, pkg.NewService(name, values))
	}

	return servicesToWatch
}

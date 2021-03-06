package cmd

import (
	"flag"
	"fmt"
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
func (c *Cli) Run() {
	installCmd := flag.NewFlagSet("install", flag.ExitOnError)

	startCmd := flag.NewFlagSet("start", flag.ExitOnError)
	configFilePath := startCmd.String("f", "config.yaml", "Path for config file")

	// checks for correct number of args
	if len(os.Args) < 2 {
		fmt.Println("Expected a command")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "install":
		installCmd.Parse(os.Args[2:])
		c.Install()
	case "start":
		startCmd.Parse(os.Args[2:])
		c.configFilePath = *configFilePath
		c.Start()
	default:
		flag.PrintDefaults()
		os.Exit(1)
	}

}

// Install configure all scripts needed for running ``who-is-down`` as a systemctl service.
func (c *Cli) Install() {
	// default paths
	// create folders if needed
	// ask
}

// Start initialize the supervisor
func (c *Cli) Start() {
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
		s, _ := pkg.NewService(name, values)
		servicesToWatch = append(servicesToWatch, s)
	}

	return servicesToWatch
}

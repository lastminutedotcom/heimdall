package properties

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// TBD
// Is it worth using an init() to parse the config from std location? Something like this is somewhat necessary
// Or should we let the library clients decide when to store and how to manipulate the Config through their application?
// A sane default would be to provide the singleton if the default file is found, log a WARN otherwise...

// Config is a singleton to fetch configuration data from
var Config *ConfigProps

// We leverage variable initialization happening BEFORE init() for all packages to load the Config singleton
var _ = func() {
	var err error
	Config, err = loadOrBuildDefaultConfig()
	if err != nil {
		log.Println("could not get default configuration from standard file, please provide one via code")
	}
}

const (
	// ErrPropertiesMissingApplicationName describes the missing name error
	ErrPropertiesMissingApplicationName = "application name cannot be empty"
)

// ConfigProps defines the application's configuration properties used to run the service
// As a minimum the configuration will have to provide an application name to pass the validation
// all other attributes are optional
type ConfigProps struct {
	filePath    string
	Application *ApplicationProps `yaml:"application" json:"application"`
}

// ApplicationProps defines the application properties for configuration
type ApplicationProps struct {
	Name       string `yaml:"name" json:"name"`
	Datasource struct {
		URL      string `yaml:"url" json:"url"`
		Username string `yaml:"username" json:"username"`
		Password string `yaml:"password" json:"password"`
	} `yaml:"datasource" json:"datasource"`
	Parameters map[string]interface{} `yaml:"parameters" json:"parameters"`
}

// ConfigsFromYaml returns the application properties loading from standard
// file location, it will return nil and error if the file is not found or YAML cannot be parsed
func ConfigsFromYaml(filePath string) (*ConfigProps, error) {
	var err error
	cfg := &ConfigProps{
		filePath: filePath,
	}
	if cfg, err = cfg.parse(); err != nil {
		return nil, fmt.Errorf("error parsing YAML configuration file: %v", err)
	}
	if err = cfg.validate(); err != nil {
		return nil, fmt.Errorf("configuration file is not valid: %v", err)
	}
	return cfg, nil
}

func (c *ConfigProps) parse() (*ConfigProps, error) {
	yamlContent, err := ioutil.ReadFile(c.filePath)
	if err != nil {
		return nil, err
	}
	if err = yaml.Unmarshal(yamlContent, &c); err != nil {
		return nil, err
	}
	return c, nil
}

func (c *ConfigProps) validate() error {
	if c.Application.Name == "" {
		return errors.New(ErrPropertiesMissingApplicationName)
	}
	return nil
}

func loadOrBuildDefaultConfig() (config *ConfigProps, err error) {
	if config, err = ConfigsFromYaml(filepath.Join("appfw", "config", "appfw-configuration.yaml")); err != nil {
		return &ConfigProps{
			Application: &ApplicationProps{
				Name: randomAppName(0),
			},
		}, nil
	}
	return config, nil
}

// AppName returns the application name or populates it with a random value
func AppName() string {
	if Config != nil && Config.Application != nil && Config.Application.Name != "" {
		return Config.Application.Name
	}
	Config = &ConfigProps{
		Application: &ApplicationProps{
			Name: randomAppName(0),
		},
	}
	return Config.Application.Name
}

// ConfigJSON prints the current configuration properties in JSON format
func ConfigJSON() ([]byte, error) {
	//TODO(fgulazzi) obfuscate sensible configs with asterisks
	return json.Marshal(Config)
}

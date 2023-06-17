package app

import (
	"encoding/json"
	"hash/adler32"
	"os"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/piotr-ku/yaml-runner-go/system"
	"gopkg.in/yaml.v3"
)

// File config.go defines data structures used for the configuration file.
//
// The file includes the following data structures:
//
// Fact: Provides a data format for the facts defined in the configuration file.
//   - Name: The name of the fact. It is a required field.
//   - Command: The command associated with the fact. It is a required field.
//   - Shell: Shell used to execute the command.
//
// Action: Provides a data format for the actions defined in the configuration file.
//   - Command: The command associated with the action. It is a required field.
//   - Rules: A slice of strings representing the rules associated with the action.
//   - Shell: Shell used to execute the command.
//
// Config: Provides a data format for the configuration file.
//   - Logging: Contains configuration settings for logging. It uses the system.LogConfig type.
//   - Facts: A slice of Fact objects representing the facts defined in the configuration file.
//   - Actions: A slice of Action objects representing the actions defined in the configuration file.
//   - Hash: A checksum value calculated based on the merged configuration data.
//
// The data structures make use of struct tags for validation purposes, ensuring that required fields are present.
// The validation tags are defined using the "validate" tag in struct fields.
//
// This file provides the necessary data structures for representing and validating a configuration file.

// Fact provides a data format for the facts defined
// in the configuration file.
type Fact struct {
	Name    string `validate:"required"` // fact name
	Command string `validate:"required"` // fact command
	Shell   string // fact shell
}

// Format provides a data format for the actions defined
// in the configuration file.
type Action struct {
	Command string   `validate:"required"` // action command
	Rules   []string // action rules
	Shell   string   // action shell
}

// Daemon provides a data format for daemon settings defined
// in the configuration file.
type Daemon struct {
	Interval string `validate:"duration"`
}

// Config provides a data format for the configuration file.
type Config struct {
	Daemon  Daemon           `validate:"dive"`
	Logging system.LogConfig `validate:"dive"`
	Facts   []Fact           `validate:"dive"`          // facts slice
	Actions []Action         `validate:"required,dive"` // actions slice
	Hash    uint32
}

// Merge merges the fields of the provided Config into the receiver Config.
func (c *Config) Merge(m Config) {
	// Merge Daemon fields
	if m.Daemon.Interval != "" {
		c.Daemon.Interval = m.Daemon.Interval
	}

	// Merge Logging fields
	if m.Logging.File != "" {
		c.Logging.File = m.Logging.File
	}
	if m.Logging.Level != "" {
		c.Logging.Level = m.Logging.Level
	}
	if m.Logging.Quiet {
		c.Logging.Quiet = m.Logging.Quiet
	}
	if m.Logging.Json {
		c.Logging.Json = m.Logging.Json
	}

	// Merge Facts
	if len(m.Facts) > 0 {
		c.Facts = append(c.Facts, m.Facts...)
	}

	// Merge Actions
	if len(m.Actions) > 0 {
		c.Actions = append(c.Actions, m.Actions...)
	}
}

// CalculateHash calculates a Adler-32 hash from the Config struct
func (c *Config) CalculateHash() {
	// ignore c.Hash from calculation
	c.Hash = 0
	// calculate a checksum
	jsonData, err := json.Marshal(c)
	// notest
	if err != nil {
		system.FatalError("ParseError", err.Error())
	}

	// create a new Adler-32 hash
	hash := adler32.New()
	hash.Write(jsonData)
	c.Hash = hash.Sum32()
}

// LoadConfigFile loads a configuration file, validates it, and returns
// the resulting Config.
func LoadConfigFile(file string) Config {
	// read configuration file
	configContent, err := os.ReadFile(file)
	// notest
	if err != nil {
		system.FatalError("IOError", err.Error())
	}

	// parse configuration file
	config, err := parseYaml(configContent)
	// notest
	if err != nil {
		system.FatalError("ParseError", err.Error())
	}

	// validate configuration file
	validate := validateConfig(config)
	// notest
	if validate != nil {
		system.FatalError("ValidationError", validate.Error())
	}

	return config
}

// parseYaml parses the provided YAML content into a Config struct and returns it.
// If an error occurs during unmarshaling, it is also returned.
func parseYaml(content []byte) (Config, error) {
	var structure Config
	err := yaml.Unmarshal([]byte(content), &structure)
	return structure, err
}

// DurationValidator is a custom validator for duration strings.
type DurationValidator struct {
	validator *validator.Validate
}

// newDurationValidator creates a new instance of DurationValidator.
func newDurationValidator() *DurationValidator {
	return &DurationValidator{
		validator: validator.New(),
	}
}

// Validate is the validation method for duration strings.
// It checks if the duration string is valid by attempting to parse it using time.ParseDuration().
func (v *DurationValidator) Validate(fl validator.FieldLevel) bool {
	durationStr := fl.Field().String()
	// field is not required
	if durationStr == "" {
		return true
	}
	_, err := time.ParseDuration(durationStr)
	return err == nil
}

// validateConfig validates the provided Config object using a validator and returns any validation errors encountered.
// If the configuration is valid, it returns nil.
func validateConfig(config Config) error {
	// Create a new instance of DurationValidator.
	v := newDurationValidator()

	// Create a validator instance.
	validate := v.validator

	// Register the custom validation function "duration" with the validator.
	validate.RegisterValidation("duration", v.Validate)

	return validate.Struct(config)
}

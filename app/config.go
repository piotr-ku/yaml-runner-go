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
// Config: Provides a data format for the configuration file.
//   - Logging: Contains configuration settings for logging. It uses
// the system.LogConfig type.
//   - Facts: A slice of Fact objects representing the facts defined in
// the configuration file.
//   - Actions: A slice of Action objects representing the actions defined in
// the configuration file.
//   - Hash: A checksum value calculated based on the merged configuration data.
//
// The data structures make use of struct tags for validation purposes, ensuring
// that required fields are present.
// The validation tags are defined using the "validate" tag in struct fields.
//
// This file provides the necessary data structures for representing and
// validating a configuration file.

var (
	mockJSONMarshal      = json.Marshal
	mockParseYaml        = parseYaml
	mockValidateConfig   = validateConfig
	mockRegisterDuration = registerDuration
	mockAdler32Hash      = adler32Hash
)

// Daemon provides a data format for daemon settings defined
// in the configuration file.
type Daemon struct {
	Interval string `validate:"duration"`
}

// Config provides a data format for the configuration file.
type Config struct {
	Daemon  Daemon           `validate:""`
	Logging system.LogConfig `validate:""`
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
	if m.Logging.JSON {
		c.Logging.JSON = m.Logging.JSON
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
	jsonData, err := mockJSONMarshal(c)
	// notest
	if err != nil {
		panic(err.Error())
	}

	// set Adler-32 hash
	hash, err := mockAdler32Hash(jsonData)
	if err != nil {
		panic(err)
	}
	c.Hash = hash
}

func adler32Hash(data []byte) (uint32, error) {
	// create a new Adler-32 hash
	hash := adler32.New()
	_, err := hash.Write(data)
	return hash.Sum32(), err
}

// LoadConfigFile loads a configuration file, validates it, and returns
// the resulting Config.
func LoadConfigFile(file string) Config {
	// read configuration file
	configContent, err := os.ReadFile(file)
	// notest
	if err != nil {
		system.FatalError("IOError", err.Error())
		return Config{}
	}

	// parse configuration file
	config, err := mockParseYaml(configContent)
	// notest
	if err != nil {
		system.FatalError("ParseError", err.Error())
		return Config{}
	}

	// validate configuration file
	validate := mockValidateConfig(config)
	// notest
	if validate != nil {
		system.FatalError("ValidationError", validate.Error())
		return Config{}
	}

	return config
}

// parseYaml parses the provided YAML content into a Config struct
// and returns it. If an error occurs during unmarshaling, it is
// also returned.
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
// It checks if the duration string is valid by attempting to parse it using
// time.ParseDuration().
func (*DurationValidator) Validate(fl validator.FieldLevel) bool {
	durationStr := fl.Field().String()
	// field is not required
	if durationStr == "" {
		return true
	}
	_, err := time.ParseDuration(durationStr)
	return err == nil
}

// validateConfig validates the provided Config object using a validator
// and returns any validation errors encountered.
// If the configuration is valid, it returns nil.
func validateConfig(config Config) error {
	// register duration validator
	validate, err := mockRegisterDuration()
	if err != nil {
		panic(err)
	}

	return validate.Struct(config)
}

// registerDuration registers a custom validation function "duration" with
// the validator and returns the validator instance and an error, if any.
func registerDuration() (*validator.Validate, error) {
	// Create a new instance of DurationValidator.
	v := newDurationValidator()

	// Create a validator instance.
	validate := v.validator

	// Register the custom validation function "duration" with
	// the validator.
	return validate, validate.RegisterValidation("duration", v.Validate)
}

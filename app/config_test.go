package app

import (
	"reflect"
	"testing"

	"github.com/piotr-ku/yaml-runner-go/system"
)

func TestParseYamlWithValidData(t *testing.T) {
	// given: We define the input, which is the contents of a valid YAML file.
	input := []byte(`
        facts:
        - name: fact1
          command: cmd1
        actions:
        - name: action1
          rules:
          - rule1
          command: cmd1
    `)

	// when: We call the parseYaml function with the input to get the result.
	result, err := parseYaml(input)

	// then: We check that the function did not return an error.
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// We check that the result object is equal to the expected data structure.
	expected := Config{
		Facts: []Fact{
			{
				Name:    "fact1",
				Command: "cmd1",
			},
		},
		Actions: []Action{
			{
				Rules:   []string{"rule1"},
				Command: "cmd1",
			},
		},
	}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("unexpected result:\n%+v\nexpected:\n%+v", result, expected)
	}
}

func TestParseYamlWithInvalidData(t *testing.T) {
	// given: We define the input, which is invalid YAML content.
	input := []byte("invalid YAML content")

	// when: We call the parseYaml function with the invalid input to get the result.
	result, err := parseYaml(input)

	// then: We check that the function returned an error.
	if err == nil {
		t.Error("expected an error, but got none")
	}

	// We check that the result object is zero-value, since the input was invalid.
	expected := Config{}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("unexpected result:\n%+v\nexpected:\n%+v", result, expected)
	}
}

func TestParseYamlWithEmptyContent(t *testing.T) {
	// given: We define the input, which is an empty byte slice.
	input := []byte{}

	// when: We call the parseYaml function with the empty input to get the result.
	result, err := parseYaml(input)

	// then: We check that the function did not return an error.
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// We check that the result object is zero-value, since there was no input.
	expected := Config{}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("unexpected result:\n%+v\nexpected:\n%+v", result, expected)
	}
}

func TestParseYamlWithMissingData(t *testing.T) {
	// given: We define the input, which is a YAML file with missing data.
	input := []byte(`
        facts:
        - name: fact1
        actions:
        - name: action1
          command: cmd1
    `)

	// when: We call the parseYaml function with the input to get the result.
	result, err := parseYaml(input)

	// then: We check that the function did not return an error.
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// We check that the result object has zero-value fields where data is missing.
	expected := Config{
		Facts: []Fact{
			{Name: "fact1", Command: ""},
		},
		Actions: []Action{
			{Rules: nil, Command: "cmd1"},
		},
	}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("unexpected result:\n%+v\nexpected:\n%+v", result, expected)
	}
}

func TestParseYamlWithMalformedData(t *testing.T) {
	// given: We define the input, which is a YAML file with a syntax error.
	input := []byte(`
        facts:
        - name: fact1
          command: cmd1
        actions:
        - name: action1
          command cmd1 # missing colon here
    `)

	// when: We call the parseYaml function with the input to get the result.
	_, err := parseYaml(input)

	// then: We check that the function returns an error.
	if err == nil {
		t.Error("expected error, but got none")
	}
}

func TestParseYamlWithEmptyInput(t *testing.T) {
	// given: We define the input, which is an empty YAML file.
	input := []byte("")

	// when: We call the parseYaml function with the empty input to get the result.
	config, err := parseYaml(input)

	// then: We check that the function returns an empty config and no error.
	if err != nil {
		t.Errorf("expected no error, but got %v", err)
	}
	if len(config.Facts) != 0 {
		t.Errorf("expected empty facts, but got %+v", config.Facts)
	}
	if len(config.Actions) != 0 {
		t.Errorf("expected empty actions, but got %+v", config.Actions)
	}
}

func TestParseYamlWithInvalidInput(t *testing.T) {
	// given: We define the input, which is an invalid YAML file.
	input := []byte("invalid_yaml_file")

	// when: We call the parseYaml function with the invalid input to get the result.
	_, err := parseYaml(input)

	// then: We check that the function returns an error.
	if err == nil {
		t.Error("expected error, but got none")
	}
}

func TestParseYamlWithValidInput(t *testing.T) {
	// given: We define the input, which is a valid YAML file.
	input := []byte(`
        facts:
          - name: fact1
            command: cmd1
            shell: /bin/bash
        actions:
          - command: cmd1
            rules: []
            shell: /bin/bash
    `)

	// when: We call the parseYaml function with the valid input to get the result.
	config, err := parseYaml(input)

	// then: We check that the function returns the expected config and no error.
	if err != nil {
		t.Errorf("expected no error, but got %v", err)
	}
	expectedConfig := Config{
		Facts:   []Fact{{Name: "fact1", Command: "cmd1", Shell: "/bin/bash"}},
		Actions: []Action{{Rules: []string{}, Command: "cmd1", Shell: "/bin/bash"}},
	}
	if !reflect.DeepEqual(config, expectedConfig) {
		t.Errorf("unexpected result:\n%+v\nexpected:\n%+v", config, expectedConfig)
	}
}

func TestParseYamlWithValidInputAndExtraFields(t *testing.T) {
	// given: We define the input, which is a valid YAML file with extra fields.
	input := []byte(`
        facts:
          - name: fact1
            command: cmd1
        actions:
          - name: action1
            rules: []
            command: cmd1
        extra_field: ignored
    `)

	// when: We call the parseYaml function with the valid input to get the result.
	config, err := parseYaml(input)

	// then: We check that the function returns the expected config and no error.
	if err != nil {
		t.Errorf("expected no error, but got %v", err)
	}
	expectedConfig := Config{
		Facts:   []Fact{{Name: "fact1", Command: "cmd1"}},
		Actions: []Action{{Rules: []string{}, Command: "cmd1"}},
	}
	if !reflect.DeepEqual(config, expectedConfig) {
		t.Errorf("unexpected result:\n%+v\nexpected:\n%+v", config, expectedConfig)
	}
}

func TestParseYamlWithInvalidYamlInput(t *testing.T) {
	// given: We define an invalid YAML input that is missing a colon after "name".
	invalidInput := []byte(`
        facts:
          - name fact1
            command: cmd1
        actions:
          - name: action1
            rules: []
            command: cmd1
    `)

	// when: We call the parseYaml function with the invalid input to get the result.
	_, err := parseYaml(invalidInput)

	// then: We check that the function returns an error.
	if err == nil {
		t.Error("expected an error, but got none")
	}
}

func TestValidateConfigWithValidData(t *testing.T) {
	// given: We define the input, which is the contents of a valid YAML file.
	input := []byte(`
        facts:
        - name: fact1
          command: cmd1
        actions:
        - name: action1
          rules:
          - rule1
          command: cmd1
    `)

	// when: We call the parseYaml function with the input to get the result.
	config, err := parseYaml(input)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	validated := validateConfig(config)

	// then: We check that the function did not return an error.
	if validated != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidateConfigWithMissingActions(t *testing.T) {
	// given: We define the input, which is the contents of a invalid YAML file.
	input := []byte(`
        facts:
        - name: fact1
          command: cmd1
    `)

	// when: We call the parseYaml function with the input to get the result.
	config, err := parseYaml(input)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	validated := validateConfig(config)

	// then: We check that the function did not return an error.
	if validated == nil {
		t.Error("expected an error, but got none")
	}
}

func TestValidateConfigWithMissingFactName(t *testing.T) {
	// given: We define the input, which is the contents of a invalid YAML file.
	input := []byte(`
        facts:
        - command: cmd1
        actions:
        - name: action1
          rules:
          - rule1
          command: cmd1
    `)

	// when: We call the parseYaml function with the input to get the result.
	config, err := parseYaml(input)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	validated := validateConfig(config)

	// then: We check that the function did not return an error.
	if validated == nil {
		t.Error("expected an error, but got none")
	}
}

func TestValidateConfigWithMissingFactCommand(t *testing.T) {
	// given: We define the input, which is the contents of a invalid YAML file.
	input := []byte(`
        facts:
        - name: fact1
        actions:
        - name: action1
          rules:
          - rule1
          command: cmd1
    `)

	// when: We call the parseYaml function with the input to get the result.
	config, err := parseYaml(input)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	validated := validateConfig(config)

	// then: We check that the function did not return an error.
	if validated == nil {
		t.Error("expected an error, but got none")
	}
}

func TestValidateConfigWithMissingActionCommand(t *testing.T) {
	// given: We define the input, which is the contents of a invalid YAML file.
	input := []byte(`
        facts:
        - name: fact1
          command: cmd1
        actions:
        - name: action1
          rules:
          - rule1
          command:
    `)

	// when: We call the parseYaml function with the input to get the result.
	config, err := parseYaml(input)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	validated := validateConfig(config)

	// then: We check that the function did not return an error.
	if validated == nil {
		t.Error("expected an error, but got none")
	}
}

func TestLoadConfigWithoutMerging(t *testing.T) {
	// given: We define the input, which is an example config file
	// and empty struct to merge.
	file := "../config-testing.yaml"

	// when: We call the LoadConfig function with the input to get the result.
	config := LoadConfigFile(file)

	// then: We check that the function returns the expected config and no error.
	if config.Logging.File != "./yaml-runner-go.log" {
		t.Errorf("unexpected result:\n%+v\nexpected:\n%+v", config.Logging.File, "./yaml-runner-go.log")
	}
	if config.Logging.Level != "error" {
		t.Errorf("unexpected result:\n%+v\nexpected:\n%+v", config.Logging.Level, "error")
	}
	if !config.Logging.Quiet {
		t.Errorf("unexpected result:\n%+v\nexpected:\n%+v", config.Logging.Quiet, false)
	}
	if !config.Logging.Json {
		t.Errorf("unexpected result:\n%+v\nexpected:\n%+v", config.Logging.Json, false)
	}
}

func TestLoadConfigWithMerging(t *testing.T) {
	// given: We define the input, which is an example config file
	// and empty struct to merge.
	file := "../config-testing.yaml"
	merge := Config{
		Daemon: Daemon{
			Interval: "2s",
		},
		Logging: system.LogConfig{
			File:  "./yaml-runner-go-merge.log",
			Level: "warn",
			Quiet: true,
			Json:  true,
		},
		Facts: []Fact{
			{Name: "MergedFact", Command: "echo mergedFact"},
		},
		Actions: []Action{
			{Command: "echo mergedAction"},
		},
	}

	// when: We call the LoadConfig function with the input to get the result.
	config := LoadConfigFile(file)

	// when: We merge configuration values.
	config.Merge(merge)

	// then: We check that the function returns the expected config and no error.
	for _, test := range []struct {
		Expected interface{}
		Got      interface{}
	}{
		{Expected: config.Daemon.Interval, Got: merge.Daemon.Interval},
		{Expected: config.Logging.File, Got: merge.Logging.File},
		{Expected: config.Logging.Level, Got: merge.Logging.Level},
		{Expected: config.Logging.Quiet, Got: merge.Logging.Quiet},
		{Expected: config.Logging.Json, Got: merge.Logging.Json},
		{Expected: config.Facts[len(config.Facts)-1].Name, Got: merge.Facts[len(merge.Facts)-1].Name},
		{Expected: config.Facts[len(config.Facts)-1].Command, Got: merge.Facts[len(merge.Facts)-1].Command},
		{Expected: config.Facts[len(config.Facts)-1].Command, Got: merge.Facts[len(merge.Facts)-1].Command},
		{Expected: config.Actions[len(config.Facts)-1].Command, Got: merge.Actions[len(merge.Actions)-1].Command},
	} {
		if test.Expected != test.Got {
			t.Errorf("unexpected result:\n%+v\nexpected:\n%+v", test.Got, test.Expected)
		}
	}
}

func TestConfigHashing(t *testing.T) {
	// given: We define the input, which is an example config
	config := Config{
		Daemon:  Daemon{Interval: "1s"},
		Logging: system.LogConfig{File: "/tmp/testing.log"},
	}

	// when: We calculate a hash for the config file
	config.CalculateHash()
	var got uint32 = config.Hash
	var expected uint32 = 3217042962

	// then: We check if hash was calculated as expected
	if got != expected {
		t.Errorf("unexpected result:\n%+v\nexpected:\n%+v", got, expected)
	}
}

func TestDurationValidator(t *testing.T) {
	// Data represents data to validate
	type Data struct {
		Duration string `validate:"duration"`
	}

	// given: We define a slice with the input, which is a Data to validate
	for _, test := range []struct {
		Duration Data
		Expected bool
	}{
		{Duration: Data{Duration: "1m"}, Expected: true},
		{Duration: Data{Duration: "2s"}, Expected: true},
		{Duration: Data{Duration: "1m30s"}, Expected: true},
		{Duration: Data{Duration: "incorrect_format"}, Expected: false},
	} {
		// then: We check validation results
		// Create a new instance of DurationValidator.
		v := newDurationValidator()

		// Create a validator instance.
		validate := v.validator

		// Register the custom validation function "duration" with the validator.
		validate.RegisterValidation("duration", v.Validate)

		// Compare got and expected result
		got := validate.Struct(test.Duration) == nil
		if test.Expected != got {
			t.Errorf("unexpected result, duration %s: got %+v expected:%+v", test.Duration, got, test.Expected)
		}
	}
}

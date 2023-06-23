package app

import (
	"reflect"
	"testing"

	"github.com/piotr-ku/yaml-runner-go/system"
)

const unexpectedResult string = "unexpected result:\n%+v\nexpected:\n%+v"
const unexpectedError string = "unexpected error: %v"
const unexpectedNone string = "expected an error, but got none"

func TestParseYamlWithValidData(t *testing.T) {
	// given: We define the input, which is the contents of a valid YAML file.
	input := []byte(`
        facts:
        - name: fact1
          command: echo revoke-lint-pluck
        actions:
        - name: action1
          rules:
          - rule1
          command: echo rectangle-fencing-unclip
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
				Command: "echo revoke-lint-pluck",
			},
		},
		Actions: []Action{
			{
				Rules:   []string{"rule1"},
				Command: "echo rectangle-fencing-unclip",
			},
		},
	}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf(unexpectedResult, result, expected)
	}
}

func TestParseYamlWithInvalidData(t *testing.T) {
	// given: We define the input, which is invalid YAML content.
	input := []byte("invalid YAML content")

	// when: We call the parseYaml function with the invalid input to get
	// the result.
	result, err := parseYaml(input)

	// then: We check that the function returned an error.
	if err == nil {
		t.Error(unexpectedNone)
	}

	// We check that the result object is zero-value, since the input
	// was invalid.
	expected := Config{}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf(unexpectedResult, result, expected)
	}
}

func TestParseYamlWithEmptyContent(t *testing.T) {
	// given: We define the input, which is an empty byte slice.
	input := []byte{}

	// when: We call the parseYaml function with the empty input to get
	// the result.
	result, err := parseYaml(input)

	// then: We check that the function did not return an error.
	if err != nil {
		t.Errorf(unexpectedError, err)
	}

	// We check that the result object is zero-value, since there was no input.
	expected := Config{}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf(unexpectedResult, result, expected)
	}
}

func TestParseYamlWithMissingData(t *testing.T) {
	// given: We define the input, which is a YAML file with missing data.
	input := []byte(`
        facts:
        - name: fact1
        actions:
        - name: action1
          command: echo carnation-secrecy-twins
    `)

	// when: We call the parseYaml function with the input to get the result.
	result, err := parseYaml(input)

	// then: We check that the function did not return an error.
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// We check that the result object has zero-value fields where data
	// is missing.
	expected := Config{
		Facts: []Fact{
			{Name: "fact1", Command: ""},
		},
		Actions: []Action{
			{Rules: nil, Command: "echo carnation-secrecy-twins"},
		},
	}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf(unexpectedResult, result, expected)
	}
}

func TestParseYamlWithMalformedData(t *testing.T) {
	// given: We define the input, which is a YAML file with a syntax error.
	input := []byte(`
        facts:
        - name: fact1
          command: echo budding-delusion-pulse
        actions:
        - name: action1
          command echo manhole-mangy-armchair # missing colon here
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

	// when: We call the parseYaml function with the empty input to get
	// the result.
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

	// when: We call the parseYaml function with the invalid input to get
	// the result.
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
            command: echo coping-huddle-creme
            shell: /bin/bash
        actions:
          - command: echo refusing-unrented-sandal
            rules: []
            shell: /bin/bash
    `)

	// when: We call the parseYaml function with the valid input to get
	// the result.
	config, err := parseYaml(input)

	// then: We check that the function returns the expected config and
	// no error.
	if err != nil {
		t.Errorf("expected no error, but got %v", err)
	}
	expectedConfig := Config{
		Facts: []Fact{
			{
				Name:    "fact1",
				Command: "echo coping-huddle-creme",
				Shell:   "/bin/bash",
			},
		},
		Actions: []Action{
			{
				Rules:   []string{},
				Command: "echo refusing-unrented-sandal",
				Shell:   "/bin/bash",
			},
		},
	}
	if !reflect.DeepEqual(config, expectedConfig) {
		t.Errorf(unexpectedResult, config, expectedConfig)
	}
}

func TestParseYamlWithValidInputAndExtraFields(t *testing.T) {
	// given: We define the input, which is a valid YAML file with extra fields.
	input := []byte(`
        facts:
          - name: fact2
            command: echo arrange-tamale-deserving
        actions:
          - name: action1
            rules: []
            command: echo diploma-fame-equity
        extra_field: ignored
    `)

	// when: We call the parseYaml function with the valid input to get
	// the result.
	config, err := parseYaml(input)

	// then: We check that the function returns the expected config and
	// no error.
	if err != nil {
		t.Errorf("expected no error, but got %v", err)
	}
	expectedConfig := Config{
		Facts: []Fact{
			{Name: "fact2", Command: "echo arrange-tamale-deserving"},
		},
		Actions: []Action{
			{Rules: []string{}, Command: "echo diploma-fame-equity"},
		},
	}
	if !reflect.DeepEqual(config, expectedConfig) {
		t.Errorf(unexpectedResult, config, expectedConfig)
	}
}

func TestParseYamlWithInvalidYamlInput(t *testing.T) {
	// given: We define an invalid YAML input that is missing a colon
	// after "name".
	invalidInput := []byte(`
        facts:
          - name fact1
            command: echo utensil-unproven-announcer
        actions:
          - name: action1
            rules: []
            command: echo humble-copier-graveness
    `)

	// when: We call the parseYaml function with the invalid input to
	// get the result.
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
          command: echo spotted-similarly-spotless
        actions:
        - name: action1
          rules:
          - rule1
          command: echo guise-recite-consult
    `)

	// when: We call the parseYaml function with the input to get the result.
	config, err := parseYaml(input)
	if err != nil {
		t.Errorf(unexpectedError, err)
	}
	validated := validateConfig(config)

	// then: We check that the function did not return an error.
	if validated != nil {
		t.Errorf(unexpectedError, err)
	}
}

func TestValidateConfigWithMissingActions(t *testing.T) {
	// given: We define the input, which is the contents of a invalid YAML file.
	input := []byte(`
        facts:
        - name: fact1
          command: echo aftermath-muzzle-thievish
    `)

	// when: We call the parseYaml function with the input to get the result.
	config, err := parseYaml(input)
	if err != nil {
		t.Errorf(unexpectedError, err)
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
        - command: echo defuse-elitism-composite
        actions:
        - name: action1
          rules:
          - rule1
          command: echo smuggling-whacking-coach
    `)

	// when: We call the parseYaml function with the input to get the result.
	config, err := parseYaml(input)
	if err != nil {
		t.Errorf(unexpectedError, err)
	}
	validated := validateConfig(config)

	// then: We check that the function did not return an error.
	if validated == nil {
		t.Error(unexpectedNone)
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
          command: echo scholar-dumping-grimacing
    `)

	// when: We call the parseYaml function with the input to get the result.
	config, err := parseYaml(input)
	if err != nil {
		t.Errorf(unexpectedError, err)
	}
	validated := validateConfig(config)

	// then: We check that the function did not return an error.
	if validated == nil {
		t.Error(unexpectedNone)
	}
}

func TestValidateConfigWithMissingActionCommand(t *testing.T) {
	// given: We define the input, which is the contents of a invalid YAML file.
	input := []byte(`
        facts:
        - name: fact1
          command: echo headlock-chaos-alibi
        actions:
        - name: action1
          rules:
          - rule1
          command:
    `)

	// when: We call the parseYaml function with the input to get the result.
	config, err := parseYaml(input)
	if err != nil {
		t.Errorf(unexpectedError, err)
	}
	validated := validateConfig(config)

	// then: We check that the function did not return an error.
	if validated == nil {
		t.Error(unexpectedNone)
	}
}

func TestLoadConfigWithoutMerging(t *testing.T) {
	// given: We define the input, which is an example config file
	// and empty struct to merge.
	file := "../config-testing.yaml"

	// when: We call the LoadConfig function with the input to get the result.
	config := LoadConfigFile(file)

	// then: We check that the function returns the expected config and
	// no error.
	if config.Logging.File != "./yaml-runner-go.log" {
		t.Errorf(unexpectedResult,
			config.Logging.File, "./yaml-runner-go.log")
	}
	if config.Logging.Level != "error" {
		t.Errorf(unexpectedResult,
			config.Logging.Level, "error")
	}
	if !config.Logging.Quiet {
		t.Errorf(unexpectedResult,
			config.Logging.Quiet, false)
	}
	if !config.Logging.JSON {
		t.Errorf(unexpectedResult,
			config.Logging.JSON, false)
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
			JSON:  true,
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

	// then: We check that the function returns the expected config and
	// no error.
	for _, test := range []struct {
		Expected interface{}
		Got      interface{}
	}{
		{
			Expected: config.Daemon.Interval,
			Got:      merge.Daemon.Interval,
		},
		{
			Expected: config.Logging.File,
			Got:      merge.Logging.File,
		},
		{
			Expected: config.Logging.Level,
			Got:      merge.Logging.Level,
		},
		{
			Expected: config.Logging.Quiet,
			Got:      merge.Logging.Quiet,
		},
		{
			Expected: config.Logging.JSON,
			Got:      merge.Logging.JSON,
		},
		{
			Expected: config.Facts[len(config.Facts)-1].Name,
			Got:      merge.Facts[len(merge.Facts)-1].Name,
		},
		{
			Expected: config.Facts[len(config.Facts)-1].Command,
			Got:      merge.Facts[len(merge.Facts)-1].Command,
		},
		{
			Expected: config.Facts[len(config.Facts)-1].Command,
			Got:      merge.Facts[len(merge.Facts)-1].Command,
		},
		{
			Expected: config.Actions[len(config.Facts)-1].Command,
			Got:      merge.Actions[len(merge.Actions)-1].Command,
		},
	} {
		if test.Expected != test.Got {
			t.Errorf(unexpectedResult,
				test.Got, test.Expected)
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
	got := config.Hash
	var expected uint32 = 2915052978

	// then: We check if hash was calculated as expected
	if got != expected {
		t.Errorf(unexpectedResult, got, expected)
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

		// Register the custom validation function "duration" with
		// the validator.
		err := validate.RegisterValidation("duration", v.Validate)
		if err != nil {
			system.FatalError("ValidationError",
				"unable to register duration function")
		}

		// Compare got and expected result
		got := validate.Struct(test.Duration) == nil
		if test.Expected != got {
			t.Errorf("unexpected result, duration %s: got %+v expected:%+v",
				test.Duration, got, test.Expected)
		}
	}
}

package app

import (
	"encoding/json"
	"errors"
	"os"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/piotr-ku/yaml-runner-go/system"
	"github.com/stretchr/testify/assert"
)

const codeIOError = 64
const codeParseError = 65
const codeValidationError = 66
const testingConfigFile = "../config-testing.yaml"

// TestParseYamlWithValidData tests the parseYaml function with valid data.
//
// It defines the input, which is the contents of a valid YAML file.
// The function calls the parseYaml function with the input to get the result.
// It checks that the function did not return an error. The function
// also checks that the result object is equal to the expected data structure.
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
	assert.Nil(t, err)

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
	assert.Equal(t, expected, result)
}

// TestParseYamlWithInvalidData tests the parseYaml function with
// invalid YAML content.
//
// The given parameter is the input, which is invalid YAML content.
// The function calls the parseYaml function with the invalid input
// to get the result. The function checks that the parseYaml function
// returned an error. The function also checks that the result object
// is a zero-value Config since the input was invalid.
func TestParseYamlWithInvalidData(t *testing.T) {
	// given: We define the input, which is invalid YAML content.
	input := []byte("invalid YAML content")

	// when: We call the parseYaml function with the invalid input to get
	// the result.
	result, err := parseYaml(input)

	// then: We check that the function returned an error.
	assert.Error(t, err)

	// We check that the result object is zero-value, since the input
	// was invalid.
	assert.Equal(t, Config{}, result)
}

// TestParseYamlWithEmptyContent is a test function that tests the behavior
// of the parseYaml function when it is called with an empty input. It checks
// that the function does not return an error and that the result object is
// a zero-value Config struct.
func TestParseYamlWithEmptyContent(t *testing.T) {
	// given: We define the input, which is an empty byte slice.
	input := []byte{}

	// when: We call the parseYaml function with the empty input to get
	// the result.
	result, err := parseYaml(input)

	// then: We check that the function did not return an error.
	assert.Nil(t, err)

	// We check that the result object is zero-value, since there was no input.
	assert.Equal(t, Config{}, result)
}

// TestParseYamlWithMissingData is a test function that tests the behavior
// of the parseYaml function when given a YAML file with missing data.
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
	assert.Nil(t, err)

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
	assert.Equal(t, expected, result)
}

// TestParseYamlWithMalformedData is a unit test function that tests
// the behavior of the parseYaml function when it is given malformed YAML data.
//
// It defines the input as a YAML file with a syntax error and then calls
// the parseYaml function with the input to get the result. It checks
// that the function returns an error.
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
	assert.Error(t, err)
}

// TestParseYamlWithEmptyInput tests the parseYaml function when given
// an empty input.
func TestParseYamlWithEmptyInput(t *testing.T) {
	// given: We define the input, which is an empty YAML file.
	input := []byte("")

	// when: We call the parseYaml function with the empty input to get
	// the result.
	config, err := parseYaml(input)

	// then: We check that the function returns an empty config and no error.
	assert.Nil(t, err)
	assert.Equal(t, 0, len(config.Facts))
	assert.Equal(t, 0, len(config.Actions))
}

// TestParseYamlWithInvalidInput tests the parseYaml function with
// invalid input.
//
// It defines the input as an invalid YAML file and calls the parseYaml
// function with the invalid input to get the result. Then it checks
// that the function returns an error.
func TestParseYamlWithInvalidInput(t *testing.T) {
	// given: We define the input, which is an invalid YAML file.
	input := []byte("invalid_yaml_file")

	// when: We call the parseYaml function with the invalid input to get
	// the result.
	_, err := parseYaml(input)

	// then: We check that the function returns an error.
	assert.Error(t, err)
}

// TestParseYamlWithValidInput tests the parseYaml function with valid input.
//
// It defines the input, which is a valid YAML file, and calls the parseYaml
// function to get the result. The function checks that the parseYaml function
// returns the expected config and no error.
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
	assert.Nil(t, err)

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
	assert.Equal(t, expectedConfig, config)
}

// TestParseYamlWithValidInputAndExtraFields tests the parseYaml function
// with valid input and extra fields.
//
// The input is a valid YAML file with extra fields. This test case checks
// that the parseYaml function correctly handles the input and returns
// the expected configuration and no error.
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
	assert.Nil(t, err)
	expectedConfig := Config{
		Facts: []Fact{
			{Name: "fact2", Command: "echo arrange-tamale-deserving"},
		},
		Actions: []Action{
			{Rules: []string{}, Command: "echo diploma-fame-equity"},
		},
	}
	assert.Equal(t, expectedConfig, config)
}

// TestParseYamlWithInvalidYamlInput is a test function that tests the parsing
// of invalid YAML input. It defines an invalid YAML input that is missing a
// colon after "name". The function calls the parseYaml function with the
// invalid input and checks that the function returns an error.
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
	assert.Error(t, err)
}

// TestValidateConfigWithValidData is a unit test for
// the validateConfig function.
//
// It tests the validation of a configuration with valid data.
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
	assert.Nil(t, err)
	validated := validateConfig(config)

	// then: We check that the function did not return an error.
	assert.Nil(t, validated)
}

// TestValidateConfigWithMissingActions is a Go function that tests
// the behavior of the validateConfig function when there are missing
// actions in the configuration.
//
// It sets up the input by defining the contents of an invalid YAML file.
// Then, it calls the parseYaml function with the input to obtain
// the configuration and any potential error. Next, it validates
// the configuration using the validateConfig function. Finally, it checks
// that the validateConfig function does not return an error.
func TestValidateConfigWithMissingActions(t *testing.T) {
	// given: We define the input, which is the contents of a invalid YAML file.
	input := []byte(`
        facts:
        - name: fact1
          command: echo aftermath-muzzle-thievish
    `)

	// when: We call the parseYaml function with the input to get the result.
	config, err := parseYaml(input)
	assert.Nil(t, err)
	validated := validateConfig(config)

	// then: We check that the function did not return an error.
	assert.NotNil(t, validated)
}

// TestValidateConfigWithMissingFactName is a Go function that tests
// the behavior of the validateConfig function when the fact name
// is missing in the input YAML file.
//
// It defines the input as the contents of an invalid YAML file
// and calls the parseYaml function to get the result. Then, it validates
// the config and checks that the function did not return an error.
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
	assert.Nil(t, err)
	validated := validateConfig(config)

	// then: We check that the function did not return an error.
	assert.NotNil(t, validated)
}

// TestValidateConfigWithMissingFactCommand tests the validateConfig function
// when the YAML input contains a missing fact command.
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
	assert.Nil(t, err)
	validated := validateConfig(config)

	// then: We check that the function did not return an error.
	assert.NotNil(t, validated)
}

// TestValidateConfigWithMissingActionCommand tests the validateConfig
// function when there is a missing action command in the input YAML file.
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
	assert.Nil(t, err)
	validated := validateConfig(config)

	// then: We check that the function did not return an error.
	assert.NotNil(t, validated)
}

// TestLoadConfigWithoutMerging is a test function that verifies the behavior
// of the LoadConfigWithoutMerging function.
//
// It tests the function by defining the input, calling the function, and
// checking the expected output.
func TestLoadConfigWithoutMerging(t *testing.T) {
	// given: We define the input, which is an example config file
	// and empty struct to merge.
	file := testingConfigFile

	// when: We call the LoadConfig function with the input to get the result.
	config := LoadConfigFile(file)

	// then: We check that the function returns the expected config and
	// no error.
	assert.Equal(t, "./yaml-runner-go.log", config.Logging.File)
	assert.Equal(t, "error", config.Logging.Level)
	assert.Equal(t, true, config.Logging.Quiet)
	assert.Equal(t, true, config.Logging.JSON)
}

// TestLoadConfigWithMerging is a test function that verifies the behavior
// of the LoadConfigWithMerging function.
//
// The function takes an input config file and an empty struct to merge.
// It loads the config file, merges it with the provided struct, and
// then checks if the merged values match the expected values.
// The function uses the testing.T object to report any failures.
func TestLoadConfigWithMerging(t *testing.T) {
	// given: We define the input, which is an example config file
	// and empty struct to merge.
	file := testingConfigFile
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
		assert.Equal(t, test.Expected, test.Got)
	}
}

// TestLoadConfiFileIOError is a unit test function that tests
// the behavior of the LoadConfigFile function when an IO error occurs.
//
// It mocks the os.Exit function, defines an input which is
// a non-existing config file, and checks that the function will cause
// a fatal error. The function then checks that the expected config
// and no error are returned.
func TestLoadConfiFileIOError(t *testing.T) {
	// mock os.Exit
	var rc int
	system.MockOsExit = func(code int) {
		rc = code
	}
	defer func() {
		system.MockOsExit = os.Exit
	}()

	// given: We define the input, which is an non-existing config file
	file := "../non-existing-file.yaml"

	// then: We check that the function will cause a fatal error
	LoadConfigFile(file)

	// then: We check that the function returns the expected config and
	// no error.
	assert.Equal(t, codeIOError, rc)
}

// TestLoadConfiFileParseError is a test function that tests the behavior
// of LoadConfigFile when encountering a parse error in the config file.
//
// The function mocks the os.Exit function, the parseYaml function,
// and defines the input as a non-existing config file. It then checks
// that the LoadConfigFile function causes a fatal error and returns
// the expected config and no error.
func TestLoadConfiFileParseError(t *testing.T) {
	// mock os.Exit
	var rc int
	system.MockOsExit = func(code int) {
		rc = code
	}
	defer func() {
		system.MockOsExit = os.Exit
	}()

	// mock parseYaml
	mockParseYaml = func(_ []byte) (Config, error) {
		return Config{}, errors.New("fake YAML error")
	}
	defer func() {
		mockParseYaml = parseYaml
	}()

	// given: We define the input, which is an non-existing config file
	file := testingConfigFile

	// then: We check that the function will cause a fatal error
	LoadConfigFile(file)

	// then: We check that the function returns the expected config and
	// no error.
	assert.Equal(t, codeParseError, rc)
}

// TestLoadConfiFileValidationError is a test function that tests
// the behavior of the LoadConfigFile function when a validation error occurs.
//
// This function mocks the os.Exit function, the validateConfig function,
// and defines the input file. It then calls the LoadConfigFile function
// and checks that it causes a fatal error. Finally, it asserts
// that the return value of the function is the expected error code.
func TestLoadConfiFileValidationError(t *testing.T) {
	// mock os.Exit
	var rc int
	system.MockOsExit = func(code int) {
		rc = code
	}
	defer func() {
		system.MockOsExit = os.Exit
	}()

	// mock validateConfig
	mockValidateConfig = func(_ Config) error {
		return errors.New("fake validation error")
	}
	defer func() {
		mockValidateConfig = validateConfig
	}()

	// given: We define the input, which is an non-existing config file
	file := testingConfigFile

	// then: We check that the function will cause a fatal error
	LoadConfigFile(file)

	// then: We check that the function returns the expected config and
	// no error.
	assert.Equal(t, codeValidationError, rc)
}

// TestConfigHashing tests the hashing functionality of the Config struct.
//
// It creates an example config with predefined values, calculates the hash
// of the config file, and compares it to the expected value to ensure that
// the hash was calculated correctly.
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
	assert.Equal(t, expected, got)
}

// TestConfigHashingJsonMarshallError is a test function that tests
// the scenario when there is an error in the json.Marshal() function call.
// It mocks the json.Marshal() function and verifies if the CalculateHash()
// function panics.
func TestConfigHashingJsonMarshallError(t *testing.T) {
	// mock json.Marshal()
	mockJSONMarshal = func(_ any) ([]byte, error) {
		return []byte{}, errors.New("json.Marshall error")
	}
	defer func() {
		mockJSONMarshal = json.Marshal
	}()

	// given: We define the input, which is an example config
	config := Config{
		Daemon:  Daemon{Interval: "1s"},
		Logging: system.LogConfig{File: "/tmp/testing.log"},
	}

	// We check if CalculateHash function panics
	assert.Panics(t, func() { config.CalculateHash() })
}

// TestConfigHashingAdler32Error is a test function that tests
// the behavior of the CalculateHash method of the Config struct when
// the adler32Hash function returns an error.
func TestConfigHashingAdler32Error(t *testing.T) {
	// mock adler32Hash
	mockAdler32Hash = func(_ []byte) (uint32, error) {
		return 0, errors.New("fake adler32 hash")
	}
	defer func() {
		mockAdler32Hash = adler32Hash
	}()

	// given: We define the input, which is an empty config
	config := Config{}

	// We check if CalculateHash function panics
	assert.Panics(t, func() { config.CalculateHash() })
}

// TestDurationValidator is a test function that validates the Duration field
// of the Data struct.
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
		assert.Nil(t, err)

		// Compare got and expected result
		assert.Equal(t, test.Expected, validate.Struct(test.Duration) == nil)
	}
}

// TestDurationValidatorRegisterError tests the duration validator
// registration error.
//
// It mocks the registerDuration function to return a fake validator error.
// It creates a new instance of DurationValidator and a validator instance.
// It registers the custom validation function "duration" with the validator.
// It defines the input as an empty Config.
// It checks that the function will cause a fatal error using assert.Panics.
func TestDurationValidatorRegisterError(t *testing.T) {
	// mock registerDuration
	mockRegisterDuration = func() (*validator.Validate, error) {
		// Create a new instance of DurationValidator.
		v := newDurationValidator()

		// Create a validator instance.
		validate := v.validator

		// Register the custom validation function "duration" with
		// the validator.
		return validate, errors.New("fake validator error")
	}
	defer func() {
		mockRegisterDuration = registerDuration
	}()

	// given: We define the input, which is an empty Config
	config := Config{}

	// then: We check that the function will cause a fatal error
	assert.Panics(t, func() { _ = validateConfig(config) })
}

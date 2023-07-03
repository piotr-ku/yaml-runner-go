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

func TestParseYamlWithInvalidInput(t *testing.T) {
	// given: We define the input, which is an invalid YAML file.
	input := []byte("invalid_yaml_file")

	// when: We call the parseYaml function with the invalid input to get
	// the result.
	_, err := parseYaml(input)

	// then: We check that the function returns an error.
	assert.Error(t, err)
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
	mockParseYaml = func(content []byte) (Config, error) {
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
	mockValidateConfig = func(config Config) error {
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

func TestConfigHashingJsonMarshallError(t *testing.T) {
	// mock json.Marshal()
	mockJSONMarshal = func(v any) ([]byte, error) {
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

func TestConfigHashingAdler32Error(t *testing.T) {
	// mock adler32Hash
	mockAdler32Hash = func(data []byte) (uint32, error) {
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
	assert.Panics(t, func() { validateConfig(config) })
}

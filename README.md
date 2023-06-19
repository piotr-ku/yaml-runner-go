# YAML Runner Go

YAML Runner Go is an application that executes commands based on the rules defined in a YAML file. It provides the flexibility to run commands either once or as a daemon at specific intervals.

![Continous Integration status](https://github.com/piotr-ku/yaml-runner-go/actions/workflows/integration.yml/badge.svg?branch=main)

## Installation

To install YAML Runner Go from sources, including any dependencies or system requirements, please follow these clear instructions:

### Prerequisites

Before installing YAML Runner Go from sources, ensure that you have the following prerequisites:

1. Go Programming Language: Ensure that you have Go installed on your system. You can download and install Go from the official Go website: [https://golang.org/](https://golang.org/). Follow the installation instructions specific to your operating system.

### Installation Steps

Follow these steps to install YAML Runner Go from sources:

1. Open a terminal or command prompt.

2. Clone the YAML Runner Go repository from GitHub using the `git clone` command:

   ```shell
   git clone https://github.com/piotr-ku/yaml-runner-go.git
   ```

   This command will clone the YAML Runner Go repository to your local machine.

3. Change into the cloned repository directory:

   ```shell
   cd yaml-runner-go
   ```

4. Use the `go build` command to build the YAML Runner Go executable:

   ```shell
   go build .
   ```

   This command will compile the source code and generate the YAML Runner Go executable.

5. (Optional) If you want to install the YAML Runner Go executable system-wide, you can use the `go install` command:

   ```shell
   go install .
   ```

   This command will build and install the YAML Runner Go executable to the Go bin directory, allowing you to run it from anywhere in the terminal.

### Dependencies

YAML Runner Go, being a standard Go application, manages its dependencies using Go modules. When you build the application using the `go build` or `go install` command, Go automatically downloads and installs the required dependencies specified in the Go.mod and Go.sum files.

There is no need for manual installation or management of dependencies.

### System Requirements

YAML Runner Go should be compatible with most major operating systems, including Windows, macOS, and Linux. It relies on the Go programming language's cross-platform support.

Ensure that your system meets the minimum requirements for installing and running Go. Refer to the Go documentation for specific system requirements based on your operating system.

By following these installation instructions, you should be able to install YAML Runner Go from sources, including any necessary dependencies, and start using it to execute commands based on YAML configuration files.

## Usage

`yaml-runner-go [command]`

## Available Commands

* completion: Generate the autocompletion script for the specified shell
* daemon: Run actions periodically in the background
* help: Help about any command
* oneshot: Runs actions once and exits

## Flags

* --config string: Specifies the configuration file in YAML format (default: "./config.yaml")
* --debug: Enables debug logging
* --help, -h: Provides help for yaml-runner-go
* --interval string: Sets the interval for the daemon
* --json: Enables JSON formatting for the output
* --log string: Enables logging to a file
* --quiet: Enables quiet mode

To get more information about a specific command, use the following syntax:

`yaml-runner-go [command] --help`

## Configuration

The YAML configuration file allows you to define the rules, actions, and command execution settings for YAML Runner Go. Here's the structure and syntax of the configuration file:

```yaml
daemon:
  interval: 5s
logging:
  file: ./yaml-runner-go.log
  quiet: false
  level: debug
  json: false
facts:
  - name: apacheIsRunning
    command: "curl --connect-timeout 1 -s http://localhost:80/; echo $?;"
  - name: loadAverage1
    command: "[[ -e /proc/loadavg ]] && awk '{print $(NF-2)}' /proc/loadavg | cut -d. -f1 || sysctl -n vm.loadavg | awk '{print $2}' | cut -d, -f1"
    shell: /bin/bash
actions:
  - command: "echo \"Stopping apache\""
    rules:
      - "[[ ${loadAverage1} -gt 15 ]]"
      - "[[ ${apacheIsRunning} -eq 0 ]]"
    shell: /bin/bash
  - command: "echo \"Starting apache\""
    rules:
      - "[[ ${loadAverage1} -lt 15 ]]"
      - "[[ ${apacheIsRunning} -ne 0 ]]"
    shell: /bin/bash
```

### Structure

The configuration file consists of the following sections:

- **daemon**: Defines the settings for the YAML Runner Go daemon, including the interval at which the actions should be executed. The interval value should be specified in a valid duration format (e.g., "5s" for 5 seconds).

- **logging**: Specifies the logging settings for the application. It includes the log file path, whether to enable quiet mode (suppressing non-error log messages), the log level (e.g., "debug", "info", "warn"), and whether to format log output in JSON.

- **facts**: Describes the facts or variables that can be used in the rules section. Each fact has a unique name and a command associated with it. The command will be executed to obtain the value of the fact.

- **actions**: Defines the actions to be executed based on the specified rules. Each action consists of a command to be executed when the rules evaluate to true. The rules are expressed using boolean expressions that can reference the facts defined earlier.

### Syntax

- **Key-Value Pairs**: The configuration file is structured using key-value pairs. Each key is followed by a colon, and the associated value is indented below it.

- **Lists**: Lists are represented using a hyphen followed by a space ("- "). In the configuration file, the `facts` and `actions` sections are represented as lists.

- **Command Execution**: Command execution is defined by providing the command to be executed as a string value under the `command` key.

- **Rules**: Rules are expressed as boolean expressions using square brackets. The expressions can reference the facts by using the `${factName}` syntax.

Ensure that the configuration file follows the proper YAML syntax and indentation rules for accurate parsing by YAML Runner Go.

Use this configuration file as a template and modify it according to your specific requirements.

## Use Cases

YAML Runner Go can be useful in various scenarios where you need to automate command execution based on a YAML file configuration. Here are some possible use cases:

1. **Scheduled Tasks**: By running YAML Runner Go as a daemon, you can schedule and automate recurring tasks. For example, you can define commands to be executed at specific intervals, such as fetching data from APIs, running backups, or performing system maintenance.

2. **Continuous Integration/Deployment**: YAML Runner Go can be integrated into your CI/CD pipelines to execute commands defined in the YAML configuration. This allows you to automate build processes, run tests, deploy applications, or trigger other actions based on specific events or conditions.

3. **System Administration**: You can use YAML Runner Go to streamline and automate various system administration tasks. For instance, you can define commands to monitor system resources, check for software updates, manage user accounts, or perform routine maintenance tasks on multiple servers.

4. **Data Processing**: If you have data processing pipelines that involve executing specific commands or scripts, YAML Runner Go can help simplify the execution and management of these pipelines. You can define commands to process data files, perform transformations, extract information, or generate reports based on your YAML configuration.

5. **Monitoring and Alerting**: YAML Runner Go can be utilized to periodically execute commands that monitor system metrics, log files, or external services. By defining actions in the YAML file to check for specific conditions or trigger alerts, you can proactively monitor your infrastructure and receive notifications when anomalies or critical events occur.

These are just a few examples of how YAML Runner Go can be applied in different scenarios. The flexibility of defining commands in a YAML file and the ability to run them either as one-shot tasks or as a recurring daemon make it a versatile tool for automating command execution in various domains.

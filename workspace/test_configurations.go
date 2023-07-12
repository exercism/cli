package workspace

import (
	"runtime"
	"strings"
)

type TestConfiguration struct {
	// The static portion of the test Command, which will be run for every test on this track. Examples include `cargo test` or `go test`.
	// Might be empty if there are platform-specific versions
	Command string

	// Windows-specific test command. Mostly relevant for tests wrapped by shell invocations. Falls back to `Command` if we're not running windows or this is empty.
	WindowsCommand string
}

func (c *TestConfiguration) GetTestCommand() (string, error) {
	var cmd string
	if runtime.GOOS == "windows" && c.WindowsCommand != "" {
		cmd = c.WindowsCommand
	} else {
		cmd = c.Command
	}

	if strings.Contains(cmd, "{{test_files}}") {
		testFiles, err := getTestFiles()
		if err != nil {
			return "", err
		}
		cmd = strings.ReplaceAll(cmd, "{{test_files}}", strings.Join(testFiles, " "))
	}

	return cmd, nil
}

var TestConfigurations = map[string]TestConfiguration{
	"8th": {
		Command:        "tester.sh",
		WindowsCommand: "tester.bat",
	},
	"ballerina": {
		Command: "bal test",
	},
	"c": {
		Command: "make",
	},
	"cfml": {
		Command: "box task run TestRunner",
	},
	"cobol": {
		Command:        "test.sh",
		WindowsCommand: "test.ps1",
	},
	"coffeescript": {
		Command: "jasmine-node --coffee {{test_files}}",
	},
	"crystal": {
		Command: "crystal spec",
	},
	"csharp": {
		Command: "dotnet test",
	},
	"dart": {
		Command: "dart test",
	},
	"elixir": {
		Command: "mix test",
	},
	"elm": {
		Command: "elm-test",
	},
	"go": {
		Command: "go test",
	},
	"rust": {
		Command: "cargo test --",
	},
	"ruby": {
		Command: "ruby {{test_files}}",
	},
}

func getTestFiles() ([]string, error) {
	testFiles, err := NewExerciseConfig(".")
	if err != nil {
		return []string{}, err
	}
	return testFiles.Files.Test, nil
}

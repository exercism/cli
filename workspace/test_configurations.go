package workspace

import "runtime"

type TestConfiguration struct {
	// The static portion of the test Command, which will be run for every test on this track. Examples include `cargo test` or `go test`.
	// Might be empty if there are platform-specific versions
	Command string

	// Windows-specific test command. Mostly relevant for tests wrapped by shell invocations. Falls back to `Command` if we're not running windows or this is empty.
	WindowsCommand string

	// Some tracks test by running a specific file, such as `ruby lasagna_test.rb`. Set this to `true` to look up and include the name of the default test file(s).
	AppendTestFiles bool
}

func (c *TestConfiguration) GetTestCommand() string {
	if runtime.GOOS == "windows" && c.WindowsCommand != "" {
		return c.WindowsCommand
	}
	return c.Command
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
		Command:         "jasmine-node --coffee",
		AppendTestFiles: true,
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
		Command:         "ruby",
		AppendTestFiles: true,
	},
}

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
		exerciseConfig, err := NewExerciseConfig(".")
		if err != nil {
			return "", err
		}

		testFiles, err := exerciseConfig.GetTestFiles()
		if err != nil {
			return "", err
		}
		cmd = strings.ReplaceAll(cmd, "{{test_files}}", strings.Join(testFiles, " "))
	}

	return cmd, nil
}

var TestConfigurations = map[string]TestConfiguration{
	"8th": {
		Command:        "bash tester.sh",
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
		Command:        "bash test.sh",
		WindowsCommand: "pwsh test.ps1",
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
	"emacs-lisp": {
		Command: "emacs -batch -l ert -l *-test.el -f ert-run-tests-batch-and-exit",
	},
	"erlang": {
		Command: "rebar3 eunit",
	},
	"fsharp": {
		Command: "dotnet test",
	},
	"gleam": {
		Command: "gleam test",
	},
	"go": {
		Command: "go test",
	},
	"groovy": {
		Command: "gradle test",
	},
	"haskell": {
		Command: "stack test",
	},
	"java": {
		Command: "gradle test",
	},
	"javascript": {
		Command: "npm run test",
	},
	"julia": {
		Command: "julia runtests.jl",
	},
	"kotlin": {
		Command:        "./gradlew test",
		WindowsCommand: "gradlew.bat test",
	},
	"lfe": {
		Command: "make test",
	},
	"lua": {
		Command: "busted",
	},
	"mips": {
		Command: "java -jar /path/to/mars.jar nc runner.mips impl.mips",
	},
	"nim": {
		Command: "nim r {{test_files}}",
	},
	"ocaml": {
		Command: "make",
	},
	"perl5": {
		Command: "prove .",
	},
	"php": {
		Command: "phpunit {{test_files}}",
	},
	"purescript": {
		Command: "spago test",
	},
	"python": {
		Command: "python3 -m pytest -o markers=task {{test_files}}",
	},
	"racket": {
		Command: "raco test {{test_files}}",
	},
	"raku": {
		Command: "prove6 {{test_files}}",
	},
	"reasonml": {
		Command: "npm run test",
	},
	"red": {
		Command: "red {{test_files}}",
	},
	"ruby": {
		Command: "ruby {{test_files}}",
	},
	"rust": {
		Command: "cargo test --",
	},
	"scala": {
		Command: "sbt test",
	},
	"sml": {
		Command: "poly -q --use {{test_files}}",
	},
	"swift": {
		Command: "swift test",
	},
	"tcl": {
		Command: "tclsh {{test_files}}",
	},
	"typescript": {
		Command: "yarn test",
	},
	"vbnet": {
		Command: "dotnet test",
	},
	"vlang": {
		Command: "v -stats test run_test.v",
	},
	"wasm": {
		Command: "npm run test",
	},
	"wren": {
		Command: "wrenc {{test_files}}",
	},
	"x86-64-assembly": {
		Command: "make",
	},
	"zig": {
		Command: "zig test {{test_files}}",
	},
}

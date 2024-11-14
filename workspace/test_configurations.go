package workspace

import (
	"fmt"
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

	// pre-declare these so we can conditionally initialize them
	var exerciseConfig *ExerciseConfig
	var err error

	if strings.Contains(cmd, "{{") {
		// only read exercise's config.json if we need it
		exerciseConfig, err = NewExerciseConfig(".")
		if err != nil {
			return "", err
		}
	}

	if strings.Contains(cmd, "{{solution_files}}") {
		if exerciseConfig == nil {
			return "", fmt.Errorf("exerciseConfig not initialize before use")
		}
		solutionFiles, err := exerciseConfig.GetSolutionFiles()
		if err != nil {
			return "", err
		}
		cmd = strings.ReplaceAll(cmd, "{{solution_files}}", strings.Join(solutionFiles, " "))
	}
	if strings.Contains(cmd, "{{test_files}}") {
		if exerciseConfig == nil {
			return "", fmt.Errorf("exerciseConfig not initialize before use")
		}
		testFiles, err := exerciseConfig.GetTestFiles()
		if err != nil {
			return "", err
		}
		cmd = strings.ReplaceAll(cmd, "{{test_files}}", strings.Join(testFiles, " "))
	}

	return cmd, nil
}

// some tracks aren't (or won't be) implemented; every track is listed either way
var TestConfigurations = map[string]TestConfiguration{
	"8th": {
		Command: "8th -f test.8th",
	},
	// abap: tests are run via "ABAP Development Tools", not the CLI
	"arm64-assembly": {
		Command: "make",
	},
	"arturo": {
		Command: "arturo tester.art",
	},
	"awk": {
		Command: "bats {{test_files}}",
	},
	"ballerina": {
		Command: "bal test",
	},
	"batch": {
		WindowsCommand: "cmd /c {{test_files}}",
	},
	"bash": {
		Command: "bats {{test_files}}",
	},
	"c": {
		Command: "make",
	},
	"cairo": {
		Command: "scarb cairo-test",
	},
	"cfml": {
		Command: "box task run TestRunner",
	},
	"clojure": {
		// chosen because the docs recommend `clj` by default and `lein` as optional
		Command: "clj -X:test",
	},
	"cobol": {
		Command:        "bash test.sh",
		WindowsCommand: "pwsh test.ps1",
	},
	"coffeescript": {
		Command: "jasmine-node --coffee {{test_files}}",
	},
	// common-lisp: tests are loaded into a "running Lisp implementation", not the CLI directly
	"cpp": {
		Command: "make",
	},
	"crystal": {
		Command: "crystal spec",
	},
	"csharp": {
		Command: "dotnet test",
	},
	"d": {
		// this always works even if the user installed DUB
		Command: "dmd source/*.d -de -w -main -unittest",
	},
	"dart": {
		Command: "dart test",
	},
	// delphi: tests are run via IDE
	"elixir": {
		Command: "mix test",
	},
	"elm": {
		Command: "elm-test",
	},
	"emacs-lisp": {
		Command: "emacs -batch -l ert -l {{test_files}} -f ert-run-tests-batch-and-exit",
	},
	"erlang": {
		Command: "rebar3 eunit",
	},
	"fortran": {
		Command: "make",
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
	"idris": {
		Command: "pack test `basename *.ipkg .ipkg`",
	},
	"j": {
		Command: `jconsole -js "exit echo unittest {{test_files}} [ load {{solution_files}}"`,
	},
	"java": {
		Command:        "./gradlew test",
		WindowsCommand: "gradlew.bat test",
	},
	"javascript": {
		Command: "npm run test",
	},
	"jq": {
		Command: "bats {{test_files}}",
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
	// objective-c: tests are run via XCode. There's a CLI option (ruby gem `objc`), but the docs note that this is an inferior experience
	"ocaml": {
		Command: "make",
	},
	"perl5": {
		Command: "prove .",
	},
	// pharo-smalltalk: tests are run via IDE
	"php": {
		Command: "phpunit {{test_files}}",
	},
	// plsql: test are run via a "mounted oracle db"
	"powershell": {
		Command: "Invoke-Pester",
	},
	"prolog": {
		Command: "swipl -f {{solution_files}} -s {{test_files}} -g run_tests,halt -t 'halt(1)'",
	},
	"purescript": {
		Command: "spago test",
	},
	"pyret": {
		Command: "pyret {{test_files}}",
	},
	"python": {
		Command: "python3 -m pytest -o markers=task {{test_files}}",
	},
	"r": {
		Command: "Rscript {{test_files}}",
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
	"roc": {
		Command: "roc test {{test_files}}",
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
	// scheme: docs present 2 equally valid test methods (`make chez` and `make guile`). So I wasn't sure which to pick
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
	"uiua": {
		Command: "uiua test {{test_files}}",
	},
	// unison: tests are run from an active UCM session
	"vbnet": {
		Command: "dotnet test",
	},
	// vimscript: tests are run from inside a vim session
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
	"yamlscript": {
		Command: "make test",
	},
	"zig": {
		Command: "zig test {{test_files}}",
	},
}

package config

import (
	"fmt"
	"path/filepath"
)

type GlobRule struct {
	Matches []string
	Except  []string
}

// The configuration rules below originaly is a port from
// Source: @petertseng - https://gist.github.com/petertseng/e3e88bf1c383865ff67f4095413993b2
//
var defaultTrackGlobs = map[string]GlobRule{
	"ballerina":  {Matches: []string{"*.bal"}},
	"bash":       {Matches: []string{"*.sh"}, Except: []string{"*_test.sh"}},
	"c":          {Matches: []string{"src/*.c", "src/*.h"}},
	"ceylon":     {Matches: []string{"source/*/*.ceylon"}, Except: []string{"module.ceylon", " Test.ceylon"}},
	"cfml":       {Matches: []string{"*.cfc"}, Except: []string{"Test.cfc", "TestRunner.cfc"}},
	"clojure":    {Matches: []string{"src"}, Except: []string{"*.clj"}},
	"coq":        {Matches: []string{"*.v"}, Except: []string{"test.v"}},
	"cpp":        {Matches: []string{"*.cpp", "*.h"}, Except: []string{"_test.cpp"}},
	"crystal":    {Matches: []string{"src/*.cr"}},
	"csharp":     {Matches: []string{"*.cs"}, Except: []string{"Test.cs"}},
	"d":          {Matches: []string{"source/*.d"}},
	"dart":       {Matches: []string{"lib/*.dart"}},
	"ecmascript": {Matches: []string{"*.js"}, Except: []string{".spec.js", "gulpfile.js"}},
	"elixir":     {Matches: []string{"lib/*.ex", "lib/*.exs"}, Except: []string{"lib/*_test.exs"}},
	"elm":        {Matches: []string{"*.elm"}, Except: []string{"Tests.elm"}},
	"erlang":     {Matches: []string{"src/*.erl"}},
	"factor":     {Matches: []string{"*.factor"}, Except: []string{"-tests.factor"}},
	"fsharp":     {Matches: []string{"*.fs"}, Except: []string{"Test.fs", "Program.fs"}},
	"go":         {Matches: []string{"*.go"}, Except: []string{"*_test.go"}},
	"groovy":     {Matches: []string{"*.groovy"}, Except: []string{"Spec.groovy"}},
	"haskell":    {Matches: []string{"src/*.hs"}},
	"haxe":       {Matches: []string{"src/*.hx"}},
	"idris":      {Matches: []string{"src/*.idr"}},
	"java":       {Matches: []string{"src/main/java/*.java"}},
	// Note: I include .ts files, typically students won"t have.,
	"javascript": {Matches: []string{"*.ts", "*.js"}, Except: []string{".spec.js"}},
	"julia":      {Matches: []string{"*.jl"}, Except: []string{"runtests.jl"}},
	"kotlin":     {Matches: []string{"src/main/kotlin/*kt"}},
	"lua":        {Matches: []string{"*.lua"}, Except: []string{"_spec.lua"}},
	"nim":        {Matches: []string{"*.nim"}, Except: []string{"_test.nim"}},
	"ocaml":      {Matches: []string{"*.ml"}, Except: []string{"test.ml"}},
	// Note: tests have .t ext, files have *.pm extension.,
	// So, "files of same extension as test file" wouldn"t work
	// This script doesn"t care, but an autodetector will fail.
	"perl5": {Matches: []string{"*.pm"}},
	// Same idea w/ "files of same extension as test file"
	"perl6":      {Matches: []string{"*.pm6"}},
	"php":        {Matches: []string{"*.php"}, Except: []string{"_test.php"}},
	"pony":       {Matches: []string{"*.pony"}, Except: []string{"test.pony"}},
	"prolog":     {Matches: []string{"*.pl"}},
	"purescript": {Matches: []string{"src/*.purs"}},
	"python":     {Matches: []string{"*.py"}, Except: []string{"_test.py"}},

	"r":        {Matches: []string{"*.R"}, Except: []string{`/test_\S+\.R/`}},
	"racket":   {Matches: []string{"*.rkt"}, Except: []string{"-test.rkt"}},
	"reasonml": {Matches: []string{"src/*.re"}},
	"ruby":     {Matches: []string{"*.rb"}, Except: []string{"_test.rb"}},
	"rust":     {Matches: []string{"src/*.rs"}},
	"scala":    {Matches: []string{"src/main"}, Except: []string{"scala/*.scala"}},
	"sml":      {Matches: []string{"*.sml"}, Except: []string{"test.sml", "testlib.sml"}},
	"swift":    {Matches: []string{"Sources/*.swift", "Sources/?*/*.swift"}},
	// At some point (https://github.com/exercism/swift/pull/363) everything got moved into directories:
	"typescript": {Matches: []string{"*.ts"}, Except: []string{".test.ts"}},
	"vimscript":  {Matches: []string{"*.vim"}},
}

// FindSolutions Infer exercism solution given a track
func FindSolutions(trackGlobs map[string]GlobRule, track string, exerciseDir string) ([]string, error) {
	var result []string
	exerciseRule, ok := trackGlobs[track]
	if !ok {
		return nil, fmt.Errorf("Cannot find file patterns for track %s", track)
	}

	for _, rule := range exerciseRule.Matches {
		paths, err := filepath.Glob(rule)
		if err != nil {
			return nil, err
		}

		if len(exerciseRule.Except) == 0 {
			result = append(result, paths...)
			continue
		}

		for _, path := range paths {
			for _, exceptRule := range exerciseRule.Except {
				match, err := filepath.Match(exceptRule, path)
				if err != nil {
					return nil, err
				}

				if !match {
					result = append(result, path)
				}
			}

		}
	}

	for i, path := range result {
		result[i] = filepath.Join(exerciseDir, path)
	}

	return result, nil
}

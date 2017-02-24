package cmd

import "testing"

func TestIsTest(t *testing.T) {
	testCases := []struct {
		name   string
		isTest bool
	}{
		{
			name:   "exercise/whatever-test.ext",
			isTest: true,
		},
		{
			name:   "exercise/whatever.ext",
			isTest: false,
		},
		{
			name:   "exercise/whatever_test.spec.ext",
			isTest: true,
		},
		{
			name:   "exercise/WhateverTest.ext",
			isTest: true,
		},
		{
			name:   "exercise/Whatever.ext",
			isTest: false,
		},
		{
			name:   "exercise/whatever_test.ext",
			isTest: true,
		},
		{
			name:   "exercise/whatever.ext",
			isTest: false,
		},
		{
			name:   "exercise/test.ext",
			isTest: true,
		},
		{
			name:   "exercise/Whatever.t", // perl
			isTest: true,
		},
		{
			name:   "whatever_spec.ext", // lua
			isTest: true,
		},
	}

	for _, tt := range testCases {
		if isTest(tt.name) != tt.isTest {
			t.Fatalf("Expected isTest(%s) to be %t", tt.name, tt.isTest)
		}
	}
}

func TestIsREADME(t *testing.T) {
	testCases := []struct {
		name     string
		isREADME bool
	}{
		{
			name:     "exercise/README.md",
			isREADME: true,
		},
		{
			name:     "exercise/README",
			isREADME: true,
		},
		{
			name:     "exercise/README.txt",
			isREADME: true,
		},
		{
			name:     "exercise/some_exercise.py",
			isREADME: false,
		},
		{
			name:     "exercise/readme.lua",
			isREADME: false,
		},
		{
			name:     "exercise/readme_spec.lua",
			isREADME: false,
		},
	}

	for _, tt := range testCases {
		if isREADME(tt.name) != tt.isREADME {
			t.Fatalf("Expected isREADME(%s) to be %t", tt.name, tt.isREADME)
		}
	}
}

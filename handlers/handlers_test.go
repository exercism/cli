package handlers

import "testing"

func TestIsTest(t *testing.T) {
	testCases := []struct {
		name   string
		isTest bool
	}{
		{
			name:   "problem/whatever-test.ext",
			isTest: true,
		},
		{
			name:   "problem/whatever.ext",
			isTest: false,
		},
		{
			name:   "problem/whatever_test.spec.ext",
			isTest: true,
		},
		{
			name:   "problem/WhateverTest.ext",
			isTest: true,
		},
		{
			name:   "problem/Whatever.ext",
			isTest: false,
		},
		{
			name:   "problem/whatever_test.ext",
			isTest: true,
		},
		{
			name:   "problem/whatever.ext",
			isTest: false,
		},
		{
			name:   "problem/test.ext",
			isTest: true,
		},
		{
			name:   "problem/Whatever.t", // perl
			isTest: true,
		},
	}

	for _, tt := range testCases {
		if isTest(tt.name) != tt.isTest {
			t.Fatalf("Expected isTest(%s) to be %t", tt.name, tt.isTest)
		}
	}
}

package bci

import (
	"bytes"
	"strings"
	"testing"
)

type TestCase struct {
	name   string
	input  []string
	output []string
}

// ---------------------------------------------------------------------------

func TestInterpreter(t *testing.T) {
	testCases := []TestCase{
		{
			name: "ASCII code to character",
			input: []string{
				"pounce on 65",
				"yowl",
			},
			output: []string{
				"A",
			},
		},
	}

	executeTests(t, testCases)
}

// ---------------------------------------------------------------------------

func executeTests(t *testing.T, testCases []TestCase) {
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			executeTestCase(t, tc)
		})
	}
}

func executeTestCase(t *testing.T, tc TestCase) {
	var buffer bytes.Buffer

	interpreter := NewInterpreterWithIO(tc.input, &buffer)

	err := interpreter.Run()
	if err != nil {
		t.Errorf("[%s] Unexpected error: %s", tc.name, err)
		return
	}

	var actualOutput []string
	if len(buffer.String()) > 0 {
		actualOutput = strings.Split(buffer.String(), "\n")
	}

	if len(tc.output) != len(actualOutput) {
		t.Errorf("[%s] Number of lines differ. Expected %d, got %d", tc.name, len(tc.output), len(actualOutput))
	}

	for i, expectedText := range tc.output {
		expectedText = strings.TrimSpace(expectedText)
		actualText := actualOutput[i]
		actualText = strings.TrimSpace(actualText)
		if actualText != expectedText {
			t.Errorf("[%s] Line %d: Lines do not match. Expected '%s', got '%s'", tc.name, i, expectedText, actualText)
		}
	}
}

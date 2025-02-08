package jsonparser

import (
	"fmt"
	"os"
	"testing"
)

func TestParse(t *testing.T) {

	validJsonFiles := []string{
		"tests/step1/valid.json",
		"tests/step2/valid.json",
		"tests/step2/valid2.json",
		"tests/step3/valid.json",
		"tests/step4/valid.json",
		"tests/step4/valid2.json",
	}

	for _, filename := range validJsonFiles {
		t.Run(fmt.Sprintf("%v should be valid json", filename), func(t *testing.T) {
			data, err := os.ReadFile(filename)
			if err != nil {
				panic(err)
			}

			res, err := Parse(data)

			if err != nil {
				t.Errorf("expected no error, got %v", err)
			}

			if res == nil {
				t.Error("expected non nil result")
			}
		})
	}

	invalidJsonFiles := []string{
		"tests/step1/invalid.json",
		"tests/step2/invalid.json",
		"tests/step2/invalid2.json",
		"tests/step3/invalid.json",
		"tests/step4/invalid.json",
	}

	for _, filename := range invalidJsonFiles {
		t.Run(fmt.Sprintf("%s should be invalid json", filename), func(t *testing.T) {
			data, err := os.ReadFile(filename)
			if err != nil {
				panic(err)
			}

			_, err = Parse(data)
			if err == nil {
				t.Errorf("expected error, got nil")
			}

		})
	}
}

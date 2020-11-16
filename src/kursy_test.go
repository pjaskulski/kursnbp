package main

import (
	"testing"
)

func TestCharactersAllowed(t *testing.T) {
	var text string = "1234-12-11"
	if !charAllowed(text, false) {
		t.Errorf("want: %t, got: %t", true, false)
	}
}

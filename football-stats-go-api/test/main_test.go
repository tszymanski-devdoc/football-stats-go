package main

import "testing"

func TestDummy(t *testing.T) {
	expected := 42
	actual := 42
	if actual != expected {
		t.Errorf("expected %d, got %d", expected, actual)
	}
}

package test

import "testing"

func TestColor(t *testing.T) {
	c, err := ParseColor("red")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c != ColorRed {
		t.Fatalf("unexpected value: %v", c)
	}

	c, err = ParseColor("RED")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c != ColorRed {
		t.Fatalf("unexpected value: %v", c)
	}

	_, err = ParseColor("x")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

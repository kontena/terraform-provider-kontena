package main

import (
	"testing"
)

func TestProvider(t *testing.T) {
	if err := provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

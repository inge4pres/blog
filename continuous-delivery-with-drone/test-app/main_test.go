package main

import (
	"strings"
	"testing"
)

func TestAFUnctionIsRunning(t *testing.T) {
	if s := testMe(); s == "" {
		t.Error("should net get empty string on test fn")
	}
	if s := testMe(); !strings.Contains(s, "👍🏻") {
		t.Error("message does not have 👍🏻")
	}
}

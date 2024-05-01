package testutil

import "testing"

func AssertEqual[T comparable](t *testing.T, got, want T) {
	t.Helper()
	if got != want {
		t.Errorf("Got: %v, Expected: %v\n", got, want)
	}
}

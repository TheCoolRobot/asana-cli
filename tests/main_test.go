package tests

import (
	"testing"
)

// TestPackageExists verifies the test package is set up correctly
func TestPackageExists(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}
	t.Log("Test package initialized successfully")
}
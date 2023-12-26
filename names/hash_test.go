package names

import "testing"

// TestJoinHash checks that nobody changed the hash function for Join.
func TestJoinHash(t *testing.T) {
	// DO NOT CHANGE THIS TEST!
	// This is intentionally a change-detection test. If it breaks, you messed up.
	parts := []string{"hello", "world"}
	want := "1dd41005"
	if got := Hash(parts); got != want {
		t.Fatalf("Hash(%v) = %q, want %q", parts, got, want)
	}
}

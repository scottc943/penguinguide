package distro

import "testing"

// This is a minimal sanity test for Detect.
// It does not care which distro you are running on,
// only that Detect works.
func TestDetectSetsFamily(t *testing.T) {
	t.Helper()

	info, err := Detect()
	if err != nil {
		t.Fatalf("Detect returned error: %v", err)
	}

	if info.Family == "" {
		t.Fatalf("expected Family to be set, got empty. Info: %+v", info)
	}
}

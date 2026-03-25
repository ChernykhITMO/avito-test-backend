package password

import "testing"

func TestHashAndCompare(t *testing.T) {
	manager := New(4)

	hash, err := manager.Hash("secret")
	if err != nil {
		t.Fatalf("Hash: %v", err)
	}

	if err := manager.Compare(hash, "secret"); err != nil {
		t.Fatalf("Compare: %v", err)
	}
}

package password

import "testing"

func TestHashAndCompare(t *testing.T) {
	hash, err := Hash("s3cret")
	if err != nil {
		t.Fatalf("Hash() error = %v", err)
	}
	if hash == "" || hash == "s3cret" {
		t.Fatalf("unexpected hash value: %q", hash)
	}

	if err := Compare(hash, "s3cret"); err != nil {
		t.Fatalf("Compare() should succeed, got %v", err)
	}

	if err := Compare(hash, "wrong"); err == nil {
		t.Fatal("Compare() should fail for wrong password")
	}
}

package auth

import "testing"

func TestGenerateAndParsePair(t *testing.T) {
	manager := NewManager("secret", "shop-go", 10, 24)

	accessToken, refreshToken, err := manager.GeneratePair(123, "user@example.com")
	if err != nil {
		t.Fatalf("GeneratePair() error = %v", err)
	}
	if accessToken == "" || refreshToken == "" {
		t.Fatal("expected non-empty tokens")
	}

	accessClaims, err := manager.Parse(accessToken)
	if err != nil {
		t.Fatalf("Parse(accessToken) error = %v", err)
	}
	if accessClaims.UserID != 123 || accessClaims.Email != "user@example.com" || accessClaims.Type != "access" {
		t.Fatalf("unexpected access claims: %+v", accessClaims)
	}

	refreshClaims, err := manager.Parse(refreshToken)
	if err != nil {
		t.Fatalf("Parse(refreshToken) error = %v", err)
	}
	if refreshClaims.Type != "refresh" {
		t.Fatalf("unexpected refresh type: %s", refreshClaims.Type)
	}
}

func TestParseRejectsInvalidToken(t *testing.T) {
	manager := NewManager("secret", "shop-go", 10, 24)

	if _, err := manager.Parse("bad.token"); err == nil {
		t.Fatal("expected invalid token error")
	}
}

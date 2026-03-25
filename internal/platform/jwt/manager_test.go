package jwt

import (
	"testing"
	"time"

	authmodel "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/auth/model"
)

func TestIssueAndParse(t *testing.T) {
	manager := New("secret", time.Hour)

	token, err := manager.Issue(authmodel.DummyUserID, authmodel.RoleUser)
	if err != nil {
		t.Fatalf("Issue: %v", err)
	}

	claims, err := manager.Parse(token)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if claims.UserID != authmodel.DummyUserID.String() {
		t.Fatalf("unexpected user id: %s", claims.UserID)
	}
	if claims.Role != string(authmodel.RoleUser) {
		t.Fatalf("unexpected role: %s", claims.Role)
	}
}

package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	authmodel "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/auth/model"
	jwtplatform "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/platform/jwt"
)

func TestAuthenticate(t *testing.T) {
	manager := jwtplatform.New("secret", time.Hour)
	token, err := manager.Issue(authmodel.DummyAdminID, authmodel.RoleAdmin)
	if err != nil {
		t.Fatalf("Issue: %v", err)
	}

	handler := Authenticate(manager)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/rooms/list", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

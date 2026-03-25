package config

import "testing"

func TestLoadUsesEnv(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/room_booking?sslmode=disable")
	t.Setenv("JWT_SECRET", "secret")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.HTTP.Port == "" {
		t.Fatal("expected default port")
	}
	if cfg.Database.URL == "" {
		t.Fatal("expected default database url")
	}
	if cfg.JWT.Secret == "" {
		t.Fatal("expected jwt secret")
	}
}

func TestLoadReturnsErrorWhenRequiredEnvMissing(t *testing.T) {
	t.Setenv("DATABASE_URL", "")
	t.Setenv("JWT_SECRET", "")

	if _, err := Load(); err == nil {
		t.Fatal("expected error")
	}
}

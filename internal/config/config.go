package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	LogLevel   string
	Database   DatabaseConfig
	JWT        JWTConfig
	Slots      SlotConfig
	Conference ConferenceConfig
	HTTP       HTTPConfig
}

type DatabaseConfig struct {
	URL                    string
	MaxConns               int32
	MinConns               int32
	MaxConnIdleTimeSeconds int
	MaxConnLifeTimeSeconds int
}

type JWTConfig struct {
	Secret     string
	TTLSeconds int
}

type SlotConfig struct {
	GenerationWindowDays  int
	RefillIntervalSeconds int
}

type HTTPConfig struct {
	Port                   string
	ReadTimeoutSeconds     int
	WriteTimeoutSeconds    int
	IdleTimeoutSeconds     int
	ShutdownTimeoutSeconds int
}

type ConferenceConfig struct {
	BaseURL  string
	MockMode string
}

func Load() (Config, error) {
	const op = "internal.config.Load"

	databaseURL, err := requiredEnv("DATABASE_URL")
	if err != nil {
		return Config{}, fmt.Errorf("%s: %w", op, err)
	}

	jwtSecret, err := requiredEnv("JWT_SECRET")
	if err != nil {
		return Config{}, fmt.Errorf("%s: %w", op, err)
	}

	return Config{
		LogLevel: getEnv("APP_LOG_LEVEL", "info"),
		Database: DatabaseConfig{
			URL:                    databaseURL,
			MaxConns:               int32(getEnvInt("DB_MAX_CONNS", 10)),
			MinConns:               int32(getEnvInt("DB_MIN_CONNS", 2)),
			MaxConnIdleTimeSeconds: getEnvInt("DB_MAX_CONN_IDLE_TIME_SECONDS", 300),
			MaxConnLifeTimeSeconds: getEnvInt("DB_MAX_CONN_LIFETIME_SECONDS", 1800),
		},
		JWT: JWTConfig{
			Secret:     jwtSecret,
			TTLSeconds: getEnvInt("JWT_TTL_SECONDS", 86400),
		},
		Slots: SlotConfig{
			GenerationWindowDays:  getEnvInt("SLOT_GENERATION_WINDOW_DAYS", 30),
			RefillIntervalSeconds: getEnvInt("SLOT_REFILL_INTERVAL_SECONDS", 300),
		},
		Conference: ConferenceConfig{
			BaseURL:  getEnv("CONFERENCE_BASE_URL", "https://conference.local/booking/"),
			MockMode: getEnv("CONFERENCE_MOCK_MODE", "ok"),
		},
		HTTP: HTTPConfig{
			Port:                   getEnv("APP_PORT", "8080"),
			ReadTimeoutSeconds:     getEnvInt("APP_READ_TIMEOUT_SECONDS", 5),
			WriteTimeoutSeconds:    getEnvInt("APP_WRITE_TIMEOUT_SECONDS", 10),
			IdleTimeoutSeconds:     getEnvInt("APP_IDLE_TIMEOUT_SECONDS", 60),
			ShutdownTimeoutSeconds: getEnvInt("APP_SHUTDOWN_TIMEOUT_SECONDS", 10),
		},
	}, nil
}

func getEnv(key string, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

func getEnvInt(key string, fallback int) int {
	value, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}

	return parsed
}

func requiredEnv(key string) (string, error) {
	value, ok := os.LookupEnv(key)
	if !ok || strings.TrimSpace(value) == "" {
		return "", fmt.Errorf("%s is required", key)
	}

	return value, nil
}

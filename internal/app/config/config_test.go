package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func setRequiredEnv(t *testing.T) {
	t.Helper()
	for name, value := range map[string]string{
		"ENV":                        "prod",
		"DB_USER":                    "runtime-user",
		"DB_PASS":                    "runtime-pass",
		"DB_HOST":                    "runtime-host",
		"DB_PORT":                    "5432",
		"DB_NAME":                    "runtime-db",
		"JWT_SECRET":                 "runtime-secret",
		"JWT_ACCESS_TOKEN_DURATION":  "15m",
		"JWT_REFRESH_TOKEN_DURATION": "168h",
	} {
		t.Setenv(name, value)
	}
}

func TestLoadUsesRuntimeEnvironmentWithoutDotEnv(t *testing.T) {
	t.Chdir(t.TempDir())
	setRequiredEnv(t)

	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Env != "prod" || cfg.Host != "runtime-host" || cfg.JwtSecret != "runtime-secret" {
		t.Fatalf("unexpected runtime config: %+v", cfg)
	}
}

func TestLoadUsesDotEnvOnlyForMissingLocalValues(t *testing.T) {
	t.Chdir(t.TempDir())
	setRequiredEnv(t)
	t.Setenv("ENV", "dev")
	t.Setenv("DB_HOST", "")
	if err := os.Unsetenv("DB_HOST"); err != nil {
		t.Fatal(err)
	}
	content := "DB_USER=file-user\nDB_HOST=file-host\n"
	if err := os.WriteFile(".env", []byte(content), 0600); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Host != "file-host" || cfg.User != "runtime-user" {
		t.Fatalf("unexpected local config: %+v", cfg)
	}
}

func TestLoadIgnoresDotEnvInProduction(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)
	setRequiredEnv(t)
	t.Setenv("DB_HOST", "")
	if err := os.Unsetenv("DB_HOST"); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, ".env"), []byte("DB_HOST=file-host\n"), 0600); err != nil {
		t.Fatal(err)
	}

	_, err := Load()
	if err == nil || !strings.Contains(err.Error(), "DB_HOST") {
		t.Fatalf("expected missing DB_HOST, got %v", err)
	}
}

func TestLoadRejectsEmptySecretAndInvalidDuration(t *testing.T) {
	t.Chdir(t.TempDir())
	setRequiredEnv(t)
	t.Setenv("JWT_SECRET", "")
	if _, err := Load(); err == nil || !strings.Contains(err.Error(), "JWT_SECRET") {
		t.Fatalf("expected missing JWT_SECRET, got %v", err)
	}
	t.Setenv("JWT_SECRET", "runtime-secret")
	t.Setenv("JWT_ACCESS_TOKEN_DURATION", "0")
	if _, err := Load(); err == nil || !strings.Contains(err.Error(), "JWT_ACCESS_TOKEN_DURATION") {
		t.Fatalf("expected invalid access duration, got %v", err)
	}
}

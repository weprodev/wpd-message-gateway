package app

import "testing"

func TestResolveConfigPath(t *testing.T) {
	t.Setenv("CONFIG_PATH", "")

	if got := resolveConfigPath("configs/custom.yml"); got != "configs/custom.yml" {
		t.Fatalf("explicit path: got %q", got)
	}

	t.Setenv("CONFIG_PATH", "configs/from-env.yml")
	if got := resolveConfigPath(""); got != "configs/from-env.yml" {
		t.Fatalf("CONFIG_PATH: got %q", got)
	}
}

func TestConfig_ListenAddr(t *testing.T) {
	cfg := &Config{Server: ServerConfig{Port: 8080}}
	if got := cfg.ListenAddr(); got != ":8080" {
		t.Fatalf("got %q", got)
	}

	cfg.Server.Port = 0
	if got := cfg.ListenAddr(); got != ":10101" {
		t.Fatalf("default port: got %q", got)
	}
}

func TestConfig_EnvironmentName(t *testing.T) {
	t.Setenv("APP_ENVIRONMENT", "")

	cfg := &Config{Environment: "staging"}
	if got := cfg.EnvironmentName(); got != "staging" {
		t.Fatalf("from config: got %q", got)
	}

	cfg.Environment = ""
	t.Setenv("APP_ENVIRONMENT", "test")
	if got := cfg.EnvironmentName(); got != "test" {
		t.Fatalf("from env: got %q", got)
	}
}

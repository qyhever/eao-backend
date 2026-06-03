package config

import "testing"

func TestNormalizeAndValidateRequiresFileAPIConfig(t *testing.T) {
	cfg := &Config{}
	if err := cfg.normalizeAndValidate(); err == nil {
		t.Fatal("expected missing file api config error")
	}
}

func TestNormalizeAndValidateDefaultsFileAPITimeout(t *testing.T) {
	cfg := &Config{
		ThirdParty: ThirdPartyConfig{
			FileAPI: FileAPIConfig{
				BaseURL: " http://localhost:6301 ",
				Secret:  " test-secret ",
			},
		},
	}

	if err := cfg.normalizeAndValidate(); err != nil {
		t.Fatalf("validate failed: %v", err)
	}
	if cfg.ThirdParty.FileAPI.TimeoutSeconds != 10 {
		t.Fatalf("expected default timeout 10, got %d", cfg.ThirdParty.FileAPI.TimeoutSeconds)
	}
	if cfg.ThirdParty.FileAPI.BaseURL != "http://localhost:6301" {
		t.Fatalf("unexpected base url: %q", cfg.ThirdParty.FileAPI.BaseURL)
	}
	if cfg.ThirdParty.FileAPI.Secret != "test-secret" {
		t.Fatalf("unexpected secret: %q", cfg.ThirdParty.FileAPI.Secret)
	}
}

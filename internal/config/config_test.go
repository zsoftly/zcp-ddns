package config_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/zsoftly/zcp-ddns/internal/config"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name            string
		envVars         map[string]string
		yamlContent     string
		wantErr         bool
		wantErrContains string
		check           func(t *testing.T, cfg *config.Config)
	}{
		{
			name: "TestLoad_ValidYAML",
			envVars: map[string]string{
				"ZCP_TOKEN": "test-token-123",
			},
			yamlContent: `
api_url: https://api.zcp.zsoftly.ca/api
interval: 300s
records:
  - zone: example.com
    name: home.example.com
    type: A
    ttl: 300
    source:
      type: public-ip
`,
			wantErr: false,
			check: func(t *testing.T, cfg *config.Config) {
				if cfg.APIURL != "https://api.zcp.zsoftly.ca/api" {
					t.Errorf("expected APIURL https://api.zcp.zsoftly.ca/api, got %s", cfg.APIURL)
				}
				if cfg.Interval != 5*time.Minute {
					t.Errorf("expected Interval 5m, got %v", cfg.Interval)
				}
				if len(cfg.Records) != 1 {
					t.Fatalf("expected 1 record, got %d", len(cfg.Records))
				}
				rec := cfg.Records[0]
				if rec.Zone != "example.com" || rec.Name != "home.example.com" || rec.Type != "A" || rec.TTL != 300 {
					t.Errorf("parsed record is incorrect: %+v", rec)
				}
				if rec.Source.Type != "public-ip" {
					t.Errorf("parsed source is incorrect: %+v", rec.Source)
				}
			},
		},
		{
			name: "TestLoad_Defaults",
			envVars: map[string]string{
				"ZCP_TOKEN": "test-token-123",
			},
			yamlContent: `
records:
  - zone: example.com
    name: home.example.com
    type: a
    source:
      type: public-ip
`,
			wantErr: false,
			check: func(t *testing.T, cfg *config.Config) {
				if cfg.APIURL != config.DefaultAPIURL {
					t.Errorf("expected default APIURL, got %s", cfg.APIURL)
				}
				if cfg.Interval != config.DefaultInterval {
					t.Errorf("expected default Interval, got %v", cfg.Interval)
				}
				if cfg.RunOnce != false {
					t.Errorf("expected default RunOnce false, got true")
				}
				if cfg.Records[0].TTL != config.DefaultTTL {
					t.Errorf("expected default TTL %d, got %d", config.DefaultTTL, cfg.Records[0].TTL)
				}
				if cfg.Records[0].Type != "A" {
					t.Errorf("expected type to be normalized to A, got %s", cfg.Records[0].Type)
				}
			},
		},
		{
			name: "TestLoad_EnvOverridesAPIURL",
			envVars: map[string]string{
				"ZCP_TOKEN":   "test-token-123",
				"ZCP_API_URL": "https://custom.api/override",
			},
			yamlContent: `
api_url: https://api.zcp.zsoftly.ca/api
records:
  - zone: example.com
    name: home.example.com
    type: A
    source:
      type: public-ip
`,
			wantErr: false,
			check: func(t *testing.T, cfg *config.Config) {
				if cfg.APIURL != "https://custom.api/override" {
					t.Errorf("expected APIURL to be overridden by env, got %s", cfg.APIURL)
				}
			},
		},
		{
			name: "TestLoad_EmptyAPIURLRestoresDefault",
			envVars: map[string]string{
				"ZCP_TOKEN": "test-token-123",
			},
			yamlContent: `
api_url: ""
records:
  - zone: example.com
    name: home.example.com
    type: A
    source:
      type: public-ip
`,
			wantErr: false,
			check: func(t *testing.T, cfg *config.Config) {
				if cfg.APIURL != config.DefaultAPIURL {
					t.Errorf("expected empty APIURL to restore default, got %s", cfg.APIURL)
				}
			},
		},
		{
			name:    "TestLoad_TokenFromEnv",
			envVars: map[string]string{"ZCP_TOKEN": "secret-token"},
			yamlContent: `
records:
  - zone: example.com
    name: home.example.com
    type: A
    source:
      type: public-ip
`,
			wantErr: false,
			check: func(t *testing.T, cfg *config.Config) {
				if cfg.Token != "secret-token" {
					t.Errorf("expected Token secret-token, got %s", cfg.Token)
				}
			},
		},
		{
			name:    "TestLoad_TokenFromYAML",
			envVars: map[string]string{"ZCP_TOKEN": ""},
			yamlContent: `
token: yaml-token-123
records:
  - zone: example.com
    name: home.example.com
    type: A
    source:
      type: public-ip
`,
			wantErr: false,
			check: func(t *testing.T, cfg *config.Config) {
				if cfg.Token != "yaml-token-123" {
					t.Errorf("expected Token from YAML, got %s", cfg.Token)
				}
			},
		},
		{
			name: "TestLoad_RunOnceFromEnv_Numeric",
			envVars: map[string]string{
				"ZCP_TOKEN":     "test-token-123",
				"ZCP_DDNS_ONCE": "1",
			},
			yamlContent: `
records:
  - zone: example.com
    name: home.example.com
    type: A
    source:
      type: public-ip
`,
			wantErr: false,
			check: func(t *testing.T, cfg *config.Config) {
				if cfg.RunOnce != true {
					t.Errorf("expected RunOnce true, got %v", cfg.RunOnce)
				}
			},
		},
		{
			name: "TestLoad_IntervalParsed",
			envVars: map[string]string{
				"ZCP_TOKEN": "test-token-123",
			},
			yamlContent: `
interval: 10m
records:
  - zone: example.com
    name: home.example.com
    type: A
    source:
      type: public-ip
`,
			wantErr: false,
			check: func(t *testing.T, cfg *config.Config) {
				if cfg.Interval != 10*time.Minute {
					t.Errorf("expected Interval 10m, got %v", cfg.Interval)
				}
			},
		},
		{
			name:    "TestLoad_MissingToken_ReturnsError",
			envVars: map[string]string{"ZCP_TOKEN": ""},
			yamlContent: `
records:
  - zone: example.com
    name: home.example.com
    type: A
    source:
      type: public-ip
`,
			wantErr:         true,
			wantErrContains: "ZCP_TOKEN is not set",
		},
		{
			name:    "TestLoad_EmptyRecords_ReturnsError",
			envVars: map[string]string{"ZCP_TOKEN": "test-token-123"},
			yamlContent: `
interval: 300s
`,
			wantErr:         true,
			wantErrContains: "at least one record",
		},
		{
			name:    "TestLoad_InvalidRecordType_ReturnsError",
			envVars: map[string]string{"ZCP_TOKEN": "test-token-123"},
			yamlContent: `
records:
  - zone: example.com
    name: home.example.com
    type: TXT
    source:
      type: public-ip
`,
			wantErr:         true,
			wantErrContains: "invalid type 'TXT'",
		},
		{
			name:    "TestLoad_InvalidSourceType_ReturnsError",
			envVars: map[string]string{"ZCP_TOKEN": "test-token-123"},
			yamlContent: `
records:
  - zone: example.com
    name: home.example.com
    type: A
    source:
      type: webhook
`,
			wantErr:         true,
			wantErrContains: "invalid source type 'webhook'",
		},
		{
			name:    "TestLoad_InterfaceSourceMissingName_ReturnsError",
			envVars: map[string]string{"ZCP_TOKEN": "test-token-123"},
			yamlContent: `
records:
  - zone: example.com
    name: home.example.com
    type: A
    source:
      type: interface
`,
			wantErr:         true,
			wantErrContains: "source name is required",
		},
		{
			name:    "TestLoad_StaticSourceMissingAddr_ReturnsError",
			envVars: map[string]string{"ZCP_TOKEN": "test-token-123"},
			yamlContent: `
records:
  - zone: example.com
    name: home.example.com
    type: A
    source:
      type: static
`,
			wantErr:         true,
			wantErrContains: "source addr is required",
		},
		{
			name:    "TestLoad_StaticSourceInvalidIP_ReturnsError",
			envVars: map[string]string{"ZCP_TOKEN": "test-token-123"},
			yamlContent: `
records:
  - zone: example.com
    name: home.example.com
    type: A
    source:
      type: static
      addr: not-an-ip
`,
			wantErr:         true,
			wantErrContains: "invalid static IP address",
		},
		{
			name:    "TestLoad_StaticSourceIPv6ForA_ReturnsError",
			envVars: map[string]string{"ZCP_TOKEN": "test-token-123"},
			yamlContent: `
records:
  - zone: example.com
    name: home.example.com
    type: A
    source:
      type: static
      addr: 2001:db8::1
`,
			wantErr:         true,
			wantErrContains: "requires an IPv4 address",
		},
		{
			name:    "TestLoad_StaticSourceIPv4ForAAAA_ReturnsError",
			envVars: map[string]string{"ZCP_TOKEN": "test-token-123"},
			yamlContent: `
records:
  - zone: example.com
    name: home.example.com
    type: AAAA
    source:
      type: static
      addr: 1.2.3.4
`,
			wantErr:         true,
			wantErrContains: "requires an IPv6 address",
		},
		{
			name:    "TestLoad_MissingZone_ReturnsError",
			envVars: map[string]string{"ZCP_TOKEN": "test-token-123"},
			yamlContent: `
records:
  - name: home.example.com
    type: A
    source:
      type: public-ip
`,
			wantErr:         true,
			wantErrContains: "zone is required",
		},
		{
			name:    "TestLoad_MissingName_ReturnsError",
			envVars: map[string]string{"ZCP_TOKEN": "test-token-123"},
			yamlContent: `
records:
  - zone: example.com
    type: A
    source:
      type: public-ip
`,
			wantErr:         true,
			wantErrContains: "name is required",
		},
		{
			name:    "TestLoad_InvalidYAML_ReturnsError",
			envVars: map[string]string{"ZCP_TOKEN": "test-token-123"},
			yamlContent: `
records:
  - zone: example.com
  bad_yaml: [
`,
			wantErr:         true,
			wantErrContains: "invalid YAML",
		},
		{
			name:    "TestLoad_NegativeInterval_ReturnsError",
			envVars: map[string]string{"ZCP_TOKEN": "test-token-123"},
			yamlContent: `
interval: -5m
records:
  - zone: example.com
    name: home.example.com
    type: A
    source:
      type: public-ip
`,
			wantErr:         true,
			wantErrContains: "interval must be positive",
		},
		{
			name:    "TestLoad_UnparseableInterval_ReturnsError",
			envVars: map[string]string{"ZCP_TOKEN": "test-token-123"},
			yamlContent: `
interval: banana
records:
  - zone: example.com
    name: home.example.com
    type: A
    source:
      type: public-ip
`,
			wantErr:         true,
			wantErrContains: "invalid interval format",
		},
		{
			name: "TestLoad_MultipleRecords",
			envVars: map[string]string{
				"ZCP_TOKEN": "test-token-123",
			},
			yamlContent: `
records:
  - zone: example.com
    name: home.example.com
    type: A
    source:
      type: public-ip
  - zone: example.com
    name: app.example.com
    type: AAAA
    source:
      type: interface
      name: eth0
`,
			wantErr: false,
			check: func(t *testing.T, cfg *config.Config) {
				if len(cfg.Records) != 2 {
					t.Fatalf("expected 2 records, got %d", len(cfg.Records))
				}
				r2 := cfg.Records[1]
				if r2.Name != "app.example.com" || r2.Type != "AAAA" || r2.Source.Type != "interface" || r2.Source.Name != "eth0" {
					t.Errorf("second record parsed incorrectly: %+v", r2)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("ZCP_TOKEN", "")
			t.Setenv("ZCP_API_URL", "")
			t.Setenv("ZCP_DDNS_ONCE", "")
			for k, v := range tt.envVars {
				t.Setenv(k, v)
			}

			dir := t.TempDir()
			configFile := filepath.Join(dir, "config.yaml")
			if tt.yamlContent != "" {
				if err := os.WriteFile(configFile, []byte(tt.yamlContent), 0644); err != nil {
					t.Fatalf("failed to write temp config file: %v", err)
				}
			}
			t.Setenv("ZCP_DDNS_CONFIG", configFile)

			cfg, err := config.Load()
			if (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.wantErrContains != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErrContains) {
					t.Errorf("expected error containing %q, got %v", tt.wantErrContains, err)
				}
			}
			if !tt.wantErr && tt.check != nil {
				tt.check(t, cfg)
			}
		})
	}
}

func TestLoad_FileNotFound_ReturnsError(t *testing.T) {
	t.Setenv("ZCP_TOKEN", "test-token-123")
	t.Setenv("ZCP_DDNS_CONFIG", "/does/not/exist.yaml")

	_, err := config.Load()
	if err == nil {
		t.Error("expected error when config file does not exist, got nil")
	} else if !strings.Contains(err.Error(), "config file not found") {
		t.Errorf("expected 'config file not found' error, got: %v", err)
	}
}

func TestLoad_EmptyFile_ReturnsError(t *testing.T) {
	t.Setenv("ZCP_TOKEN", "test-token-123")

	dir := t.TempDir()
	configFile := filepath.Join(dir, "empty.yaml")
	if err := os.WriteFile(configFile, []byte(""), 0644); err != nil {
		t.Fatalf("failed to write temp config file: %v", err)
	}
	t.Setenv("ZCP_DDNS_CONFIG", configFile)

	_, err := config.Load()
	if err == nil {
		t.Error("expected error for empty file, got nil")
	} else if !strings.Contains(err.Error(), "at least one record") {
		t.Errorf("expected 'at least one record' error, got: %v", err)
	}
}

func TestConfig_StringRedactsToken(t *testing.T) {
	cfg := config.Config{
		APIURL:   config.DefaultAPIURL,
		Token:    "super-secret-token-123",
		Interval: config.DefaultInterval,
	}

	s := cfg.String()
	if strings.Contains(s, "super-secret-token-123") {
		t.Errorf("Config.String() leaked the token: %s", s)
	}
	if !strings.Contains(s, "[REDACTED]") {
		t.Errorf("Config.String() missing [REDACTED] marker: %s", s)
	}
}

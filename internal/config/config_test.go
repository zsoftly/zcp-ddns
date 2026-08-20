package config_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/zsoftly/zcp-ddns/internal/config"
)

func TestLoad(t *testing.T) {
	// Table-driven test cases
	tests := []struct {
		name        string
		envVars     map[string]string
		yamlContent string
		wantErr     bool
		check       func(t *testing.T, cfg *config.Config)
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
    type: A
    source:
      type: public-ip
`,
			wantErr: false,
			check: func(t *testing.T, cfg *config.Config) {
				if cfg.APIURL != "https://api.zcp.zsoftly.ca/api" {
					t.Errorf("expected default APIURL https://api.zcp.zsoftly.ca/api, got %s", cfg.APIURL)
				}
				if cfg.Interval != 5*time.Minute {
					t.Errorf("expected default Interval 5m, got %v", cfg.Interval)
				}
				if cfg.RunOnce != false {
					t.Errorf("expected default RunOnce false, got true")
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
			name: "TestLoad_TokenFromEnv",
			envVars: map[string]string{
				"ZCP_TOKEN": "secret-token",
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
				if cfg.Token != "secret-token" {
					t.Errorf("expected Token secret-token, got %s", cfg.Token)
				}
			},
		},
		{
			name: "TestLoad_RunOnceFromEnv",
			envVars: map[string]string{
				"ZCP_TOKEN":     "test-token-123",
				"ZCP_DDNS_ONCE": "true",
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
			envVars: map[string]string{}, // No ZCP_TOKEN
			yamlContent: `
records:
  - zone: example.com
    name: home.example.com
    type: A
    source:
      type: public-ip
`,
			wantErr: true,
			check:   nil,
		},
		{
			name: "TestLoad_EmptyRecords_ReturnsError",
			envVars: map[string]string{
				"ZCP_TOKEN": "test-token-123",
			},
			yamlContent: `
interval: 300s
`,
			wantErr: true,
			check:   nil,
		},
		{
			name: "TestLoad_InvalidRecordType_ReturnsError",
			envVars: map[string]string{
				"ZCP_TOKEN": "test-token-123",
			},
			yamlContent: `
records:
  - zone: example.com
    name: home.example.com
    type: TXT
    source:
      type: public-ip
`,
			wantErr: true,
			check:   nil,
		},
		{
			name: "TestLoad_InvalidSourceType_ReturnsError",
			envVars: map[string]string{
				"ZCP_TOKEN": "test-token-123",
			},
			yamlContent: `
records:
  - zone: example.com
    name: home.example.com
    type: A
    source:
      type: webhook
`,
			wantErr: true,
			check:   nil,
		},
		{
			name: "TestLoad_InterfaceSourceMissingName_ReturnsError",
			envVars: map[string]string{
				"ZCP_TOKEN": "test-token-123",
			},
			yamlContent: `
records:
  - zone: example.com
    name: home.example.com
    type: A
    source:
      type: interface
`,
			wantErr: true,
			check:   nil,
		},
		{
			name: "TestLoad_MissingZone_ReturnsError",
			envVars: map[string]string{
				"ZCP_TOKEN": "test-token-123",
			},
			yamlContent: `
records:
  - name: home.example.com
    type: A
    source:
      type: public-ip
`,
			wantErr: true,
			check:   nil,
		},
		{
			name: "TestLoad_MissingName_ReturnsError",
			envVars: map[string]string{
				"ZCP_TOKEN": "test-token-123",
			},
			yamlContent: `
records:
  - zone: example.com
    type: A
    source:
      type: public-ip
`,
			wantErr: true,
			check:   nil,
		},
		{
			name: "TestLoad_InvalidYAML_ReturnsError",
			envVars: map[string]string{
				"ZCP_TOKEN": "test-token-123",
			},
			yamlContent: `
records:
  - zone: example.com
  bad_yaml: [
`,
			wantErr: true,
			check:   nil,
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
			// Clear all env vars to ensure isolated tests
			os.Clearenv()
			for k, v := range tt.envVars {
				os.Setenv(k, v)
			}

			// Create temporary YAML file
			dir := t.TempDir()
			configFile := filepath.Join(dir, "config.yaml")
			if tt.yamlContent != "" {
				if err := os.WriteFile(configFile, []byte(tt.yamlContent), 0644); err != nil {
					t.Fatalf("failed to write temp config file: %v", err)
				}
			}
			os.Setenv("ZCP_DDNS_CONFIG", configFile)

			cfg, err := config.Load()
			if (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && tt.check != nil {
				tt.check(t, cfg)
			}
		})
	}
}

func TestLoad_FileNotFound_ReturnsError(t *testing.T) {
	os.Clearenv()
	os.Setenv("ZCP_TOKEN", "test-token-123")
	os.Setenv("ZCP_DDNS_CONFIG", "/does/not/exist.yaml")

	_, err := config.Load()
	if err == nil {
		t.Error("expected error when config file does not exist, got nil")
	}
}

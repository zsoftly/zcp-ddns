package config

import (
	"errors"
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	APIURL   string        `yaml:"api_url"`
	Token    string        `yaml:"-"`
	Interval time.Duration `yaml:"-"`
	RunOnce  bool          `yaml:"-"`
	Records  []Record      `yaml:"records"`

	RawInterval string `yaml:"interval"`
	RawToken    string `yaml:"token"`
}

type Record struct {
	Zone   string `yaml:"zone"`
	Name   string `yaml:"name"`
	Type   string `yaml:"type"`
	TTL    int    `yaml:"ttl"`
	Source Source `yaml:"source"`
}

type Source struct {
	Type string `yaml:"type"`
	Name string `yaml:"name"`
	Addr string `yaml:"addr"`
}

func Load() (*Config, error) {
	cfg := &Config{
		APIURL:      "https://api.zcp.zsoftly.ca/api",
		RawInterval: "300s",
	}

	configPath := os.Getenv("ZCP_DDNS_CONFIG")
	if configPath == "" {
		configPath = "./zcp-ddns.yaml"
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("invalid YAML: %w", err)
	}

	if envAPI := os.Getenv("ZCP_API_URL"); envAPI != "" {
		cfg.APIURL = envAPI
	}

	cfg.Token = os.Getenv("ZCP_TOKEN")
	if cfg.Token == "" {
		cfg.Token = cfg.RawToken
	}

	if os.Getenv("ZCP_DDNS_ONCE") == "true" {
		cfg.RunOnce = true
	}

	if cfg.RawInterval == "" {
		cfg.RawInterval = "300s"
	}
	dur, err := time.ParseDuration(cfg.RawInterval)
	if err != nil {
		return nil, fmt.Errorf("invalid interval format: %w", err)
	}
	cfg.Interval = dur

	if cfg.Token == "" {
		return nil, errors.New("missing ZCP_TOKEN")
	}
	if len(cfg.Records) == 0 {
		return nil, errors.New("at least one record must be configured")
	}

	for i, r := range cfg.Records {
		if r.Zone == "" {
			return nil, fmt.Errorf("record %d: zone is required", i)
		}
		if r.Name == "" {
			return nil, fmt.Errorf("record %d: name is required", i)
		}
		if r.Type != "A" && r.Type != "AAAA" {
			return nil, fmt.Errorf("record %d: invalid type '%s', must be A or AAAA", i, r.Type)
		}
		if r.Source.Type != "public-ip" && r.Source.Type != "interface" && r.Source.Type != "static" {
			return nil, fmt.Errorf("record %d: invalid source type '%s'", i, r.Source.Type)
		}
		if r.Source.Type == "interface" && r.Source.Name == "" {
			return nil, fmt.Errorf("record %d: source name is required for interface type", i)
		}
	}

	return cfg, nil
}

// Package config loads and validates the zcp-ddns runtime configuration.
package config

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// DefaultAPIURL is the base URL for the ZCP API.
const DefaultAPIURL = "https://api.zcp.zsoftly.ca/api"

// DefaultInterval is the default record update interval.
const DefaultInterval = 5 * time.Minute

// DefaultConfigPath is the default path to the YAML configuration file.
const DefaultConfigPath = "./zcp-ddns.yaml"

// DefaultTTL is the default time-to-live for DNS records if not specified.
const DefaultTTL = 300

// Config holds the validated runtime configuration for zcp-ddns.
type Config struct {
	APIURL   string
	Token    string
	Interval time.Duration
	RunOnce  bool
	Records  []Record
}

// String implements fmt.Stringer to redact sensitive credentials from logs and output.
func (c Config) String() string {
	return fmt.Sprintf("{APIURL:%s Token:[REDACTED] Interval:%s RunOnce:%v Records:%v}",
		c.APIURL, c.Interval, c.RunOnce, c.Records)
}

// Record represents a single DNS record to be managed.
type Record struct {
	Zone   string
	Name   string
	Type   string
	TTL    int
	Source Source
}

// Source represents the method for determining the IP address.
type Source struct {
	Type string
	Name string
	Addr string
}

// fileConfig is a private struct used only for YAML unmarshaling to prevent
// raw token or interval fields from leaking into the public API.
type fileConfig struct {
	APIURL   string   `yaml:"api_url"`
	Interval string   `yaml:"interval"`
	Token    string   `yaml:"token"`
	Records  []Record `yaml:"records"`
}

// Load reads the YAML configuration file and environment variables,
// applies precedence rules, and returns a validated Config.
func Load() (*Config, error) {
	configPath := os.Getenv("ZCP_DDNS_CONFIG")
	if configPath == "" {
		configPath = DefaultConfigPath
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("config file not found at %s - set ZCP_DDNS_CONFIG or create the file", configPath)
		}
		return nil, fmt.Errorf("reading config file %s: %w", configPath, err)
	}

	var fc fileConfig
	if err := yaml.Unmarshal(data, &fc); err != nil {
		return nil, fmt.Errorf("invalid YAML in %s: %w", configPath, err)
	}

	cfg := &Config{
		APIURL:  fc.APIURL,
		Token:   fc.Token,
		Records: fc.Records,
	}

	// 1. Precedence: Env -> YAML -> Default
	if envAPI := os.Getenv("ZCP_API_URL"); envAPI != "" {
		cfg.APIURL = envAPI
	}
	if cfg.APIURL == "" {
		cfg.APIURL = DefaultAPIURL
	}

	if envToken := os.Getenv("ZCP_TOKEN"); envToken != "" {
		cfg.Token = envToken
	}

	if v := os.Getenv("ZCP_DDNS_ONCE"); v != "" {
		once, err := strconv.ParseBool(v)
		if err != nil {
			return nil, fmt.Errorf("invalid ZCP_DDNS_ONCE value %q: %w", v, err)
		}
		cfg.RunOnce = once
	}

	// Interval parsing and validation
	if fc.Interval == "" {
		cfg.Interval = DefaultInterval
	} else {
		dur, err := time.ParseDuration(fc.Interval)
		if err != nil {
			return nil, fmt.Errorf("invalid interval format: %w", err)
		}
		if dur <= 0 {
			return nil, fmt.Errorf("interval must be positive, got %s", fc.Interval)
		}
		cfg.Interval = dur
	}

	// 2. Global Validation
	if cfg.Token == "" {
		return nil, errors.New("ZCP_TOKEN is not set and no token key is in the config file - set the ZCP_TOKEN environment variable")
	}
	if len(cfg.Records) == 0 {
		return nil, errors.New("at least one record must be configured")
	}

	// 3. Record Validation
	for i, r := range cfg.Records {
		if r.Zone == "" {
			return nil, fmt.Errorf("record %d: zone is required", i)
		}
		if r.Name == "" {
			return nil, fmt.Errorf("record %d: name is required", i)
		}

		rType := strings.ToUpper(r.Type)
		if rType != "A" && rType != "AAAA" {
			return nil, fmt.Errorf("record %d: invalid type '%s', must be A or AAAA", i, r.Type)
		}
		cfg.Records[i].Type = rType

		if r.TTL == 0 {
			cfg.Records[i].TTL = DefaultTTL
		}

		sType := strings.ToLower(r.Source.Type)
		if sType != "public-ip" && sType != "interface" && sType != "static" {
			return nil, fmt.Errorf("record %d: invalid source type '%s'", i, r.Source.Type)
		}
		cfg.Records[i].Source.Type = sType

		if sType == "interface" && r.Source.Name == "" {
			return nil, fmt.Errorf("record %d: source name is required for interface type", i)
		}

		if sType == "static" {
			if r.Source.Addr == "" {
				return nil, fmt.Errorf("record %d: source addr is required for static type", i)
			}
			ip := net.ParseIP(r.Source.Addr)
			if ip == nil {
				return nil, fmt.Errorf("record %d: invalid static IP address '%s'", i, r.Source.Addr)
			}
			if rType == "A" && ip.To4() == nil {
				return nil, fmt.Errorf("record %d: type A requires an IPv4 address, got %s", i, r.Source.Addr)
			}
			if rType == "AAAA" && ip.To4() != nil {
				return nil, fmt.Errorf("record %d: type AAAA requires an IPv6 address, got %s", i, r.Source.Addr)
			}
		}
	}

	return cfg, nil
}

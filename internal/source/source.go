// Package source provides interfaces and implementations for IP address resolution.
package source

import (
	"context"
	"fmt"
	"net"
	"strings"
)

// Source represents a method for determining the current IP address for a DNS record.
type Source interface {
	// Resolve returns the current IP address for the record.
	Resolve(ctx context.Context) (net.IP, error)
}

// StaticSource implements Source by returning a fixed IP address.
type StaticSource struct {
	ip net.IP
}

// NewStaticSource creates a new StaticSource.
// It validates that addr is a valid IP and matches the family expected by recordType (A or AAAA).
func NewStaticSource(recordType, addr string) (*StaticSource, error) {
	recordType = strings.ToUpper(recordType)
	if recordType != "A" && recordType != "AAAA" {
		return nil, fmt.Errorf("invalid record type '%s', must be A or AAAA", recordType)
	}

	ip := net.ParseIP(addr)
	if ip == nil {
		return nil, fmt.Errorf("invalid static IP address: '%s'", addr)
	}
	if recordType == "A" && ip.To4() == nil {
		return nil, fmt.Errorf("type A requires an IPv4 address, got %s", addr)
	}
	if recordType == "AAAA" && ip.To4() != nil {
		return nil, fmt.Errorf("type AAAA requires an IPv6 address, got %s", addr)
	}
	return &StaticSource{ip: ip}, nil
}

// Resolve returns the fixed IP address.
func (s *StaticSource) Resolve(ctx context.Context) (net.IP, error) {
	return append(net.IP(nil), s.ip...), nil
}

package source_test

import (
	"context"
	"strings"
	"testing"

	"github.com/zsoftly/zcp-ddns/internal/source"
)

// TestStaticSource verifies the constructor and resolution logic for the static source provider.
func TestStaticSource(t *testing.T) {
	tests := []struct {
		name            string
		recordType      string
		addr            string
		wantErr         bool
		wantErrContains string
		expectedIP      string
	}{
		{
			name:       "TestStaticSource_ValidIPv4",
			recordType: "A",
			addr:       "192.0.2.1",
			wantErr:    false,
			expectedIP: "192.0.2.1",
		},
		{
			name:       "TestStaticSource_ValidIPv6",
			recordType: "AAAA",
			addr:       "2001:db8::1",
			wantErr:    false,
			expectedIP: "2001:db8::1",
		},
		{
			name:       "TestStaticSource_ValidIPv4_LowercaseType",
			recordType: "a",
			addr:       "192.0.2.1",
			wantErr:    false,
			expectedIP: "192.0.2.1",
		},
		{
			name:            "TestStaticSource_InvalidType_ReturnsError",
			recordType:      "CNAME",
			addr:            "192.0.2.1",
			wantErr:         true,
			wantErrContains: "invalid record type",
		},
		{
			name:            "TestStaticSource_InvalidAddr_ReturnsError",
			recordType:      "A",
			addr:            "not-an-ip",
			wantErr:         true,
			wantErrContains: "invalid static IP address",
		},
		{
			name:            "TestStaticSource_IPv6AddrForIPv4Record_ReturnsError",
			recordType:      "A",
			addr:            "2001:db8::1",
			wantErr:         true,
			wantErrContains: "requires an IPv4 address",
		},
		{
			name:            "TestStaticSource_IPv4AddrForIPv6Record_ReturnsError",
			recordType:      "AAAA",
			addr:            "192.0.2.1",
			wantErr:         true,
			wantErrContains: "requires an IPv6 address",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := source.NewStaticSource(tt.recordType, tt.addr)
			if (err != nil) != tt.wantErr {
				t.Fatalf("NewStaticSource() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				if !strings.Contains(err.Error(), tt.wantErrContains) {
					t.Errorf("error message %q does not contain %q", err.Error(), tt.wantErrContains)
				}
				return
			}

			ip, err := s.Resolve(context.Background())
			if err != nil {
				t.Fatalf("unexpected Resolve error: %v", err)
			}
			if ip.String() != tt.expectedIP {
				t.Errorf("expected IP %s, got %s", tt.expectedIP, ip.String())
			}

			// Verify immutability
			ip[0] = 255
			ip2, _ := s.Resolve(context.Background())
			if ip2.String() == ip.String() {
				t.Errorf("expected internal IP to be immutable, but was changed")
			}
		})
	}
}

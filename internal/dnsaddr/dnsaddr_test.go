package dnsaddr

import "testing"

func TestNormalize(t *testing.T) {
	tests := []struct {
		name    string
		in      string
		want    string
		wantErr bool
	}{
		{"bare ipv4", "8.8.8.8", "8.8.8.8:53", false},
		{"ipv4 with port", "8.8.8.8:5353", "8.8.8.8:5353", false},
		{"ipv4 trailing colon", "8.8.8.8:", "8.8.8.8:53", false},
		{"bare ipv6", "2001:4860:4860::8888", "[2001:4860:4860::8888]:53", false},
		{"bracketed ipv6 no port", "[2001:4860:4860::8888]", "[2001:4860:4860::8888]:53", false},
		{"bracketed ipv6 with port", "[2001:4860:4860::8888]:53", "[2001:4860:4860::8888]:53", false},
		{"ipv6 loopback bracketed", "[::1]", "[::1]:53", false},
		{"ipv4 mapped ipv6", "::ffff:8.8.8.8", "[::ffff:8.8.8.8]:53", false},
		{"ipv6 ending in port-like group", "2001:db8::1:53", "[2001:db8::1:53]:53", false},
		{"whitespace trimmed", "  8.8.8.8  ", "8.8.8.8:53", false},
		{"empty", "", "", true},
		{"whitespace only", "   ", "", true},
		{"non-ip host", "dns.example.com", "", true},
		{"non-ip host with port", "dns.example.com:53", "", true},
		{"port lower boundary", "8.8.8.8:1", "8.8.8.8:1", false},
		{"port upper boundary", "8.8.8.8:65535", "8.8.8.8:65535", false},
		// Port 0 dials without error but every query is discarded, so it must
		// not pass validation.
		{"port zero", "8.8.8.8:0", "", true},
		{"port just over max", "8.8.8.8:65536", "", true},
		{"port out of range", "8.8.8.8:70000", "", true},
		{"negative port", "8.8.8.8:-1", "", true},
		{"signed port", "8.8.8.8:+53", "", true},
		{"non-numeric port", "8.8.8.8:abc", "", true},
		{"named port", "8.8.8.8:domain", "", true},
		{"invalid ipv4 octet", "256.256.256.256", "", true},
		{"too many ipv4 octets", "1.2.3.4.5", "", true},
		{"unbracketed ipv6 with port", "2001:db8::1:5353:9999:1:2:3:4", "", true},
		{"garbage", "not-an-ip", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Normalize(tt.in)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("Normalize(%q): expected error, got nil (result %q)", tt.in, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("Normalize(%q): unexpected error: %v", tt.in, err)
			}
			if got != tt.want {
				t.Fatalf("Normalize(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

// TestNormalizeIsIdempotent guards the contract that the returned value is
// already canonical: feeding it back must not change it.
func TestNormalizeIsIdempotent(t *testing.T) {
	for _, in := range []string{"8.8.8.8", "  8.8.8.8:5353 ", "[::1]", "2001:db8::1"} {
		once, err := Normalize(in)
		if err != nil {
			t.Fatalf("Normalize(%q): unexpected error: %v", in, err)
		}
		twice, err := Normalize(once)
		if err != nil {
			t.Fatalf("Normalize(%q) (second pass on %q): unexpected error: %v", once, in, err)
		}
		if once != twice {
			t.Errorf("Normalize not idempotent for %q: %q -> %q", in, once, twice)
		}
	}
}

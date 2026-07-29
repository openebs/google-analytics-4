package usage

import (
	"context"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"
)

// fakeDNSServer starts a UDP listener that stands in for a DNS server and
// returns its address plus a channel closed when it receives a packet.
func fakeDNSServer(t *testing.T) (addr string, contacted <-chan struct{}) {
	t.Helper()

	listener, err := net.ListenPacket("udp4", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start UDP listener: %v", err)
	}
	t.Cleanup(func() { _ = listener.Close() })

	ch := make(chan struct{})
	go func() {
		buf := make([]byte, 512)
		if _, _, err := listener.ReadFrom(buf); err == nil {
			close(ch)
		}
	}()

	return listener.LocalAddr().String(), ch
}

// assertResolverContacts dials a hostname through client and reports whether the
// custom resolver reached the fake DNS server. The dial itself is expected to
// fail — only the resolver traffic matters.
func assertResolverContacts(t *testing.T, client *http.Client, contacted <-chan struct{}) {
	t.Helper()

	transport, ok := client.Transport.(*http.Transport)
	if !ok {
		t.Fatal("expected *http.Transport")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	go func() {
		conn, err := transport.DialContext(ctx, "tcp", "nonexistent.test:80")
		if err == nil {
			_ = conn.Close()
		}
	}()

	select {
	case <-contacted:
		// custom DNS server was reached — resolver is correctly wired
	case <-time.After(5 * time.Second):
		t.Fatal("custom DNS server was never contacted — resolver not wired")
	}
}

func TestHttpClientWithDns_EmptyDNS(t *testing.T) {
	// An empty DNS address is not handled here (the fallback to the default
	// transport lives in New); it must be rejected as invalid.
	_, err := httpClientWithDns("")
	if err == nil {
		t.Fatal("expected error for empty DNS address, got nil")
	}
}

func TestHttpClientWithDns_WithDNS(t *testing.T) {
	addr, contacted := fakeDNSServer(t)

	client, err := httpClientWithDns(addr)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	assertResolverContacts(t, client, contacted)
}

// TestHttpClientWithDns_DialsNormalizedAddr pins the resolver to the normalized
// address rather than the raw input: the untrimmed form below is not dialable,
// so the fake DNS server is only reached if normalization is applied first.
func TestHttpClientWithDns_DialsNormalizedAddr(t *testing.T) {
	addr, contacted := fakeDNSServer(t)

	padded := "  " + addr + "  "
	if _, err := net.Dial("udp", padded); err == nil {
		t.Fatalf("precondition failed: %q is unexpectedly dialable as-is", padded)
	}

	client, err := httpClientWithDns(padded)
	if err != nil {
		t.Fatalf("expected no error for %q, got %v", padded, err)
	}
	assertResolverContacts(t, client, contacted)
}

func TestHttpClientWithDns_NoPort(t *testing.T) {
	// A bare IP without a port is accepted; the port defaults to 53.
	client, err := httpClientWithDns("8.8.8.8")
	if err != nil {
		t.Fatalf("expected no error for portless DNS address, got %v", err)
	}
	if client == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestHttpClientWithDns_InvalidDNS(t *testing.T) {
	// Rejected inputs are covered exhaustively by the dnsaddr package tests;
	// here we only check that a bad value is refused and that the message names
	// the env var an operator has to fix.
	for _, in := range []string{"not-an-ip", "dns.example.com", "8.8.8.8:0", "8.8.8.8:70000"} {
		_, err := httpClientWithDns(in)
		if err == nil {
			t.Errorf("expected error for DNS address %q, got nil", in)
			continue
		}
		if !strings.Contains(err.Error(), DnsEnv) {
			t.Errorf("error for %q does not mention %s: %v", in, DnsEnv, err)
		}
	}
}

func TestToHumanSize(t *testing.T) {
	tests := map[string]struct {
		stringSize   string
		expectedSize string
		positiveTest bool
	}{
		"One Hundred Twenty Three thousand Four Hundred Fifty Six Tebibytes": {
			"123456 TiB",
			"121 PiB",
			true,
		},
		"One Gibibyte": {
			"1 GiB",
			"1.0 GiB",
			true,
		},
		"One Megabyte": {
			"1 MB",
			"977 KiB",
			true,
		},
		"One hundred four point five gigabyte": {
			"104.5 GB",
			"97 GiB",
			true,
		},
	}

	for testKey, testSuite := range tests {
		gotValue, err := toHumanSize(testSuite.stringSize)
		if (gotValue != testSuite.expectedSize || err != nil) && testSuite.positiveTest {
			t.Fatalf("Tests failed for %s, expected=%s, got=%s", testKey, testSuite.expectedSize, gotValue)
		}
	}
}

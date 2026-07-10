package usage

import (
	"context"
	"net"
	"net/http"
	"testing"
	"time"
)

func TestHttpClientWithDns_EmptyDNS(t *testing.T) {
	// An empty DNS address is no longer handled here (the fallback to the
	// default transport lives in New); it must be rejected as invalid.
	_, err := httpClientWithDns("")
	if err == nil {
		t.Fatal("expected error for empty DNS address, got nil")
	}
}

func TestHttpClientWithDns_WithDNS(t *testing.T) {
	// Start a local UDP listener acting as the fake DNS server.
	listener, err := net.ListenPacket("udp4", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start UDP listener: %v", err)
	}
	defer listener.Close()

	client, err := httpClientWithDns(listener.LocalAddr().String())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	transport, ok := client.Transport.(*http.Transport)
	if !ok {
		t.Fatal("expected *http.Transport")
	}

	// contacted is closed when the fake DNS server receives a packet.
	contacted := make(chan struct{})
	go func() {
		buf := make([]byte, 512)
		_ = listener.SetReadDeadline(time.Now().Add(2 * time.Second))
		if _, _, err := listener.ReadFrom(buf); err == nil {
			close(contacted)
		}
	}()

	// Dialing a hostname (not an IP) forces the custom resolver to contact our
	// fake DNS server. The dial itself will fail — that's expected.
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	transport.DialContext(ctx, "tcp", "nonexistent.test:80")

	select {
	case <-contacted:
		// custom DNS server was reached — resolver is correctly wired
	case <-time.After(2 * time.Second):
		t.Fatal("custom DNS server was never contacted — resolver not wired")
	}
}

func TestHttpClientWithDns_InvalidDNS(t *testing.T) {
	_, err := httpClientWithDns("8.8.8.8")
	if err == nil {
		t.Fatal("expected error for DNS address missing port, got nil")
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

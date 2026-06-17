/*
Copyright 2023 The OpenEBS Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package usage

import (
	"context"
	"net"
	"net/http"
	"testing"
	"time"
)

func TestNewHTTPClient_NoDNS(t *testing.T) {
	t.Setenv(DnsEnv, "")
	client, err := newHTTPClient()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if client == nil {
		t.Fatal("expected non-nil http.Client")
	}
	transport, ok := client.Transport.(*http.Transport)
	if !ok {
		t.Fatal("expected *http.Transport")
	}
	if transport.DialContext == nil {
		t.Fatal("expected DialContext to be set")
	}
}

func TestNewHTTPClient_WithDNS(t *testing.T) {
	// Start a local UDP listener acting as the fake DNS server.
	listener, err := net.ListenPacket("udp4", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start UDP listener: %v", err)
	}
	defer listener.Close()

	t.Setenv(DnsEnv, listener.LocalAddr().String())

	client, err := newHTTPClient()
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

func TestNewHTTPClient_InvalidDNS(t *testing.T) {
	t.Setenv(DnsEnv, "8.8.8.8")
	_, err := newHTTPClient()
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

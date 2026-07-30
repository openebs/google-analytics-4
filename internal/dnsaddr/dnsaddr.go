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

// Package dnsaddr normalizes DNS server addresses. It is internal so that the
// usage package and the example share one implementation: they previously kept
// separate copies that drifted apart.
package dnsaddr

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
)

// DefaultPort is applied when an address carries no port of its own.
const DefaultPort = "53"

// Normalize accepts a DNS server address with or without a port and returns it
// in host:port form, defaulting the port to 53 when omitted. Accepted forms:
//
//	8.8.8.8                      -> 8.8.8.8:53
//	8.8.8.8:5353                 -> 8.8.8.8:5353
//	2001:4860:4860::8888         -> [2001:4860:4860::8888]:53
//	[2001:4860:4860::8888]       -> [2001:4860:4860::8888]:53
//	[2001:4860:4860::8888]:5353  -> [2001:4860:4860::8888]:5353
//
// Surrounding whitespace is trimmed. The host must be a valid IP address —
// hostnames are rejected because resolving the resolver's own address would be
// circular. The port must be a plain decimal in the range 1-65535; port 0 is
// rejected because it dials successfully but silently discards every query.
func Normalize(dns string) (string, error) {
	dns = strings.TrimSpace(dns)
	if dns == "" {
		return "", errors.New("empty dns address")
	}

	host, port, err := net.SplitHostPort(dns)
	if err != nil {
		// No parseable port present: treat the whole value as a bare host and
		// apply the default DNS port. A bare IPv6 may still be bracketed
		// ("[::1]"), which SplitHostPort rejects for having no port.
		host, port = strings.TrimSuffix(strings.TrimPrefix(dns, "["), "]"), DefaultPort
	}
	if port == "" {
		port = DefaultPort
	}

	if net.ParseIP(host) == nil {
		return "", fmt.Errorf("invalid dns address %q: host must be a valid IP address", dns)
	}
	// ParseUint (unlike Atoi) rejects a sign prefix, so "+53" and "-1" do not
	// slip through as valid ports, and the bitSize caps the value at 65535.
	p, perr := strconv.ParseUint(port, 10, 16)
	if perr != nil || p == 0 {
		return "", fmt.Errorf("invalid dns address %q: port must be in the range 1-65535", dns)
	}

	// net.JoinHostPort brackets IPv6 hosts as needed.
	return net.JoinHostPort(host, strconv.FormatUint(p, 10)), nil
}

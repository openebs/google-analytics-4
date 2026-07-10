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
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/dustin/go-humanize"
)

func httpClientWithDns(dns string) (*http.Client, error) {
	if _, _, err := net.SplitHostPort(dns); err != nil {
		return nil, fmt.Errorf("invalid %s address %q: must be host:port", DnsEnv, dns)
	}
	dialer := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
		Resolver: &net.Resolver{
			PreferGo: true,
			Dial: func(ctx context.Context, network string, address string) (net.Conn, error) {
				d := net.Dialer{Timeout: 5 * time.Second}
				return d.DialContext(ctx, network, dns)
			},
		},
	}
	tr := http.DefaultTransport.(*http.Transport).Clone()
	tr.DialContext = dialer.DialContext
	return &http.Client{Transport: tr}, nil
}

// toHumanSize converts sizes to legible human sizes in IEC units.
func toHumanSize(size string) (string, error) {
	sizeInBytes, err := humanize.ParseBigBytes(size)
	return humanize.BigIBytes(sizeInBytes), err
}

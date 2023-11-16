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

package client

import (
	"context"
	"net"
	"net/http"
	"regexp"

	"github.com/openebs/lib-csi/pkg/common/errors"
)

var measurementIDMatcher = regexp.MustCompile(`^G-[a-zA-Z0-9]+$`)

type MeasurementClientOption func(*MeasurementClient) error

type MeasurementClient struct {
	HttpClient    *http.Client
	apiSecret     string
	measurementId string
	clientId      string
}

func NewMeasurementClient(opts ...MeasurementClientOption) (*MeasurementClient, error) {
	dialer := &net.Dialer{
		Resolver: &net.Resolver{
			PreferGo: true,
			Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
				dialer := net.Dialer{}
				return dialer.DialContext(ctx, network, "8.8.8.8:53")
			},
		},
	}
	c := &MeasurementClient{
		HttpClient: &http.Client{
			Transport: &http.Transport{
				DialContext: dialer.DialContext,
			},
		},
	}

	var err error
	for _, opt := range opts {
		err = opt(c)
		if err != nil {
			return nil, errors.Wrap(err, "failed to build MeasurementClient")
		}
	}

	return c, nil
}

func WithApiSecret(secret string) MeasurementClientOption {
	return func(s *MeasurementClient) error {
		if len(secret) == 0 {
			return errors.Errorf("failed to set api_secret: secret is an empty string")
		}

		s.apiSecret = secret
		return nil
	}
}

func WithMeasurementId(measurementId string) MeasurementClientOption {
	return func(s *MeasurementClient) error {
		if len(measurementId) == 0 {
			return errors.Errorf("failed to set measurement_id: id is an empty string")
		}

		if !measurementIDMatcher.MatchString(measurementId) {
			return errors.Errorf("Invalid measurement_id: %s", measurementId)
		}

		s.measurementId = measurementId
		return nil
	}
}

func WithClientId(clientId string) MeasurementClientOption {
	return func(s *MeasurementClient) error {
		if len(clientId) == 0 {
			return errors.Errorf("failed to set client_id: id is an empty string")
		}

		s.clientId = clientId
		return nil
	}
}

func (client *MeasurementClient) SetClientId(clientId string) {
	client.clientId = clientId
}

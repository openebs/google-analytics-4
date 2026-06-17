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
	"net/http"
	"regexp"

	"github.com/openebs/lib-csi/pkg/common/errors"
)

var MeasurementIDMatcher = regexp.MustCompile(`^G-[a-zA-Z0-9]+$`)

type MeasurementClientOption func(*MeasurementClient) error

type MeasurementClient struct {
	HttpClient    *http.Client
	apiSecret     string
	measurementId string
	clientId      string
}

func NewMeasurementClient(opts ...MeasurementClientOption) (*MeasurementClient, error) {

	c := &MeasurementClient{HttpClient: &http.Client{}}

	for _, opt := range opts {
		if err := opt(c); err != nil {
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

func WithHttpClient(client *http.Client) MeasurementClientOption {
	return func(s *MeasurementClient) error {
		if client == nil {
			return errors.Errorf("failed to set http client: client is nil")
		}
		s.HttpClient = client
		return nil
	}
}

func WithMeasurementId(measurementId string) MeasurementClientOption {
	return func(s *MeasurementClient) error {
		if len(measurementId) == 0 {
			return errors.Errorf("failed to set measurement_id: id is an empty string")
		}

		if !MeasurementIDMatcher.MatchString(measurementId) {
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

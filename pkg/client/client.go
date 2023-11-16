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
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/openebs/lib-csi/pkg/common/errors"

	"github.com/openebs/google-analytics-4/pkg/event"
	"github.com/openebs/google-analytics-4/pkg/payload"
)

func (c *MeasurementClient) Copy() *MeasurementClient {
	cpy := *c
	return &cpy
}

func (c *MeasurementClient) addFields(v url.Values) {
	v.Add("api_secret", c.apiSecret)
	v.Add("measurement_id", c.measurementId)
}

func (c *MeasurementClient) Send(event *event.OpenebsEvent) error {

	client := c.Copy()

	dataPayload, err := payload.NewPayload(
		payload.WithClientId(client.clientId),
		payload.WithOpenebsEvent(event),
	)

	if err != nil {
		return err
	}
	jsonData, err := json.Marshal(dataPayload)
	if err != nil {
		return err
	}

	gaUrl := "https://www.google-analytics.com/mp/collect"

	req, err := http.NewRequest("POST", gaUrl, bytes.NewReader(jsonData))
	v := req.URL.Query()
	client.addFields(v)
	req.URL.RawQuery = v.Encode()

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.HttpClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode/100 != 2 {
		return errors.Errorf("Rejected by Google with code %d", resp.StatusCode)
	}

	return nil
}

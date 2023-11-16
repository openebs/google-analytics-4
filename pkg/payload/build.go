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

package payload

import (
	"github.com/openebs/google-analytics-4/pkg/event"
	"github.com/openebs/lib-csi/pkg/common/errors"
)

type PayloadOption func(*Payload) error

type Payload struct {
	ClientId string     `json:"client_id"`
	Events   []ApiEvent `json:"events"`
}

type ApiEvent struct {
	Name   string             `json:"name"`
	Params event.OpenebsEvent `json:"params"`
}

func NewPayload(opts ...PayloadOption) (*Payload, error) {
	p := &Payload{}

	var err error
	for _, opt := range opts {
		err = opt(p)
		if err != nil {
			return nil, errors.Wrap(err, "failed to build Payload")
		}
	}

	return p, nil
}

func WithClientId(clientId string) PayloadOption {
	return func(p *Payload) error {
		if len(clientId) == 0 {
			return errors.Errorf("failed to set Payload clientId: id is an empty string")
		}

		p.ClientId = clientId
		return nil
	}
}

func WithOpenebsEvent(event *event.OpenebsEvent) PayloadOption {
	return func(p *Payload) error {
		p.Events = append(p.Events, ApiEvent{
			Name:   NormalizedEventName(event.CategoryStr()),
			Params: *event,
		})
		return nil
	}
}

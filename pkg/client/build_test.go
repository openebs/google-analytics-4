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
	"testing"
)

func TestWithHttpClient_NilClient(t *testing.T) {
	_, err := NewMeasurementClient(
		WithApiSecret("test-secret"),
		WithMeasurementId("G-TEST12345"),
		WithHttpClient(nil),
	)
	if err == nil {
		t.Fatal("expected error for nil http client, got nil")
	}
}

func TestWithHttpClient_ValidClient(t *testing.T) {
	custom := &http.Client{}
	client, err := NewMeasurementClient(
		WithApiSecret("test-secret"),
		WithMeasurementId("G-TEST12345"),
		WithHttpClient(custom),
	)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if client.HttpClient != custom {
		t.Fatal("expected HttpClient to be the provided client")
	}
}

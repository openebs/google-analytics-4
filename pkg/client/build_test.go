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

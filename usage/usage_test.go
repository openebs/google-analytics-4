package usage

import (
	"encoding/base64"
	"os"
	"reflect"
	"testing"
	"unicode/utf8"
)

func TestApiCreds(t *testing.T) {
	testCases := map[string]struct {
		idEnvValue     string
		idSet          bool
		secretEnvValue string
		secretSet      bool
		expectedId     string
		expectedSecret string
	}{
		"Missing both ENVs": {
			idEnvValue:     "",
			idSet:          false,
			secretEnvValue: "",
			secretSet:      false,
			expectedId:     DefaultMeasurementId,
			expectedSecret: DefaultApiSecret,
		},
		"Missing id ENV": {
			idEnvValue:     "",
			idSet:          false,
			secretEnvValue: base64.StdEncoding.EncodeToString([]byte("testSecret")),
			secretSet:      true,
			expectedId:     DefaultMeasurementId,
			expectedSecret: DefaultApiSecret,
		},
		"Missing secret ENV": {
			idEnvValue:     base64.StdEncoding.EncodeToString([]byte("G-testId")),
			idSet:          true,
			secretEnvValue: "",
			secretSet:      false,
			expectedId:     DefaultMeasurementId,
			expectedSecret: DefaultApiSecret,
		},
		"Empty id ENV": {
			idEnvValue:     "",
			idSet:          true,
			secretEnvValue: base64.StdEncoding.EncodeToString([]byte("testSecret")),
			secretSet:      true,
			expectedId:     DefaultMeasurementId,
			expectedSecret: DefaultApiSecret,
		},
		"Empty secret ENV": {
			idEnvValue:     base64.StdEncoding.EncodeToString([]byte("G-testId")),
			idSet:          true,
			secretEnvValue: "",
			secretSet:      true,
			expectedId:     DefaultMeasurementId,
			expectedSecret: DefaultApiSecret,
		},
		"Invalid base64 value for id ENV": {
			idEnvValue:     trimLastChar(base64.StdEncoding.EncodeToString([]byte("G-testId"))),
			idSet:          true,
			secretEnvValue: base64.StdEncoding.EncodeToString([]byte("testSecret")),
			secretSet:      true,
			expectedId:     DefaultMeasurementId,
			expectedSecret: DefaultApiSecret,
		},
		"Invalid base64 value for secret ENV": {
			idEnvValue:     base64.StdEncoding.EncodeToString([]byte("testId")),
			idSet:          true,
			secretEnvValue: trimLastChar(base64.StdEncoding.EncodeToString([]byte("testSecret"))),
			secretSet:      true,
			expectedId:     DefaultMeasurementId,
			expectedSecret: DefaultApiSecret,
		},
		"Valid ENVs set": {
			idEnvValue:     base64.StdEncoding.EncodeToString([]byte("G-testId")),
			idSet:          true,
			secretEnvValue: base64.StdEncoding.EncodeToString([]byte("testSecret")),
			secretSet:      true,
			expectedId:     "G-testId",
			expectedSecret: "testSecret",
		},
	}

	for k, v := range testCases {
		k, v := k, v
		t.Run(k, func(t *testing.T) {
			if v.idSet {
				if os.Setenv(MeasurementIdEnv, v.idEnvValue) != nil {
					t.Errorf("failed to set env '%s' to '%s' for test case '%s'",
						MeasurementIdEnv, v.idEnvValue, k,
					)
				}
			} else {
				if os.Unsetenv(MeasurementIdEnv) != nil {
					t.Errorf("failed to unset env '%s' for test case '%s'", MeasurementIdEnv, k)
				}
			}
			if v.secretSet {
				if os.Setenv(ApiSecretEnv, v.secretEnvValue) != nil {
					t.Errorf("failed to set env '%s' to '%s' for test case '%s'",
						ApiSecretEnv, v.secretEnvValue, k,
					)
				}
			} else {
				if os.Unsetenv(ApiSecretEnv) != nil {
					t.Errorf("failed to unset env '%s' for test case '%s'", ApiSecretEnv, k)
				}
			}
			observedId, observedSecret := apiCreds()
			if !reflect.DeepEqual(observedId, v.expectedId) {
				t.Errorf("apiCreds() id mismatch: expected '%s', observed '%s'",
					v.expectedId, observedId,
				)
			}
			if !reflect.DeepEqual(observedSecret, v.expectedSecret) {
				t.Errorf("apiCreds() secret mismatch: expected '%s', observed '%s'",
					v.expectedSecret, observedSecret,
				)
			}
		})
	}
}

func TestNewDnsEnv(t *testing.T) {
	// New must never return nil because of the DNS env var: consumers chain
	// directly off it (analytics.New().CommonBuild(...)) and would panic.
	testCases := map[string]struct {
		dnsEnvValue string
		dnsSet      bool
		// customTransport is true when the DNS override is expected to be
		// installed, false when New must fall back to the default transport.
		customTransport bool
	}{
		"Missing dns ENV":    {dnsSet: false, customTransport: false},
		"Empty dns ENV":      {dnsEnvValue: "", dnsSet: true, customTransport: false},
		"Whitespace dns ENV": {dnsEnvValue: "   ", dnsSet: true, customTransport: false},
		"Non-IP dns ENV":     {dnsEnvValue: "dns.example.com", dnsSet: true, customTransport: false},
		"Bad port dns ENV":   {dnsEnvValue: "8.8.8.8:70000", dnsSet: true, customTransport: false},
		"Zero port dns ENV":  {dnsEnvValue: "8.8.8.8:0", dnsSet: true, customTransport: false},
		"Valid dns ENV":      {dnsEnvValue: "8.8.8.8", dnsSet: true, customTransport: true},
		"Valid dns ENV port": {dnsEnvValue: "8.8.8.8:5353", dnsSet: true, customTransport: true},
		"Valid ipv6 dns ENV": {dnsEnvValue: "[2001:4860:4860::8888]", dnsSet: true, customTransport: true},
	}

	for k, v := range testCases {
		k, v := k, v
		t.Run(k, func(t *testing.T) {
			// Use the default measurement credentials for every case; t.Setenv
			// restores whatever the environment held before.
			t.Setenv(MeasurementIdEnv, "")
			t.Setenv(ApiSecretEnv, "")
			if v.dnsSet {
				t.Setenv(DnsEnv, v.dnsEnvValue)
			} else {
				t.Setenv(DnsEnv, "")
				if err := os.Unsetenv(DnsEnv); err != nil {
					t.Fatalf("failed to unset env '%s': %v", DnsEnv, err)
				}
			}

			u := New()
			if u == nil {
				t.Fatalf("New() returned nil for test case '%s'", k)
			}
			if u.OpenebsEventBuilder == nil {
				t.Errorf("New() returned a nil OpenebsEventBuilder for test case '%s'", k)
			}
			if u.AnalyticsClient == nil {
				t.Fatalf("New() returned a nil AnalyticsClient for test case '%s'", k)
			}
			if u.AnalyticsClient.HttpClient == nil {
				t.Fatalf("New() returned a nil HttpClient for test case '%s'", k)
			}

			// The DNS override is only observable as a transport that differs
			// from the stock one, so assert on that rather than on non-nilness.
			isCustom := u.AnalyticsClient.HttpClient.Transport != nil
			if isCustom != v.customTransport {
				t.Errorf("test case '%s': custom transport installed = %v, want %v",
					k, isCustom, v.customTransport,
				)
			}
		})
	}
}

func trimLastChar(s string) string {
	r, size := utf8.DecodeLastRuneInString(s)
	if r == utf8.RuneError && (size == 0 || size == 1) {
		size = 0
	}
	return s[:len(s)-size]
}

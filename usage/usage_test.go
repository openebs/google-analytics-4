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

func trimLastChar(s string) string {
	r, size := utf8.DecodeLastRuneInString(s)
	if r == utf8.RuneError && (size == 0 || size == 1) {
		size = 0
	}
	return s[:len(s)-size]
}

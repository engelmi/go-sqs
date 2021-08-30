// +build integration

package internal_test

import (
	"go-sqs/internal"
	"testing"
)

func Test_NewClient(t *testing.T) {

	type clientConfig struct {
		region   string
		endpoint string
		queue    string
	}

	type testcase struct {
		input     clientConfig
		wantError bool
	}

	testcases := map[string]testcase{
		"should create client for valid configuration": {
			input: clientConfig{
				region:   "eu-central-1",
				endpoint: "http://localhost:9324",
				queue:    "pink_panther",
			},
			wantError: false,
		},
		"should return error for unknown endpoint": {
			input: clientConfig{
				region:   "eu-central-1",
				endpoint: "http://unknown_endpoint:9324",
				queue:    "pink_panther",
			},
			wantError: true,
		},
		"should return error for unknown queue": {
			input: clientConfig{
				region:   "eu-central-1",
				endpoint: "http://localhost:9324",
				queue:    "unknown_queue",
			},
			wantError: true,
		},
	}

	for tcName, tc := range testcases {
		t.Run(tcName, func(t *testing.T) {
			_, err := internal.NewClient(tc.input.region, tc.input.endpoint, tc.input.queue)

			if tc.wantError && err == nil {
				t.Fatal("Expected error, but got <nil>")
			}
			if !tc.wantError && err != nil {
				t.Fatalf("Expected no error, but got '%s'", err.Error())
			}
		})
	}
}

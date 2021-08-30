// +build unit

package gosqs_test

import (
	gosqs "go-sqs"
	"testing"
)

func Test_MarshalToJson(t *testing.T) {

	type samplePayload struct {
		Text   string  `json:"text"`
		Number float64 `json:"number"`
	}

	type testcase struct {
		payload interface{}
		wantErr bool
	}

	zeroValue := 0.0
	testcases := map[string]testcase{
		"should marshal a valid payload successfully ": {
			payload: samplePayload{
				Text:   "payload with inf value",
				Number: 3.0,
			},
			wantErr: false,
		},
		"should return error if marshalling of payload fails ": {
			payload: samplePayload{
				Text:   "payload with inf value",
				Number: 3.0 / zeroValue,
			},
			wantErr: true,
		},
	}

	for tcName, tc := range testcases {
		t.Run(tcName, func(t *testing.T) {
			_, err := gosqs.MarshalToJson(tc.payload)

			if tc.wantErr && err == nil {
				t.Fatal("Expected error, but got <nil>")
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("Expected no error, but got '%s'", err.Error())
			}
		})
	}
}

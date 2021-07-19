package handlers

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type echoTest struct {
	input    string
	output   string
	httpCode int
}

var (
	// Input for table driven tests
	echoTests = []echoTest{
		{`{"upload": "xyz", "username": "xyz"}`,
			`{"echoed": "true", "upload": "xyz", "username": "xyz"}`,
			200,
		},
		{`{"echoed": "false", "upload": "xyz", "username": "xyz"}`,
			`{"echoed": "true", "upload": "xyz", "username": "xyz"}`,
			200,
		},
		{`{}`,
			`{"echoed": "true"}`,
			200,
		},
		{`{"upload",,,,,,,,,,,,: "xyz", "username": "xyz"}`,
			`{"error":"Malformed request, input must be valid JSON"}`,
			400,
		},
		{`{"echoed": "true", "upload": "xyz", "username": "xyz"}`,
			`{"error":"Top level echoed field is already set to true"}`,
			400,
		},
	}
)

// sendReq helper function to interact with a test mux Router and send
// an artibrary slice of bytes. It returns the HTTP response code and
// the response body as string
func sendReq(r *mux.Router, p []byte) (int, string) {
	req, _ := http.NewRequest("POST", "/api/v1/echo", bytes.NewBuffer(p))
	res := httptest.NewRecorder()
	r.ServeHTTP(res, req)
	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)
	return res.Code, bodyString
}

// TestEchoHandler integration test for the doEcho function
// we test the doEcho function via a mux handler and an httptest.ResponseRecorder
func TestEchoHandler(t *testing.T) {
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/echo", EchoHandler).Methods(http.MethodPost)

	for _, tt := range echoTests {
		t.Logf("Running test with input %v", tt.input)
		t.Logf("Expecting output %v", tt.output)
		resCode, bodyString := sendReq(r, []byte(tt.input))
		assert.Equal(t, tt.httpCode, resCode)
		require.JSONEq(t, tt.output, bodyString)
	}
}

// TestDoEcho unit test for the doEcho function
func TestDoEcho(t *testing.T) {
	for _, tt := range echoTests {
		t.Logf("Running test with input %v", tt.input)
		t.Logf("Expecting output %v", tt.output)
		r := ioutil.NopCloser(strings.NewReader(tt.input))
		ret := doEcho(r)
		assert.Equal(t, tt.httpCode, ret.HttpCode)
		require.JSONEq(t, tt.output, string(ret.Output))

	}
}

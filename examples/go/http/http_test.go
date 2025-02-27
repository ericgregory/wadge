//go:generate go run go.wasmcloud.dev/wadge/cmd/wadge-bindgen-go

package main_test

import (
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.wasmcloud.dev/wadge"
	incominghandler "go.wasmcloud.dev/wadge/examples/go/http/bindings/wasi/http/incoming-handler"
	"go.wasmcloud.dev/wadge/wadgehttp"
)

func TestIncomingHandler(t *testing.T) {
	wadge.RunTest(t, func() {
		req, err := http.NewRequest("", "test", nil)
		if err != nil {
			t.Fatalf("failed to create new HTTP request: %s", err)
		}
		req.Header.Add("foo", "bar")
		req.Header.Add("foo", "baz")
		req.Header.Add("key", "value")
		resp, err := wadgehttp.HandleIncomingRequest(incominghandler.Exports.Handle, req)
		if err != nil {
			t.Fatalf("failed to handle incoming HTTP request: %s", err)
		}
		assert.Equal(t, 200, resp.StatusCode)
		assert.Equal(t, http.Header{
			"foo": {
				"bar",
				"baz",
			},
			"key": {
				"value",
			},
		}, resp.Header)
		buf, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("failed to read HTTP response body: %s", err)
		}
		if err := resp.Body.Close(); err != nil {
			t.Fatalf("failed to close response body: %s", err)
		}
		assert.Equal(t, []byte("hello world"), buf)
	})
}

//go:build test

package testhelp

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
	"testing"
)

func DescribeHttpResponse(response *http.Response, body *bytes.Buffer) string {
	sb := &strings.Builder{}
	fmt.Fprintln(sb, "Response:")

	fmt.Fprintf(sb, "%s %s\n", response.Proto, response.Status)

	fmt.Fprintln(sb, "Headers:")
	for k, v := range response.Header {
		fmt.Fprintf(sb, "  %s: %s\n", k, strings.Join(v, ", "))
	}

	fmt.Fprintf(sb, "Content Length: %d\n", response.ContentLength)

	if body.Len() > 0 {
		fmt.Fprintln(sb, "--- Body [START] ---")
		sb.Write(body.Bytes())
		fmt.Fprintln(sb, "--- Body [END] ---")
	}

	return sb.String()
}

func TestLogHttpResponse(t testing.TB, response *http.Response, body *bytes.Buffer) {
	t.Helper()
	t.Log(DescribeHttpResponse(response, body))
}

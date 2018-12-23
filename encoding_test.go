package httpencoding

import (
	"net/http"
	"reflect"
	"testing"
)

type testEncodings []string

func (t *testEncodings) Handle(encoding Encoding) bool {
	*t = append(*t, string(encoding))
	return false
}

func TestOrder(t *testing.T) {
	for n, test := range []struct {
		AcceptEncoding string
		Encodings      testEncodings
	}{
		{"", testEncodings{""}},
		{"gzip", testEncodings{"gzip", ""}},
		{"gzip, deflate, br", testEncodings{"gzip", "deflate", "br", ""}},
		{"gzip, deflate;q=0.5, br;q=1.9", testEncodings{"br", "gzip", "deflate", ""}},
		{"identity, gzip, deflate;q=0.5, br;q=1.9", testEncodings{"br", "", "gzip", "deflate"}},
		{"gzip, br, identity;q=0", testEncodings{"gzip", "br"}},
		{"gzip, br, identity;q=0, *", testEncodings{"gzip", "br"}},
		{"gzip, *, br", testEncodings{"gzip", "", "br"}},
	} {
		te := make(testEncodings, 0, len(test.Encodings))
		HandleEncoding(&http.Request{
			Header: http.Header{
				acceptEncoding: []string{test.AcceptEncoding},
			},
		}, &te)
		if !reflect.DeepEqual(te, test.Encodings) {
			t.Errorf("test %d: expecting %v, got %v", n+1, test.Encodings, te)
		}
	}
}

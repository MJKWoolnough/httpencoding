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
		{"", testEncodings{"*", ""}},
		{"gzip", testEncodings{"gzip", ""}},
		{"gzip, deflate, br", testEncodings{"gzip", "deflate", "br", ""}},
		{"gzip, deflate;q=0.5, br;q=0.9", testEncodings{"gzip", "br", "deflate", ""}},
		{"identity, gzip, deflate;q=0.5, br;q=0.9", testEncodings{"", "gzip", "br", "deflate"}},
		{"gzip, br, identity;q=0", testEncodings{"gzip", "br"}},
		{"gzip, br, identity;q=0, *", testEncodings{"gzip", "br", "*;"}},
		{"gzip, *, br", testEncodings{"gzip", "*", "br", ""}},
		{"gzip, *, br;q=0, bzip;q=0", testEncodings{"gzip", "*;br;bzip", ""}},
		{"gzip, *, br;q=0, identity;q=0, bzip;q=0", testEncodings{"gzip", "*;br;;bzip"}},
		{"gzip;", testEncodings{""}},
		{"gzip, *, gzip;q=0", testEncodings{"gzip", "*", ""}},
		{"gzip, , bzip", testEncodings{"gzip", "bzip", ""}},
		{"*;q=0", testEncodings{}},
		{"*, *", testEncodings{"*", ""}},
		{"gzip;q=1, bzip;q=2", testEncodings{"gzip", ""}},
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

func TestIsDisallowedInWildcard(t *testing.T) {
	for n, test := range []struct {
		Wildcard, Enc Encoding
		Match         bool
	}{
		{"*", "gzip", false},
		{"*;gzip", "gzip", true},
		{"*;bzip;gzip", "gzip", true},
		{"*;bzip", "gzip", false},
		{"*;bzip", "", false},
		{"*;bzip;", "", true},
		{"*;;bzip", "", true},
	} {
		if IsDisallowedInWildcard(test.Wildcard, test.Enc) != test.Match {
			t.Errorf("test %d: expecting %v, got %v", n+1, test.Match, !test.Match)
		}
	}
}

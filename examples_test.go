package httpencoding_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"

	"vimagination.zapto.org/httpencoding"
)

func Example() {
	handler := func(w http.ResponseWriter, r *http.Request) {
		enc, ok := httpencoding.Negotiate(r, "gzip", "br", "")
		if !ok {
			enc = "none"
		} else if enc == "" {
			enc = "identity"
		}

		io.WriteString(w, string(enc))
	}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/", nil)
	handler(w, r)
	fmt.Println(w.Body)

	w = httptest.NewRecorder()
	r.Header.Set("Accept-encoding", "identity")
	handler(w, r)
	fmt.Println(w.Body)

	w = httptest.NewRecorder()
	r.Header.Set("Accept-encoding", "gzip, identity")
	handler(w, r)
	fmt.Println(w.Body)

	w = httptest.NewRecorder()
	r.Header.Set("Accept-encoding", "gzip;q=0.5, identity;q=0.6")
	handler(w, r)
	fmt.Println(w.Body)

	w = httptest.NewRecorder()
	r.Header.Set("Accept-encoding", "identity;q=0")
	handler(w, r)
	fmt.Println(w.Body)

	// Output:
	// gzip
	// identity
	// gzip
	// identity
	// none
}

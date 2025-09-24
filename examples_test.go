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
		if !httpencoding.HandleEncoding(r, httpencoding.HandlerFunc(func(e httpencoding.Encoding) bool {
			if e == "gzip" || httpencoding.IsWildcard(e) && !httpencoding.IsDisallowedInWildcard(e, "gzip") {
				io.WriteString(w, "gzip")
			} else if e == "" || httpencoding.IsWildcard(e) && !httpencoding.IsDisallowedInWildcard(e, "") {
				io.WriteString(w, "identity")
			} else {
				return false
			}

			return true
		})) {
			io.WriteString(w, "none")
		}
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

# httpencoding

[![CI](https://github.com/MJKWoolnough/httpencoding/actions/workflows/go-checks.yml/badge.svg)](https://github.com/MJKWoolnough/httpencoding/actions)
[![Go Reference](https://pkg.go.dev/badge/vimagination.zapto.org/httpencoding.svg)](https://pkg.go.dev/vimagination.zapto.org/httpencoding)
[![Go Report Card](https://goreportcard.com/badge/vimagination.zapto.org/httpencoding)](https://goreportcard.com/report/vimagination.zapto.org/httpencoding)

--
    import "vimagination.zapto.org/httpencoding"

Package httpencoding provides a function to deal with the Accept-Encoding header.

## Highlights

 - Simple handling of `Accept-Encoding` HTTP header.
 - Supports identity, wildcards, and q-values.

## Usage

```go
package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"

	"vimagination.zapto.org/httpencoding"
)

func main() {
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
```

## Documentation

Full API docs can be found at:

https://pkg.go.dev/vimagination.zapto.org/httpencoding

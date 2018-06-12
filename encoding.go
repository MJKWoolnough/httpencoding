// Package httpencoding provides a function to deal with the Accept-Encoding
// header.
package httpencoding // import "vimagination.zapto.org/httpencoding"

import (
	"net/http"
	"sort"
	"strconv"
	"strings"
)

const (
	acceptEncoding   = "Accept-Encoding"
	anyEncoding      = "*"
	identityEncoding = "identity"
	acceptSplit      = ","
	partSplit        = ";"
	weightPrefix     = "q="
)

type encodings []encoding

func (e encodings) Len() int {
	return len(e)
}

func (e encodings) Less(i, j int) bool {
	return e[j].weight < e[i].weight
}

func (e encodings) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

type encoding struct {
	encoding string
	weight   uint16
}

// Handler provides an interface to handle an encoding.
//
// The encoding string (e.g. gzip, br, deflate) is passed to the handler, which
// is expected to return true if no more encodings are required and false
// otherwise.
//
// The empty string "" is used to signify the identity encoding, or plain text
type Handler interface {
	Handle(encoding string) bool
}

// HandlerFunc wraps a func to make it satisfy the Handler interface
type HandlerFunc func(string) bool

// Handle calls the underlying func
func (h HandlerFunc) Handle(e string) bool {
	return h(e)
}

// InvalidEncoding writes the 406 header
func InvalidEncoding(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotAcceptable)
}

// HandleEncoding will process the Accept-Encoding header and calls the given
// handler for each encoding until the handler returns true.
//
// This function returns true when the Handler returns true, false otherwise
//
// For the identity (plain text) encoding the encoding string will be the
// empty string.
//
// The wildcard encoding (*) is currently treated as identity when there is no
// independent identity encoding specified; otherwise, it is ignored.
func HandleEncoding(r *http.Request, h Handler) bool {
	acceptHeader := r.Header.Get(acceptEncoding)
	accepts := make(encodings, 0, strings.Count(acceptHeader, acceptSplit)+1)
	allowIdentity := true
	hasIdentity := false
Loop:
	for _, accept := range strings.Split(acceptHeader, acceptSplit) {
		parts := strings.Split(strings.TrimSpace(accept), partSplit)
		name := strings.TrimSpace(parts[0])
		if name == "" {
			continue
		}
		var (
			qVal float64 = 1
			err  error
		)
		for _, part := range parts[1:] {
			if strings.HasPrefix(strings.TrimSpace(part), weightPrefix) {
				qVal, err = strconv.ParseFloat(part[len(weightPrefix):], 32)
				if err != nil || qVal < 0 || qVal >= 2 {
					continue Loop
				}
				break
			}
		}
		name = strings.ToLower(name)
		weight := uint16(qVal * 1000)
		if name == identityEncoding {
			allowIdentity = weight != 0
			hasIdentity = true
		}
		accepts = append(accepts, encoding{
			encoding: name,
			weight:   weight,
		})
	}
	sort.Stable(accepts)
	for _, accept := range accepts {
		switch accept.encoding {
		case identityEncoding:
			if accept.weight != 0 {
				if h.Handle("") {
					return true
				}
			}
			allowIdentity = false
		case anyEncoding:
			if !hasIdentity {
				if accept.weight != 0 {
					if h.Handle("") {
						return true
					}
				}
				allowIdentity = false
			}
		default:
			if h.Handle(accept.encoding) {
				return true
			}
		}
	}
	if allowIdentity {
		if h.Handle("") {
			return true
		}
	}
	return false
}

// ClearEncoding removes the Accept-Encoding header so that any further
// attempts to establish an encoding will simply used the default, plain text,
// encoding.
//
// Useful when you don't want a handler down the chain to also handle encoding
func ClearEncoding(r *http.Request) {
	r.Header.Del(acceptEncoding)
}

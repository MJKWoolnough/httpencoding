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
	partSplit        = ';'
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
	encoding Encoding
	weight   uint16
}

// Encoding represents an encoding string as used by the client. Examples are
// gzip, br and deflate.
type Encoding string

// Handler provides an interface to handle an encoding.
//
// The encoding string (e.g. gzip, br, deflate) is passed to the handler, which
// is expected to return true if no more encodings are required and false
// otherwise.
//
// The empty string "" is used to signify the identity encoding, or plain text.
type Handler interface {
	Handle(encoding Encoding) bool
}

// HandlerFunc wraps a func to make it satisfy the Handler interface.
type HandlerFunc func(Encoding) bool

// Handle calls the underlying func.
func (h HandlerFunc) Handle(e Encoding) bool {
	return h(e)
}

// InvalidEncoding writes the 406 header.
func InvalidEncoding(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotAcceptable)
}

// HandleEncoding will process the Accept-Encoding header and calls the given
// handler for each encoding until the handler returns true.
//
// This function returns true when the Handler returns true, false otherwise.
//
// For the identity (plain text) encoding the encoding string will be the
// empty string.
//
// The wildcard encoding (*) is currently treated as identity when there is no
// independent identity encoding specified; otherwise, it is ignored.
func HandleEncoding(r *http.Request, h Handler) bool {
	acceptHeader := strings.TrimSpace(r.Header.Get(acceptEncoding))

	if len(acceptHeader) == 0 {
		acceptHeader = anyEncoding
	}

	accepts := make(encodings, 0, strings.Count(acceptHeader, acceptSplit)+2)
	hasIdentity := false
	hasNoAny := false
	anyPos := -1

	var nots strings.Builder

	nots.WriteString("*;")

Loop:
	for _, accept := range strings.Split(acceptHeader, acceptSplit) {
		hasQ := true
		split := strings.IndexByte(accept, partSplit)
		if split == -1 {
			split = len(accept)
			hasQ = false
		}

		name := strings.ToLower(strings.TrimSpace(accept[:split]))
		if name == "" {
			continue
		}

		var (
			qVal float64 = 1
			err  error
		)

		if hasQ {
			if part := strings.TrimSpace(accept[split+1:]); strings.HasPrefix(part, weightPrefix) {
				qVal, err = strconv.ParseFloat(part[len(weightPrefix):], 32)
				if err != nil || qVal < 0 || qVal > 1 {
					continue Loop
				}
			}
		}

		weight := uint16(qVal * 1000)

		if name == identityEncoding {
			hasIdentity = true
			name = ""
		} else if name == anyEncoding && qVal == 0 {
			hasNoAny = true
		}

		if weight == 0 {
			nots.WriteString(name)
			nots.WriteByte(';')
		} else {
			if name == anyEncoding {
				if anyPos != -1 {
					continue
				}

				anyPos = len(accepts)
			}

			accepts = append(accepts, encoding{
				encoding: Encoding(name),
				weight:   weight,
			})
		}
	}

	if anyPos != -1 {
		accepts[anyPos].encoding = Encoding(nots.String())
	}

	sort.Stable(accepts)

	if !hasIdentity && !hasNoAny {
		accepts = append(accepts, encoding{
			encoding: "",
			weight:   1,
		})
	}

	for _, accept := range accepts {
		if h.Handle(accept.encoding) {
			return true
		}
	}

	return false
}

// ClearEncoding removes the Accept-Encoding header so that any further
// attempts to establish an encoding will simply used the default, plain text,
// encoding.
//
// Useful when you don't want a handler down the chain to also handle encoding.
func ClearEncoding(r *http.Request) {
	r.Header.Del(acceptEncoding)
}

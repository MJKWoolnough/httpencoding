package httpencoding

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

type Handler interface {
	Handle(w http.ResponseWriter, r *http.Request, encoding string) bool
}

type HandlerFunc func(http.ResponseWriter, *http.Request, string) bool

func (h HandlerFunc) Handle(w http.ResponseWriter, r *http.Request, e string) bool {
	return h(w, r, e)
}

func InvalidEncoding(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotAcceptable)
}

func HandleEncoding(w http.ResponseWriter, r *http.Request, h Handler) bool {
	acceptHeader := r.Header.Get(acceptEncoding)
	accepts := make(encodings, 0, strings.Count(acceptHeader, acceptSplit)+1)
Loop:
	for _, accept := range strings.Split(acceptHeader, acceptSplit) {
		parts := strings.Split(strings.TrimSpace(accept), partSplit)
		name := strings.TrimSpace(parts[0])
		if name == "" {
			continue
		}
		var (
			weight float64 = 1
			err    error
		)
		for _, part := range parts[1:] {
			if strings.HasPrefix(strings.TrimSpace(part), weightPrefix) {
				weight, err = strconv.ParseFloat(part[len(weightPrefix):], 32)
				if err != nil || weight < 0 || weight >= 2 { // return an malformed header response?
					continue Loop
				}
				break
			}
		}
		accepts = append(accepts, encoding{
			encoding: strings.ToLower(name),
			weight:   uint16(weight * 1000),
		})
	}
	sort.Sort(accepts)
	allowIdentity := true
	for _, accept := range accepts {
		switch accept.encoding {
		case identityEncoding:
			allowIdentity = accept.weight != 0
			break
		case anyEncoding:
			allowIdentity = accept.weight != 0
		default:
			if h.Handle(w, r, accept.encoding) {
				return true
			}
		}
	}
	if allowIdentity {
		h.Handle(w, r, "")
		return true
	}
	return false
}

# httpencoding
--
    import "vimagination.zapto.org/httpencoding"

Package httpencoding provides a function to deal with the Accept-Encoding
header.

## Usage

#### func  ClearEncoding

```go
func ClearEncoding(r *http.Request)
```
ClearEncoding removes the Accept-Encoding header so that any further attempts to
establish an encoding will simply used the default, plain text, encoding.

Useful when you don't want a handler down the chain to also handle encoding.

#### func  HandleEncoding

```go
func HandleEncoding(r *http.Request, h Handler) bool
```
HandleEncoding will process the Accept-Encoding header and calls the given
handler for each encoding until the handler returns true.

This function returns true when the Handler returns true, false otherwise.

For the identity (plain text) encoding the encoding string will be the empty
string.

The wildcard encoding (*) will, after the '*', contain a semi-colon seperated
list of all disallowed encodings (q=0).

#### func  InvalidEncoding

```go
func InvalidEncoding(w http.ResponseWriter)
```
InvalidEncoding writes the 406 header.

#### func  IsDisallowedInWildcard

```go
func IsDisallowedInWildcard(accept, encoding string) bool
```
IsDisallowedInWildcard will return true if the given encoding is disallowed in
the given accept string.

#### func  IsWildcard

```go
func IsWildcard(accept string) bool
```
IsWildcard returns true when the given accept string is a wildcard match.

#### type Encoding

```go
type Encoding string
```

Encoding represents an encoding string as used by the client. Examples are gzip,
br and deflate.

#### type Handler

```go
type Handler interface {
	Handle(encoding Encoding) bool
}
```

Handler provides an interface to handle an encoding.

The encoding string (e.g. gzip, br, deflate) is passed to the handler, which is
expected to return true if no more encodings are required and false otherwise.

The empty string "" is used to signify the identity encoding, or plain text.

#### type HandlerFunc

```go
type HandlerFunc func(Encoding) bool
```

HandlerFunc wraps a func to make it satisfy the Handler interface.

#### func (HandlerFunc) Handle

```go
func (h HandlerFunc) Handle(e Encoding) bool
```
Handle calls the underlying func.

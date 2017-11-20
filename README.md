# httpencoding
--
    import "github.com/MJKWoolnough/httpencoding"

Package httpencoding provides a function to deal with the Accept-Encoding
header.

## Usage

#### func  HandleEncoding

```go
func HandleEncoding(r *http.Request, h Handler) bool
```
HandleEncoding will process the Accept-Encoding header and calls the given
handler for each encoding until the handler returns true.

This function returns true when the Handler returns true, false otherwise

For the identity (plain text) encoding the encoding string will be the empty
string.

The wildcard encoding (*) is currently treated as identity when there is no
independent identity encoding specified; otherwise, it is ignored.

#### func  InvalidEncoding

```go
func InvalidEncoding(w http.ResponseWriter)
```
InvalidEncoding writes the 406 header

#### type Handler

```go
type Handler interface {
	Handle(encoding string) bool
}
```

Handler provides an interface to handle an encoding

#### type HandlerFunc

```go
type HandlerFunc func(string) bool
```

HandlerFunc wraps a func to make it satisfy the Handler interface

#### func (HandlerFunc) Handle

```go
func (h HandlerFunc) Handle(e string) bool
```
Handle calls the underlying func

package httpencoding

import "vimagination.zapto.org/parser"

const (
	alpha       = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	digit       = "0123456789"
	tchar       = "!#$%&'*+-.^_`|~" + digit + alpha
	ows         = " \t"
	delim       = ","
	weightDelim = ";"
	qvalPrefix  = "q="
)

const (
	tokenCoding parser.TokenType = iota
	tokenWeight
	tokenInvalidWeight
)

func parseAccept(accept string) *parser.Parser {
	t := parser.New(parser.NewStringTokeniser(accept))

	t.TokeniserState(parseList)

	return &t
}

func parseList(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	for {
		t.AcceptRun(ows)
		t.Get()

		if t.Accept(tchar) {
			t.AcceptRun(tchar)

			return t.Return(tokenCoding, parseQVal)
		} else if t.Peek() == -1 {
			return t.Done()
		}

		t.ExceptRun(delim)
		t.Accept(delim)
	}
}

func parseQVal(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	t.AcceptRun(ows)
	t.Get()

	if !t.Accept(weightDelim) {
		t.ExceptRun(delim)
		t.Accept(delim)

		return parseList(t)
	}

	t.AcceptRun(ows)
	t.Get()

	if t.AcceptString(qvalPrefix, false) == len(qvalPrefix) {
		t.Get()

		if t.Accept("0") {
			return parse0(t)
		} else if t.AcceptString("1.000", false) > 0 {
			return parseWeightEnd(t)
		}
	}

	t.ExceptRun(delim)
	t.Accept(delim)

	return t.Return(tokenInvalidWeight, parseList)
}

func parse0(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	if t.Accept(".") {
		t.Accept(digit)
		t.Accept(digit)
		t.Accept(digit)
	}

	return parseWeightEnd(t)
}

func parseWeightEnd(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	data := t.Get()

	t.AcceptRun(ows)

	if t.Accept(delim) || t.Peek() == -1 {
		return parser.Token{Type: tokenWeight, Data: data}, parseList
	}

	t.ExceptRun(delim)
	t.Accept(delim)
	t.Get()

	return t.Return(tokenInvalidWeight, parseList)
}

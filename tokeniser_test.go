package httpencoding

import (
	"testing"

	"vimagination.zapto.org/parser"
)

func TestTokeniser(t *testing.T) {
	for n, test := range [...]struct {
		Accept string
		Tokens []parser.Token
	}{
		{
			"",
			[]parser.Token{
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{
			"a",
			[]parser.Token{
				{Type: tokenCoding, Data: "a"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{
			"a,b",
			[]parser.Token{
				{Type: tokenCoding, Data: "a"},
				{Type: tokenCoding, Data: "b"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{
			"a\t , \tb",
			[]parser.Token{
				{Type: tokenCoding, Data: "a"},
				{Type: tokenCoding, Data: "b"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{
			"a;q=1",
			[]parser.Token{
				{Type: tokenCoding, Data: "a"},
				{Type: tokenWeight, Data: "1"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{
			"a ; q=1",
			[]parser.Token{
				{Type: tokenCoding, Data: "a"},
				{Type: tokenWeight, Data: "1"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{
			"a ; q=1.",
			[]parser.Token{
				{Type: tokenCoding, Data: "a"},
				{Type: tokenWeight, Data: "1."},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{
			"a ; q=1.0",
			[]parser.Token{
				{Type: tokenCoding, Data: "a"},
				{Type: tokenWeight, Data: "1.0"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{
			"a ; q=1.000",
			[]parser.Token{
				{Type: tokenCoding, Data: "a"},
				{Type: tokenWeight, Data: "1.000"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{
			"a ; q=1.0000",
			[]parser.Token{
				{Type: tokenCoding, Data: "a"},
				{Type: tokenInvalidWeight, Data: ""},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{
			"a ; q=1.100",
			[]parser.Token{
				{Type: tokenCoding, Data: "a"},
				{Type: tokenInvalidWeight, Data: ""},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{
			"a ; q=2",
			[]parser.Token{
				{Type: tokenCoding, Data: "a"},
				{Type: tokenInvalidWeight, Data: "2"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{
			"a ; q=2",
			[]parser.Token{
				{Type: tokenCoding, Data: "a"},
				{Type: tokenInvalidWeight, Data: "2"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{
			"a ; q=0",
			[]parser.Token{
				{Type: tokenCoding, Data: "a"},
				{Type: tokenWeight, Data: "0"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{
			"a ; q=0.123",
			[]parser.Token{
				{Type: tokenCoding, Data: "a"},
				{Type: tokenWeight, Data: "0.123"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{
			"a ; q=0.1234",
			[]parser.Token{
				{Type: tokenCoding, Data: "a"},
				{Type: tokenInvalidWeight, Data: ""},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{
			"a ; q=0.a",
			[]parser.Token{
				{Type: tokenCoding, Data: "a"},
				{Type: tokenInvalidWeight, Data: ""},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{
			"*;q=0,identity;q=0.512",
			[]parser.Token{
				{Type: tokenCoding, Data: "*"},
				{Type: tokenWeight, Data: "0"},
				{Type: tokenCoding, Data: "identity"},
				{Type: tokenWeight, Data: "0.512"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
	} {
		p := parseAccept(test.Accept)

		for m, tkn := range test.Tokens {
			if tk, _ := p.GetToken(); tk.Type != tkn.Type {
				if tk.Type == parser.TokenError {
					t.Errorf("test %d.%d: unexpected error: %s", n+1, m+1, tk.Data)
				} else {
					t.Errorf("test %d.%d: Incorrect type, expecting %d, got %d", n+1, m+1, tkn.Type, tk.Type)
				}

				break
			} else if tk.Data != tkn.Data {
				t.Errorf("test %d.%d: Incorrect data, expecting %q, got %q", n+1, m+1, tkn.Data, tk.Data)

				break
			}
		}
	}
}

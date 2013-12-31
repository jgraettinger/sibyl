package parser

import (
	"io"
)

type Token string

type TokenStream interface {
	// Produces the next token in a token sequence. io.EOF is returned
	// to signal a graceful end of input (eg, end of the sentence).
	NextToken() (Token, error)
}

type TokenArray []Token

func (a *TokenArray) NextToken() (Token, error) {
	if len(*a) != 0 {
		token := (*a)[0]
		*a = (*a)[1:len(*a)]
		return token, nil
	}
	return "", io.EOF
}

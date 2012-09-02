package chart

type Token string

func (token Token) IsStoppingPunctuaction() bool {

	switch token {
	case ",", ".", ";", ":", "--":
		return true
	}
	return false
}

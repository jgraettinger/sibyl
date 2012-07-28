package chart

func IsStoppingPunctuaction(token string) bool {

    switch token {
    case ",", ".", ";", ":", "--":
        return true
    }
    return false
}


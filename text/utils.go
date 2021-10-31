package text

import "unicode"

func grabLine(t []rune) []rune {
	for i, r := range t {
		if r == '\n' {
			return t[:i]
		}
	}
	return t
}

func discardLastWord(t []rune) []rune {
	w := false
	for i := len(t) - 1; i >= 0; i-- {
		if unicode.IsSpace(t[i]) && !w {
			w = true
		}
		if !unicode.IsSpace(t[i]) && w {
			return t[:i+1]
		}
	}
	// Discard char by char
	return t[:len(t)-1]
}

func trnmSpaces(t []rune) []rune {
	return trimFunc(t, unicode.IsSpace)
}

func trimFunc(t []rune, fn func(rune) bool) []rune {
	res := t
	for i := 0; i < len(res); i++ {
		if !fn(res[i]) {
			res = res[i:]
			break
		}
	}
	for i := len(res) - 1; i >= 0; i-- {
		if !fn(res[i]) {
			res = res[:i+1]
			break
		}
	}
	return res
}

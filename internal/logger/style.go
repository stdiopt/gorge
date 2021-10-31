package logger

import "fmt"

// Style options for specific style
type Style struct {
	Prefix       string
	Suffix       string
	IncrementPad bool
	Pad          int
}

// Get styled message
func (s *Style) Get(msg interface{}) string {
	m := fmt.Sprint(msg)

	if s.IncrementPad {
		if t := len(m); t > s.Pad {
			s.Pad = t
		}
		m = fmt.Sprintf("%*s", s.Pad, m)
	}

	return fmt.Sprintf("%s%s%s", s.Prefix, m, s.Suffix)
}

// GetCustom uses the custom prefix and suffix for the message, it will pad if
// necessary.
func (s *Style) GetCustom(pre, suf string, msg interface{}) string {
	m := fmt.Sprint(msg)

	if s.IncrementPad {
		if t := len(m); t > s.Pad {
			s.Pad = t
		}
		m = fmt.Sprintf("%*s", s.Pad, m)
	}

	return fmt.Sprintf("%s%s%s", pre, m, suf)
}

package utils

import "strings"

func EscapeString(s string) string {
	var res strings.Builder
	for _, r := range s {
		res.WriteString(EscapeRune(r))
	}
	return res.String()
}

func EscapeRune(r rune) string {
	switch r {
	case '\'':
		return "\\'"
	case '"':
		return "\\\""
	case '\\':
		return "\\\\"
	case '\n':
		return "\\n"
	case '\r':
		return "\\r"
	case '\t':
		return "\\t"
	case '\b':
		return "\\b"
	case '\f':
		return "\\f"
	}
	return string(r)
}

func UnescapeRune(s string) rune {
	switch s {
	case "\\'":
		return '\''
	case "\\\"":
		return '"'
	case "\\\\":
		return '\\'
	case "\\n":
		return '\n'
	case "\\r":
		return '\r'
	case "\\t":
		return '\t'
	case "\\b":
		return '\b'
	case "\\f":
		return '\f'
	}
	return rune(s[0])
}

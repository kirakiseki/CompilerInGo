package utils

import "strings"

// EscapeString 将字符串中的特殊字符转义
func EscapeString(s string) string {
	// 使用strings.Builder拼接字符串
	var res strings.Builder

	for _, r := range s {
		// 遍历字符串中的每一个字符并转义
		res.WriteString(EscapeRune(r))
	}

	// 返回转义后的字符串
	return res.String()
}

// EscapeRune 将字符转义
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

	// 不需转义字符，使用字符串返回
	return string(r)
}

// UnescapeRune 将单个转义字符组成的字符串还原为字符
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

	// 不需还原字符，返回第一个字符
	return rune(s[0])
}

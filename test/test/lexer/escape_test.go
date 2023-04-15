package lexer

import (
	"CompilerInGo/utils"
	"testing"
)

func TestEscapeRune(t *testing.T) {
	// Right Case
	var rightCase = map[rune]string{
		'\'': "\\'",
		'"':  "\\\"",
		'\\': "\\\\",
		'\n': "\\n",
		'\r': "\\r",
		'\t': "\\t",
		'\b': "\\b",
		'\f': "\\f",

		'a': "a",
		'1': "1",
		' ': " ",
		'!': "!",
		'@': "@",
		'#': "#",
		'$': "$",
	}

	for k, v := range rightCase {
		if res := utils.EscapeRune(k); res != v {
			t.Error("EscapeRune failed")
			t.Error("Input: ", k)
			t.Error("Expected: ", v)
			t.Error("Actual: ", res)
		}
	}
}

func TestEscapeString(t *testing.T) {
	// Right Case
	var rightCase = map[string]string{
		"abc":        "abc",
		"123":        "123",
		" ":          " ",
		"!":          "!",
		"\n":         "\\n",
		"abc\n123":   "abc\\n123",
		"\\":         "\\\\",
		"\"":         "\\\"",
		"'":          "\\'",
		"\"'":        "\\\"\\'",
		"a\b\n\\'\"": "a\\b\\n\\\\\\'\\\"",
	}

	for k, v := range rightCase {
		if res := utils.EscapeString(k); res != v {
			t.Error("EscapeString failed")
			t.Error("Input: ", k)
			t.Error("Expected: ", v)
			t.Error("Actual: ", res)
		}
	}
}

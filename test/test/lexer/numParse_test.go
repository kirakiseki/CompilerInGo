package lexer

import (
	"CompilerInGo/utils"
	"testing"
)

func TestParseInt(t *testing.T) {
	// Right Case
	var rightCase = map[string]int64{
		"123":         123,
		"0":           0,
		"-1":          -1,
		"1":           1,
		"1234567890":  1234567890,
		"-1234567890": -1234567890,
	}

	for k, v := range rightCase {
		if num, err := utils.ParseInt(k); num != v || err != nil {
			t.Error("ParseInt failed")
			t.Error("Input: ", k)
			t.Error("Expected: ", v)
			t.Error("Actual: ", num)
			t.Error("Error: ", err)
		}
	}

	// Wrong Case
	var wrongCase = []string{
		"",
		"-",
		"abc",
		"123abc",
		"abc123",
		"123.456",
		"123.456.789",
		"-123.456",
		"0x123",
		"0b123",
	}

	for _, v := range wrongCase {
		if _, err := utils.ParseInt(v); err == nil {
			t.Error("ParseInt failed")
			t.Error("Input: ", v)
			t.Error("Expected: ", "error")
			t.Error("Actual: ", "nil")
		}
	}
}

func TestParseFloat(t *testing.T) {
	// Right Case
	var rightCase = map[string]float64{
		"123":   123,
		"0":     0,
		"-1":    -1,
		"1.0":   1,
		"0.0":   0,
		"1.23":  1.23,
		"-0.01": -0.01,
		"-0.0":  0,
		"123.":  123,
		"9.":    9,
	}

	for k, v := range rightCase {
		if num, err := utils.ParseFloat(k); num != v || err != nil {
			t.Error("ParseFloat failed")
			t.Error("Input: ", k)
			t.Error("Expected: ", v)
			t.Error("Actual: ", num)
			t.Error("Error: ", err)
		}
	}

	// Wrong Case
	var wrongCase = []string{
		"",
		"-",
		"abc",
		"123abc",
		"abc123",
		"123.456.789",
		"0x123",
		"0b123",
		".123",
		"78.@",
	}

	for _, v := range wrongCase {
		if num, err := utils.ParseFloat(v); err == nil {
			t.Error("ParseFloat failed")
			t.Error("Input: ", v)
			t.Error("Expected: ", "error")
			t.Error("Actual: ", "nil")
			t.Error("Actual Result: ", num)
		}
	}
}

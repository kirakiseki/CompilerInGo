package utils

import (
	"errors"
)

// ParseInt 将字符串转换为int64
func ParseInt(s string) (int64, error) {
	// 空字符串
	if s == "" {
		return 0, errors.New("failed to ParseInt: empty string")
	}

	// 负数标志
	neg := false

	// 判断是否为负数
	if s[0] == '-' {
		// 只含负号
		if len(s) == 1 {
			return 0, errors.New("failed to ParseInt: invalid number string")
		}
		// 设置负数标志
		neg = true
		s = s[1:]
	}

	// 结果
	var n int64

	// 遍历数字部分
	for _, ch := range s {
		if ch < '0' || ch > '9' {
			return 0, errors.New("failed to ParseInt: invalid char")
		}
		n = n*10 + int64(ch-'0')
	}

	if neg {
		// 返回负数
		return -n, nil
	}
	return n, nil
}

// ParseFloat 将字符串转换为float64
func ParseFloat(s string) (float64, error) {
	// 空字符串
	if s == "" {
		return 0, errors.New("failed to ParseFloat: empty string")
	}

	// 首字符为小数点
	if s[0] == '.' {
		return 0, errors.New("failed to ParseFloat: invalid number string")
	}

	// 负数标志
	neg := false

	// 判断是否为负数
	if s[0] == '-' {
		// 只含负号
		if len(s) == 1 {
			return 0, errors.New("failed to ParseFloat: invalid number string")
		}
		neg = true
		s = s[1:]
	}

	// 整数部分
	var integ float64

	// 整数长度i
	var i int

	for i = 0; i < len(s); i++ {
		// 遇到小数点
		if s[i] == '.' {
			break
		}
		// 非数字，非法
		if s[i] < '0' || s[i] > '9' {
			return 0, errors.New("failed to ParseFloat: invalid char")
		}

		// 拼接整数部分
		integ = integ*10 + float64(s[i]-'0')
	}

	// 无小数部分
	if i == len(s) {
		if neg {
			return -integ, nil
		}
		return integ, nil
	}

	// 小数部分
	var frac float64
	for j := len(s) - 1; j > i; j-- {
		if s[j] < '0' || s[j] > '9' {
			return 0, errors.New("failed to ParseFloat: invalid char")
		}
		frac = frac/10 + float64(s[j]-'0')
	}

	if neg {
		return -(integ + frac/10), nil
	}
	return integ + frac/10, nil
}

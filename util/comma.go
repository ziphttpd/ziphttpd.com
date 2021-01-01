package util

import (
	"strconv"
	"strings"
)

func slice(str string, slen int) []string {
	ret := []string{}
	ls := len(str)
	start := 0
	for start < ls {
		end := start + slen
		if ls < end {
			end = ls
		}
		ret = append(ret, str[start:end])
		start += slen
	}
	return ret
}
func splitEnd(str string, size int) []string {
	ret := []string{}
	maxlen := len(str)
	start := maxlen % size
	if start != 0 {
		ret = append(ret, str[:start])
	}
	if start < maxlen {
		ret = append(ret, slice(str[start:], size)...)
	}
	return ret
}

// Comma は数字をカンマ区切り文字列に変換します。
func Comma(num int) string {
	return strings.Join(splitEnd(strconv.Itoa(num), 3), ",")
}

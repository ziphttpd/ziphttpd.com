package util

// Strcat は文字列を切り詰めます。
func Strcat(str string, size int) string {
	rn := []rune(str)
	len := len(rn)
	if len > size {
		len = size
	}
	return string(rn[:len])
}

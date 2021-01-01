package util

import "time"

const ymdFMT = "20060102150405" // 日付書式（golangの決まり）

// Time2ymdhms はTimeをyyyymmddhhmmss形式文字列に変換します。
func Time2ymdhms(t time.Time) string {
	return t.Format(ymdFMT)
}

// Ymdhms2time はyyyymmddhhmmss形式文字列をTimeに変換します。
func Ymdhms2time(ymdhms string) (time.Time, error) {
	return time.Parse(ymdFMT, ymdhms)
}

// Ymdhms2GMT はyyyymmddhhmmss形式文字列をHTTPでのGMT表記に変換します。
func Ymdhms2GMT(ymdhms string) (string, error) {
	utc, err := Ymdhms2time(ymdhms)
	if err != nil {
		return "", err
	}
	headtime := utc.Format("Mon, 02 Jan 2006 15:04:05 GMT")
	return headtime, nil
}

package tools

import "strconv"

func ByteToInt64(str []byte) int64 {
	num := StrToInt64(string(str))
	return num
}
func StrToInt64(str string) int64 {
	num, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0
	}
	return num
}
func Int64ToStr(i int64) string {
	str := strconv.FormatInt(i, 10)
	return str
}
func Int64ToByte(i int64) []byte {
	str := []byte(Int64ToStr(i))
	return str
}

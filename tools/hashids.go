package tools

import (
	"github.com/speps/go-hashids"
)

func HashIds(salt string, id int64) (string, error) {
	hd := hashids.NewData()
	//id 可用字符
	hd.Alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890"
	//盐
	hd.Salt = salt
	//长度
	hd.MinLength = 10

	h, err := hashids.NewWithData(hd)
	if err != nil {
		return "", err
	}
	e, err := h.EncodeInt64([]int64{id})
	if err != nil {
		return "", err
	}
	return e, nil
}
func HashIdsDecode(salt, str string) ([]int64, error) {
	hd := hashids.NewData()
	hd.Alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890"
	//盐
	hd.Salt = salt
	//长度
	hd.MinLength = 10
	h, err := hashids.NewWithData(hd)
	if err != nil {
		return []int64{0}, err
	}
	e, err := h.DecodeInt64WithError(str)
	if err != nil {
		return []int64{0}, err
	}
	return e, nil
}

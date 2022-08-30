package models

import (
	"fmt"
	"go-fly-muti/tools"
	"testing"
)

func TestUpdateUser(t *testing.T) {
	md5Pass := tools.Md5("caonima!@#A")

	fmt.Println(md5Pass)
}

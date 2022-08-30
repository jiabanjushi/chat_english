package types

import "testing"

func TestApiCode(t *testing.T) {
	t.Log(ApiCode.GetMessage(ApiCode.SUCCESS))
}

package types

import (
	"testing"
	"time"
)

func TestTime_MarshalJSON(t *testing.T) {
	now := time.Now()
	ti := &Time{
		Time: now,
	}
	res, _ := ti.MarshalJSON()
	t.Log(string(res))
}

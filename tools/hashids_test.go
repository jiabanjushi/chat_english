package tools

import (
	"testing"
	"time"
)

func TestHashIds(t *testing.T) {
	salt := "taoshihan"
	id := time.Now().Unix()
	hash, err := HashIds(salt, id)
	t.Log(hash, err)
	ids, err := HashIdsDecode(salt, hash)
	t.Log(ids, err)
}

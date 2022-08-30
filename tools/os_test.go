package tools

import (
	"log"
	"os"
	"testing"
)

func TestOs(t *testing.T) {
	log.Println(os.Getwd())
	log.Println(os.Executable())
}

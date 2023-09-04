package junit

import (
	_ "embed"
	"testing"
)

//go:embed data/junit.log
var raw []byte

func TestOne(t *testing.T) {
	defer t.Fail()

	Ingest("idx", raw)
}

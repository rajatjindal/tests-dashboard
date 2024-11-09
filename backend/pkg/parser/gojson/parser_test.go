package gojson

import (
	_ "embed"
	"testing"

	"github.com/stretchr/testify/assert"
)

//go:embed data/logs.json.txt
var data []byte

func TestOne(t *testing.T) {
	defer t.Fail()

	_, _, err := Ingest("idx", data)
	assert.Nil(t, err)
}

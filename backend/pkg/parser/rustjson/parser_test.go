package rustjson

import (
	_ "embed"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

//go:embed data/logs.json.txt
var data []byte

func TestOne(t *testing.T) {
	defer t.Fail()

	summary, suites, err := Ingest("idx", data)
	assert.Nil(t, err)

	fmt.Printf("%#v\n\n", summary)

	for _, suite := range suites {
		for _, test := range suite.TestsTree {
			fmt.Printf("%#v\n\n", test)
		}
	}

}

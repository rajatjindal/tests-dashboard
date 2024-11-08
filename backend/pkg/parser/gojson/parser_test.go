package gojson

import (
	"context"
	_ "embed"
	"testing"

	"github.com/rajatjindal/tests-dashboard/backend/pkg/storage"
	"github.com/rajatjindal/tests-dashboard/backend/pkg/types"
	"github.com/stretchr/testify/assert"
)

//go:embed data/logs.json.txt
var data []byte

func TestOne(t *testing.T) {
	defer t.Fail()

	summary, suites, err := Ingest("idx", data)
	assert.Nil(t, err)

	metadata := types.Metadata{
		RunId: "unique",
	}
	err = storage.IngestTestRun(context.Background(), &metadata, summary, suites)
	if err != nil {
	}
}

package rustjson

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/rajatjindal/test-dashboard/backend/pkg/types"
)

func Ingest(runId string, data []byte) (*types.Summary, []types.Suite, error) {
	p := &eventProcessor{
		Suites: []types.Suite{},
		Summary: &types.Summary{
			RunId:     runId,
			CreatedAt: time.Now().Format(time.RFC3339),
		},
	}

	scanner := bufio.NewScanner(bytes.NewReader(data))
	scanner.Split(bufio.ScanLines)

	index := 0
	for scanner.Scan() {
		index++
		var event types.Event

		text := scanner.Text()
		err := json.Unmarshal([]byte(text), &event)
		if err != nil {
			fmt.Println(index, " ", text)
			fmt.Printf("ignore error index: %d, txt: %s, err: %v", index, text, err)
			continue
		}

		p.processEvent(runId, event)
	}

	return p.Summary, p.Suites, nil
}

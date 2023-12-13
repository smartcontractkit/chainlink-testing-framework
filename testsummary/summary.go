package testsummary

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-testing-framework/utils/runid"
)

var (
	SUMMARY_FOLDER = ".test_summary"
	SUMMARY_FILE   string
	mu             sync.Mutex
)

type SummaryKeys map[string][]KeyContent

type KeyContent struct {
	TestName string `json:"test_name"`
	Value    string `json:"value"`
}

func init() {
	runId, err := runid.GetOrGenerateRunId(nil)
	if err != nil {
		panic(err)
	}
	SUMMARY_FILE = fmt.Sprintf("%s/test_summary-%s-%s.json", SUMMARY_FOLDER, time.Now().Format("2006-01-02T15-04-05"), runId)
}

// TODO in future allow value to be also []string or map[string]string?
func AddEntry(testName, key string, value interface{}) error {
	mu.Lock()
	defer mu.Unlock()

	if _, ok := value.(string); !ok {
		return errors.Errorf("type '%T' not supported", value)
	}
	strValue := value.(string)

	if err := os.MkdirAll(SUMMARY_FOLDER, 0755); err != nil {
		return err
	}

	f, err := os.OpenFile(SUMMARY_FILE, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	fc, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	var entries SummaryKeys
	err = json.Unmarshal(fc, &entries)
	if err != nil {
		if !strings.Contains(err.Error(), "unexpected end of JSON input") {
			return err
		}

		entries = make(SummaryKeys)
	}

	if entry, ok := entries[key]; ok {
		testFound := false
		for idx, testValue := range entry {
			// overwrite if entry for test exists
			if testValue.TestName == testName {
				entry[idx].Value = strValue
				testFound = true
				break
			}
		}

		// add new entry to existing key if no entry for test exists
		if !testFound {
			entries[key] = append(entries[key], KeyContent{TestName: testName, Value: strValue})
		}
	} else {
		entries[key] = []KeyContent{{TestName: testName, Value: strValue}}
	}

	_, err = f.Seek(0, 0)
	if err != nil {
		return err
	}

	err = f.Truncate(0)
	if err != nil {
		return err
	}

	encoder := json.NewEncoder(f)
	err = encoder.Encode(entries)
	if err != nil {
		return err
	}

	return nil
}

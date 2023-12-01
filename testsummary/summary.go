package testsummary

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/pkg/errors"
)

var (
	SUMMARY_FOLDER = ".test_summary"
	SUMMARY_FILE   = fmt.Sprintf("%s/test_summary.json", SUMMARY_FOLDER)
)

type SummaryKeys map[string][]KeyContent

type KeyContent struct {
	TestName string `json:"test_name"`
	Value    string `json:"value"`
}

var mu sync.Mutex

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
		return err
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

	encoder := json.NewEncoder(f)
	err = encoder.Encode(entries)
	if err != nil {
		return err
	}

	return nil
}

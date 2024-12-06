package testsummary

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	SUMMARY_FOLDER = ".test_summary"
	mu             sync.Mutex
)

var SUMMARY_FILE = fmt.Sprintf("%s/test_summary-%s.json", SUMMARY_FOLDER, time.Now().Format("2006-01-02T15-04-05"))

type SummaryKeys map[string][]KeyContent

type KeyContent struct {
	TestName string `json:"test_name"`
	Value    string `json:"value"`
}

// TODO in future allow value to be also []string or map[string]string?
func AddEntry(testName, key string, value interface{}) error {
	mu.Lock()
	defer mu.Unlock()

	if _, ok := value.(string); !ok {
		return fmt.Errorf("type '%T' not supported", value)
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

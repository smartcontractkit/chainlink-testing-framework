package benchspy

import (
	"encoding/json"
	"fmt"
)

// StandardReport is a report that contains all the necessary data for a performance test
type StandardReport struct {
	BasicData
	LocalReportStorage
	ResourceReporter
	QueryExecutors []QueryExecutor `json:"query_executors"`
}

func (b *StandardReport) Store() (string, error) {
	return b.LocalReportStorage.Store(b.TestName, b.CommitOrTag, b)
}

func (b *StandardReport) Load() error {
	return b.LocalReportStorage.Load(b.TestName, b.CommitOrTag, b)
}

func (b *StandardReport) Fetch() error {
	basicErr := b.BasicData.Validate()
	if basicErr != nil {
		return basicErr
	}

	// TODO parallelize it
	for _, queryExecutor := range b.QueryExecutors {
		queryExecutor.TimeRange(b.TestStart, b.TestEnd)

		if validateErr := queryExecutor.Validate(); validateErr != nil {
			return validateErr
		}

		if execErr := queryExecutor.Execute(); execErr != nil {
			return execErr
		}
	}

	resourceErr := b.FetchResources()
	if resourceErr != nil {
		return resourceErr
	}

	return nil
}

func (b *StandardReport) IsComparable(otherReport StandardReport) error {
	basicErr := b.BasicData.IsComparable(otherReport.BasicData)
	if basicErr != nil {
		return basicErr
	}

	if resourceErr := b.CompareResources(&otherReport.ResourceReporter); resourceErr != nil {
		return resourceErr
	}

	for i, queryExecutor := range b.QueryExecutors {
		queryErr := queryExecutor.IsComparable(otherReport.QueryExecutors[i])
		if queryErr != nil {
			return queryErr
		}
	}

	return nil
}

func (s *StandardReport) UnmarshalJSON(data []byte) error {
	// helper struct with QueryExecutors as json.RawMessage
	type Alias StandardReport
	var raw struct {
		Alias
		QueryExecutors []json.RawMessage `json:"query_executors"`
	}

	// unmarshal into the helper struct to populate other fields automatically
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	var queryExecutors []QueryExecutor

	// manually decide, which QueryExecutor implementation to use based on the "kind" field
	for _, rawExecutor := range raw.QueryExecutors {
		var typeIndicator struct {
			Kind string `json:"kind"`
		}
		if err := json.Unmarshal(rawExecutor, &typeIndicator); err != nil {
			return err
		}

		var executor QueryExecutor
		switch typeIndicator.Kind {
		case "loki":
			executor = &LokiQuery{}
		default:
			return fmt.Errorf("unknown query executor type: %s", typeIndicator.Kind)
		}

		if err := json.Unmarshal(rawExecutor, executor); err != nil {
			return err
		}

		queryExecutors = append(s.QueryExecutors, executor)
	}

	*s = StandardReport(raw.Alias)
	s.QueryExecutors = queryExecutors
	return nil
}

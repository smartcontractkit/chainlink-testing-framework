package benchspy

import (
	"context"
	"encoding/json"
	"fmt"

	"golang.org/x/sync/errgroup"
)

// StandardReport is a report that contains all the necessary data for a performance test
type StandardReport struct {
	BasicData
	LocalStorage
	ResourceReporter
	QueryExecutors []QueryExecutor `json:"query_executors"`
}

func (b *StandardReport) Store() (string, error) {
	return b.LocalStorage.Store(b.TestName, b.CommitOrTag, b)
}

func (b *StandardReport) Load() error {
	return b.LocalStorage.Load(b.TestName, b.CommitOrTag, b)
}

func (b *StandardReport) FetchData(ctx context.Context) error {
	startEndErr := b.BasicData.FillStartEndTimes()
	if startEndErr != nil {
		return startEndErr
	}

	basicErr := b.BasicData.Validate()
	if basicErr != nil {
		return basicErr
	}

	errGroup, errCtx := errgroup.WithContext(ctx)

	for _, queryExecutor := range b.QueryExecutors {
		errGroup.Go(func() error {
			queryExecutor.TimeRange(b.TestStart, b.TestEnd)

			if validateErr := queryExecutor.Validate(); validateErr != nil {
				return validateErr
			}

			if execErr := queryExecutor.Execute(errCtx); execErr != nil {
				return execErr
			}

			return nil
		})
	}

	if err := errGroup.Wait(); err != nil {
		return err
	}

	resourceErr := b.FetchResources(ctx)
	if resourceErr != nil {
		return resourceErr
	}

	return nil
}

func (b *StandardReport) IsComparable(otherReport Reporter) error {
	if _, ok := otherReport.(*StandardReport); !ok {
		return fmt.Errorf("expected type %s, got %T", "*StandardReport", otherReport)
	}

	asStandardReport := otherReport.(*StandardReport)

	basicErr := b.BasicData.IsComparable(asStandardReport.BasicData)
	if basicErr != nil {
		return basicErr
	}

	if resourceErr := b.CompareResources(&asStandardReport.ResourceReporter); resourceErr != nil {
		return resourceErr
	}

	for i, queryExecutor := range b.QueryExecutors {
		queryErr := queryExecutor.IsComparable(asStandardReport.QueryExecutors[i])
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
			executor = &LokiQueryExecutor{}
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

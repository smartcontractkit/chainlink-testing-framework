package comparator

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types/container"
	tc "github.com/testcontainers/testcontainers-go"
	"golang.org/x/sync/errgroup"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/client"
	"github.com/smartcontractkit/chainlink-testing-framework/wasp"
)

type ExecutionEnvironment string

const (
	ExecutionEnvironment_Docker                  ExecutionEnvironment = "docker"
	ExecutionEnvironment_k8sExecutionEnvironment                      = "k8s"
)

type Report interface {
	// Store stores the report in a persistent storage and returns the path to it, or an error
	Store() (string, error)
	// Load loads the report from a persistent storage and returns it, or an error
	Load() error
	// Fetch populates the report with the data from the test
	Fetch() error
	// IsComparable checks whether both reports can be compared (e.g. test config is the same, app's resources are the same, queries or metrics used are the same, etc.), and returns a map of the differences and an error (if any difference is found)
	IsComparable(otherReport Report) (bool, map[string]string, error)
	//// Compare compares two reports and returns a map of the differences and an error (if any difference is found). You should call it only after IsComparable returns true
	//Compare(other Report) (bool, map[string]string, error)
}

var directory = "performance_reports"

type BasicReport struct {
	TestName             string               `json:"test_name"`
	ExecutionEnvironment ExecutionEnvironment `json:"execution_environment"`

	// Test metrics
	CommitOrTag string    `json:"commit_or_tag"`
	TestStart   time.Time `json:"test_start_timestamp"`
	TestEnd     time.Time `json:"test_end_timestamp"`

	// all, generator settings, including segments
	GeneratorConfigs map[string]*wasp.Config `json:"generator_configs"`

	// AUT metrics
	PodsResources            map[string]*PodResources    `json:"pods_resources"`
	ContainerResources       map[string]*DockerResources `json:"container_resources"`
	ResourceSelectionPattern string                      `json:"resource_selection_pattern"`

	// Performance queries
	// a map of name to query template, ex: avg(rate(cpu_usage_seconds_total[5m])) @ 1663915200
	LokiQueries map[string]string `json:"loki_queries"`
	// Performance queries results
	// can be anything, avg RPS, amount of errors, 95th percentile of CPU utilization, etc
	Results map[string][]string `json:"results"`
	// In case something went wrong
	Errors []error `json:"errors"`

	LokiConfig *wasp.LokiConfig `json:"-"`
}

type DockerResources struct {
	NanoCPUs   int64
	CpuShares  int64
	Memory     int64
	MemorySwap int64
}

type PodResources struct {
	CPU    int64
	Memory int64
}

func (b *BasicReport) Store() (string, error) {
	asJson, err := json.MarshalIndent(b, "", " ")
	if err != nil {
		return "", err
	}

	if _, err := os.Stat(directory); os.IsNotExist(err) {
		if err := os.MkdirAll(directory, 0755); err != nil {
			return "", errors.Wrapf(err, "failed to create directory %s", directory)
		}
	}

	reportFilePath := filepath.Join(directory, fmt.Sprintf("%s-%s.json", b.TestName, b.CommitOrTag))
	reportFile, err := os.Create(reportFilePath)
	if err != nil {
		return "", errors.Wrapf(err, "failed to create file %s", reportFilePath)
	}
	defer func() { _ = reportFile.Close() }()

	reader := bytes.NewReader(asJson)
	_, err = io.Copy(reportFile, reader)
	if err != nil {
		return "", errors.Wrapf(err, "failed to write to file %s", reportFilePath)
	}

	abs, err := filepath.Abs(reportFilePath)
	if err != nil {
		return reportFilePath, nil
	}

	return abs, nil
}

func (b *BasicReport) Load() error {
	if b.TestName == "" {
		return errors.New("test name is empty. Please set it and try again")
	}

	var reportFilePath string
	if b.CommitOrTag == "" {
		// TODO load the latest one, but we need to think how to know which one is the latest one
		// with tags it is easy, just use semver, but with commits not so straightforward, although we could probably comb the repo to find the commit and then its commit date
		panic("not implemented")
	} else {
		reportFilePath = filepath.Join(directory, fmt.Sprintf("%s-%s.json", b.TestName, b.CommitOrTag))
	}

	reportFile, err := os.Open(reportFilePath)
	if err != nil {
		return errors.Wrapf(err, "failed to open file %s", reportFilePath)
	}

	decoder := json.NewDecoder(reportFile)
	if err := decoder.Decode(b); err != nil {
		return errors.Wrapf(err, "failed to decode file %s", reportFilePath)
	}

	return nil
}

func (b *BasicReport) Fetch() error {
	if len(b.LokiQueries) == 0 {
		return errors.New("there are no Loki queries, there's nothing to fetch. Please set them and try again")
	}
	if b.LokiConfig == nil {
		return errors.New("loki config is missing. Please set it and try again")
	}
	if b.TestStart.IsZero() {
		return errors.New("test start time is missing. We cannot query Loki without a time range. Please set it and try again")
	}
	if b.TestEnd.IsZero() {
		return errors.New("test end time is missing. We cannot query Loki without a time range. Please set it and try again")
	}

	splitAuth := strings.Split(b.LokiConfig.BasicAuth, ":")
	var basicAuth client.LokiBasicAuth
	if len(splitAuth) == 2 {
		basicAuth = client.LokiBasicAuth{
			Login:    splitAuth[0],
			Password: splitAuth[1],
		}
	}

	b.Results = make(map[string][]string)

	for name, query := range b.LokiQueries {
		queryParams := client.LokiQueryParams{
			Query:     query,
			StartTime: b.TestStart,
			EndTime:   b.TestEnd,
			Limit:     1000, //TODO make this configurable
		}

		parsedLokiUrl, err := url.Parse(b.LokiConfig.URL)
		if err != nil {
			return errors.Wrapf(err, "failed to parse Loki URL %s", b.LokiConfig.URL)
		}

		lokiUrl := parsedLokiUrl.Scheme + "://" + parsedLokiUrl.Host
		lokiClient := client.NewLokiClient(lokiUrl, b.LokiConfig.TenantID, basicAuth, queryParams)

		ctx, cancelFn := context.WithTimeout(context.Background(), b.LokiConfig.Timeout)
		rawLogs, err := lokiClient.QueryLogs(ctx)
		if err != nil {
			b.Errors = append(b.Errors, err)
			cancelFn()
			continue
		}

		cancelFn()
		b.Results[name] = []string{}
		for _, log := range rawLogs {
			b.Results[name] = append(b.Results[name], log.Log)
		}
	}

	if len(b.Errors) > 0 {
		return errors.New("there were errors while fetching the results. Please check the errors and try again")
	}

	resourceErr := b.fetchResources()
	if resourceErr != nil {
		return resourceErr
	}

	return nil
}

func (b *BasicReport) fetchResources() error {
	//TODO in both cases we'd need to know some mask to filter out the resources we need
	if b.ExecutionEnvironment == ExecutionEnvironment_Docker {
		// fetch docker resources
		// get all containers and their resources
	} else {
		// fetch k8s resources
		// get all pods and their resources
	}

	return nil
}

func (b *BasicReport) fetchDockerResources() error {
	provider, err := tc.NewDockerProvider()
	if err != nil {
		return fmt.Errorf("failed to create Docker provider: %w", err)
	}

	containers, err := provider.Client().ContainerList(context.Background(), container.ListOptions{All: true})
	if err != nil {
		return fmt.Errorf("failed to list Docker containers: %w", err)
	}

	eg := &errgroup.Group{}
	pattern := regexp.MustCompile(b.ResourceSelectionPattern)

	var dockerResources = make(map[string]*DockerResources)

	for _, containerInfo := range containers {
		eg.Go(func() error {
			containerName := containerInfo.Names[0]
			if !pattern.Match([]byte(containerName)) {
				return nil
			}

			ctx, cancelFn := context.WithTimeout(context.Background(), 30*time.Second)
			info, err := provider.Client().ContainerInspect(ctx, containerInfo.ID)
			if err != nil {
				cancelFn()
				return errors.Wrapf(err, "failed to inspect container %s", containerName)
			}

			cancelFn()
			dockerResources[containerName] = &DockerResources{
				NanoCPUs:   info.HostConfig.NanoCPUs,
				CpuShares:  info.HostConfig.CPUShares,
				Memory:     info.HostConfig.Memory,
				MemorySwap: info.HostConfig.MemorySwap,
			}

			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return errors.Wrapf(err, "failed to fetch Docker resources")
	}

	b.ContainerResources = dockerResources

	return nil
}

func (b *BasicReport) IsComparable(otherReport BasicReport) (bool, []error) {
	// check if generator configs are the same
	// are all configs present? do they have the same schedule type? do they have the same segments?
	// is call timeout the same?
	// is rate limit timeout the same?
	// would be good to be able to check if Gun and VU are the same, but idk yet how we could do that easily [hash the code?]

	if len(b.GeneratorConfigs) != len(otherReport.GeneratorConfigs) {
		return false, []error{fmt.Errorf("generator configs count is different. Expected %d, got %d", len(b.GeneratorConfigs), len(otherReport.GeneratorConfigs))}
	}

	for name1, cfg1 := range b.GeneratorConfigs {
		if cfg2, ok := otherReport.GeneratorConfigs[name1]; !ok {
			return false, []error{fmt.Errorf("generator config %s is missing from the other report", name1)}
		} else {
			if err := compareGeneratorConfigs(cfg1, cfg2); err != nil {
				return false, []error{err}
			}
		}
	}

	for name2 := range otherReport.GeneratorConfigs {
		if _, ok := b.GeneratorConfigs[name2]; !ok {
			return false, []error{fmt.Errorf("generator config %s is missing from the current report", name2)}
		}
	}

	if b.ExecutionEnvironment != otherReport.ExecutionEnvironment {
		return false, []error{fmt.Errorf("execution environments are different. Expected %s, got %s", b.ExecutionEnvironment, otherReport.ExecutionEnvironment)}
	}

	// check if pods resources are the same
	// are all pods present? do they have the same resources?
	if b.ExecutionEnvironment == ExecutionEnvironment_Docker {
		err := compareDockerResources(b.ContainerResources, otherReport.ContainerResources)
		if err != nil {
			return false, []error{err}
		}
	} else {
		err := comparePodResources(b.PodsResources, otherReport.PodsResources)
		if err != nil {
			return false, []error{err}
		}
	}

	// check if queries are the same
	// are all queries present? do they have the same template?
	lokiQueriesErr := compareLokiQueries(b.LokiQueries, otherReport.LokiQueries)
	if lokiQueriesErr != nil {
		return false, []error{lokiQueriesErr}
	}

	return true, nil
}

func compareLokiQueries(this, other map[string]string) error {
	if len(this) != len(other) {
		return fmt.Errorf("queries count is different. Expected %d, got %d", len(this), len(other))
	}

	for name1, query1 := range this {
		if query2, ok := other[name1]; !ok {
			return fmt.Errorf("query %s is missing from the other report", name1)
		} else {
			if query1 != query2 {
				return fmt.Errorf("query %s is different. Expected %s, got %s", name1, query1, query2)
			}
		}
	}

	for name2 := range other {
		if _, ok := this[name2]; !ok {
			return fmt.Errorf("query %s is missing from the current report", name2)
		}
	}

	return nil
}

func comparePodResources(this, other map[string]*PodResources) error {
	if len(this) != len(other) {
		return fmt.Errorf("pod resources count is different. Expected %d, got %d", len(this), len(other))
	}

	for name1, res1 := range this {
		if res2, ok := other[name1]; !ok {
			return fmt.Errorf("pod resource %s is missing from the other report", name1)
		} else {
			if res1 == nil {
				return fmt.Errorf("pod resource %s is nil in the current report", name1)
			}
			if res2 == nil {
				return fmt.Errorf("pod resource %s is nil in the other report", name1)
			}
			if *res1 != *res2 {
				return fmt.Errorf("pod resource %s is different. Expected %v, got %v", name1, res1, res2)
			}
		}
	}

	for name2 := range other {
		if _, ok := this[name2]; !ok {
			return fmt.Errorf("pod resource %s is missing from the current report", name2)
		}
	}

	return nil
}

func compareDockerResources(this, other map[string]*DockerResources) error {
	if len(this) != len(other) {
		return fmt.Errorf("container resources count is different. Expected %d, got %d", len(this), len(other))
	}

	for name1, res1 := range this {
		if res2, ok := other[name1]; !ok {
			return fmt.Errorf("container resource %s is missing from the other report", name1)
		} else {
			if res1 == nil {
				return fmt.Errorf("container resource %s is nil in the current report", name1)
			}
			if res2 == nil {
				return fmt.Errorf("container resource %s is nil in the other report", name1)
			}
			if *res1 != *res2 {
				return fmt.Errorf("container resource %s is different. Expected %v, got %v", name1, res1, res2)
			}
		}
	}

	for name2 := range other {
		if _, ok := this[name2]; !ok {
			return fmt.Errorf("container resource %s is missing from the current report", name2)
		}
	}

	return nil
}

func compareGeneratorConfigs(cfg1, cfg2 *wasp.Config) error {
	if cfg1.LoadType != cfg2.LoadType {
		return fmt.Errorf("load types are different. Expected %s, got %s", cfg1.LoadType, cfg2.LoadType)
	}

	if len(cfg1.Schedule) != len(cfg2.Schedule) {
		return fmt.Errorf("schedules are different. Expected %d, got %d", len(cfg1.Schedule), len(cfg2.Schedule))
	}

	for i, segment1 := range cfg1.Schedule {
		segment2 := cfg2.Schedule[i]
		if segment1 == nil {
			return fmt.Errorf("schedule at index %d is nil in the current report", i)
		}
		if segment2 == nil {
			return fmt.Errorf("schedule at index %d is nil in the other report", i)
		}
		if *segment1 != *segment2 {
			return fmt.Errorf("schedules at index %d are different. Expected %s, got %s", i, mustMarshallSegment(segment1), mustMarshallSegment(segment2))
		}
	}

	if cfg1.CallTimeout != cfg2.CallTimeout {
		return fmt.Errorf("call timeouts are different. Expected %s, got %s", cfg1.CallTimeout, cfg2.CallTimeout)
	}

	if cfg1.RateLimitUnitDuration != cfg2.RateLimitUnitDuration {
		return fmt.Errorf("rate limit unit durations are different. Expected %s, got %s", cfg1.RateLimitUnitDuration, cfg2.RateLimitUnitDuration)
	}

	return nil
}

func mustMarshallSegment(segment *wasp.Segment) string {
	segmentBytes, err := json.MarshalIndent(segment, "", " ")
	if err != nil {
		panic(err)
	}

	return string(segmentBytes)
}

//type DifferentResult struct {
//	Expected []string
//	Actual   []string
//}
//
//// this could also just be a generic var, so we don't couple Result with this struct
//type ComparisonReport struct {
//	missingResults    []string
//	unexpectedResults []string
//	differentResults  map[string]DifferentResult
//}
//
//func (c *ComparisonReport) HasDifferences() bool {
//	return len(c.missingResults) > 0 || len(c.unexpectedResults) > 0 || len(c.differentResults) > 0
//}
//
//func (c *ComparisonReport) MissingResults() []string {
//	return c.missingResults
//}
//
//func (c *ComparisonReport) UnexpectedResults() []string {
//	return c.unexpectedResults
//}
//
//func (c *ComparisonReport) DifferentResults() map[string]DifferentResult {
//	return c.differentResults
//}
//
//func (b *BasicReport) Compare(other BasicReport) (bool, ComparisonReport, error) {
//	// check if there are no errors
//
//	// check if results are the same
//	if len(b.Results) != len(other.Results) {
//		return false, ComparisonReport{}, fmt.Errorf("result count is different. Expected %d, got %d", len(b.Results), len(other.Results))
//	}
//
//	comparisonReport := ComparisonReport{
//		differentResults: make(map[string]DifferentResult),
//	}
//
//	for name, result := range b.Results {
//		if otherResult, ok := other.Results[name]; !ok {
//			comparisonReport.missingResults = append(comparisonReport.missingResults, name)
//		} else {
//
//			comparisonReport.differentResults[name] = DifferentResult{
//				Expected: result,
//				Actual:   otherResult,
//			}
//		}
//	}
//
//	for name, otherResult := range other.Results {
//		if _, ok := b.Results[name]; !ok {
//			comparisonReport.unexpectedResults[name] = otherResult
//		}
//	}
//
//	return comparisonReport.HasDifferences(), comparisonReport, nil
//}

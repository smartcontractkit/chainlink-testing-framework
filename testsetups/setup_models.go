// package testsetups compresses common test setups and more complicated setups like performance and chaos tests.
package testsetups

// PerformanceTest enables generic interaction and usage of performance tests
type PerformanceTest interface {
	Setup()
	Run()
}

package client

type TransactionsData struct {
	Data []TransactionData    `json:"data"`
	Meta TransactionsMetaData `json:"meta"`
}

type TransactionData struct {
	Type       string                `json:"type"`
	ID         string                `json:"id"`
	Attributes TransactionAttributes `json:"attributes"`
}

type TransactionAttributes struct {
	State    string `json:"state"`
	Data     string `json:"data"`
	From     string `json:"from"`
	To       string `json:"to"`
	Value    string `json:"value"`
	ChainID  string `json:"evmChainID"`
	GasLimit string `json:"gasLimit"`
	GasPrice string `json:"gasPrice"`
	Hash     string `json:"hash"`
	RawHex   string `json:"rawHex"`
	Nonce    string `json:"nonce"`
	SentAt   string `json:"sentAt"`
}

type TransactionsMetaData struct {
	Count int `json:"count"`
}

// ChainlinkProfileResults holds the results of asking the Chainlink node to run a PPROF session
type ChainlinkProfileResults struct {
	Reports                 []*ChainlinkProfileResult
	ScheduledProfileSeconds int // How long the profile was scheduled to last
	ActualRunSeconds        int // How long the target function to profile actually took to execute
	NodeIndex               int
}

// ChainlinkProfileResult contains the result of a single PPROF run
type ChainlinkProfileResult struct {
	Type string
	Data []byte
}

// NewBlankChainlinkProfileResults returns all the standard types of profile results with blank data
func NewBlankChainlinkProfileResults() *ChainlinkProfileResults {
	results := &ChainlinkProfileResults{
		Reports: make([]*ChainlinkProfileResult, 0),
	}
	profileStrings := []string{
		"allocs", // A sampling of all past memory allocations
		"block",  // Stack traces that led to blocking on synchronization primitives
		// "cmdline",      // The command line invocation of the current program
		"goroutine",    // Stack traces of all current goroutines
		"heap",         // A sampling of memory allocations of live objects.
		"mutex",        // Stack traces of holders of contended mutexes
		"profile",      // CPU profile.
		"threadcreate", // Stack traces that led to the creation of new OS threads
		"trace",        // A trace of execution of the current program.
	}
	for _, profile := range profileStrings {
		results.Reports = append(results.Reports, &ChainlinkProfileResult{Type: profile})
	}
	return results
}

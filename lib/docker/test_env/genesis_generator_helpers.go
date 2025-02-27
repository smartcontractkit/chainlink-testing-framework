package test_env

import (
	"bytes"
	"context"
	"html/template"
	"time"

	tcwait "github.com/testcontainers/testcontainers-go/wait"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/config"
)

func generateEnvValues(c *config.EthereumChainConfig) (string, error) {
	// GenesisTimestamp needs to be exported in order to be used in the template
	// but I don't want to expose it in config struct, user should not set it manually
	data := struct {
		config.EthereumChainConfig
		GenesisTimestamp int
	}{
		EthereumChainConfig: *c,
		GenesisTimestamp:    c.GenesisTimestamp,
	}
	tmpl, err := template.New("valuesEnv").Funcs(funcMap).Parse(valuesEnv)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

var funcMap = template.FuncMap{
	"decrement": func(i int) int {
		return i - 1
	},
}

var valuesEnv = `
export PRESET_BASE="mainnet"
export CHAIN_ID="{{.ChainID}}"
export DEPOSIT_CONTRACT_ADDRESS="0x4242424242424242424242424242424242424242"
export EL_AND_CL_MNEMONIC="giant issue aisle success illegal bike spike question tent bar rely arctic volcano long crawl hungry vocal artwork sniff fantasy very lucky have athlete"
export CL_EXEC_BLOCK="0"
export SLOT_DURATION_IN_SECONDS={{.SecondsPerSlot}}
export DEPOSIT_CONTRACT_BLOCK="0x0000000000000000000000000000000000000000000000000000000000000000"
export NUMBER_OF_VALIDATORS={{.ValidatorCount}}
export GENESIS_FORK_VERSION="0x10000038"
export ALTAIR_FORK_VERSION="0x20000038"
export BELLATRIX_FORK_VERSION="0x30000038"
export CAPELLA_FORK_VERSION="0x40000038"
export CAPELLA_FORK_EPOCH="0"
export DENEB_FORK_VERSION="0x50000038"
export DENEB_FORK_EPOCH="{{.HardForkEpochs.Deneb}}"
export ELECTRA_FORK_VERSION="0x60000038"
#export ELECTRA_FORK_EPOCH="{{.HardForkEpochs.Electra}}"
export EIP7594_FORK_VERSION="0x70000000"
#export EIP7594_FORK_EPOCH="{{.HardForkEpochs.EOF}}"
export WITHDRAWAL_TYPE="0x00"
export WITHDRAWAL_ADDRESS=0xf97e180c050e5Ab072211Ad2C213Eb5AEE4DF134
export BEACON_STATIC_ENR="enr:-Iq4QJk4WqRkjsX5c2CXtOra6HnxN-BMXnWhmhEQO9Bn9iABTJGdjUOurM7Btj1ouKaFkvTRoju5vz2GPmVON2dffQKGAX53x8JigmlkgnY0gmlwhLKAlv6Jc2VjcDI1NmsxoQK6S-Cii_KmfFdUJL2TANL3ksaKUnNXvTCv1tLwXs0QgIN1ZHCCIyk"
export SHADOW_FORK_RPC=""
export SHADOW_FORK_FILE=""
export GENESIS_TIMESTAMP={{.GenesisTimestamp}}
export GENESIS_DELAY={{.GenesisDelay}}
export MAX_PER_EPOCH_ACTIVATION_CHURN_LIMIT="8"
export CHURN_LIMIT_QUOTIENT="65536"
export EJECTION_BALANCE=16000000000
export SLOTS_PER_EPOCH={{.SlotsPerEpoch}}
export PREMINE_ADDRS={{ if .AddressesToFund }}'
{{- $lastIndex := decrement (len .AddressesToFund) }}
{{- range $i, $addr := .AddressesToFund }}
  "{{ $addr }}": 1000000000ETH
{{- end }}'
{{ else }}{}
{{ end }}
export ETH1_FOLLOW_DISTANCE="2048"
export MIN_VALIDATOR_WITHDRAWABILITY_DELAY="256"
export SHARD_COMMITTEE_PERIOD="256"
export SAMPLES_PER_SLOT="8"
export CUSTODY_REQUIREMENT="1"
export DATA_COLUMN_SIDECAR_SUBNET_COUNT="32"
export TARGET_NUMBER_OF_PEERS="70"
export ADDITIONAL_PRELOADED_CONTRACTS="{}"
`

var elGenesisConfig = `
preset_base: ${PRESET_BASE}
chain_id: ${CHAIN_ID}
deposit_contract_address: "${DEPOSIT_CONTRACT_ADDRESS}"
mnemonic: ${EL_AND_CL_MNEMONIC}
el_premine:
  "m/44'/60'/0'/0/0": 1000000000ETH
  "m/44'/60'/0'/0/1": 1000000000ETH
  "m/44'/60'/0'/0/2": 1000000000ETH
  "m/44'/60'/0'/0/3": 1000000000ETH
  "m/44'/60'/0'/0/4": 1000000000ETH
  "m/44'/60'/0'/0/5": 1000000000ETH
  "m/44'/60'/0'/0/6": 1000000000ETH
  "m/44'/60'/0'/0/7": 1000000000ETH
  "m/44'/60'/0'/0/8": 1000000000ETH
  "m/44'/60'/0'/0/9": 1000000000ETH
  "m/44'/60'/0'/0/10": 1000000000ETH
  "m/44'/60'/0'/0/11": 1000000000ETH
  "m/44'/60'/0'/0/12": 1000000000ETH
  "m/44'/60'/0'/0/13": 1000000000ETH
  "m/44'/60'/0'/0/14": 1000000000ETH
  "m/44'/60'/0'/0/15": 1000000000ETH
  "m/44'/60'/0'/0/16": 1000000000ETH
  "m/44'/60'/0'/0/17": 1000000000ETH
  "m/44'/60'/0'/0/18": 1000000000ETH
  "m/44'/60'/0'/0/19": 1000000000ETH
  "m/44'/60'/0'/0/20": 1000000000ETH
el_premine_addrs: ${PREMINE_ADDRS}
genesis_timestamp: ${GENESIS_TIMESTAMP}
genesis_delay: ${GENESIS_DELAY}
genesis_gaslimit: ${GENESIS_GASLIMIT}
deneb_fork_epoch: ${DENEB_FORK_EPOCH}
additional_preloaded_contracts: ${ADDITIONAL_PRELOADED_CONTRACTS}
slot_duration_in_seconds: ${SLOT_DURATION_IN_SECONDS}
electra_fork_epoch: ${ELECTRA_FORK_EPOCH}
eof_activation_epoch: ${EOF_ACTIVATION_EPOCH}
slots_per_epoch: ${SLOTS_PER_EPOCH}
`

var clGenesisConfig = `
# Extends the mainnet preset
PRESET_BASE: $PRESET_BASE
CONFIG_NAME: testnet # needs to exist because of Prysm. Otherwise it conflicts with mainnet genesis

# Genesis
# ---------------------------------------------------------------
# 2**14 (= 16,384)
MIN_GENESIS_ACTIVE_VALIDATOR_COUNT: $NUMBER_OF_VALIDATORS
# Mar-01-2021 08:53:32 AM +UTC
# This is an invalid valid and should be updated when you create the genesis
MIN_GENESIS_TIME: $GENESIS_TIMESTAMP
GENESIS_FORK_VERSION: $GENESIS_FORK_VERSION
GENESIS_DELAY: $GENESIS_DELAY


# Forking
# ---------------------------------------------------------------
# Some forks are disabled for now:
#  - These may be re-assigned to another fork-version later
#  - Temporarily set to max uint64 value: 2**64 - 1

# Altair
ALTAIR_FORK_VERSION: $ALTAIR_FORK_VERSION
ALTAIR_FORK_EPOCH: 0
# Merge
BELLATRIX_FORK_VERSION: $BELLATRIX_FORK_VERSION
BELLATRIX_FORK_EPOCH: 0
TERMINAL_TOTAL_DIFFICULTY: 0
TERMINAL_BLOCK_HASH: 0x0000000000000000000000000000000000000000000000000000000000000000
TERMINAL_BLOCK_HASH_ACTIVATION_EPOCH: 18446744073709551615

# Capella
CAPELLA_FORK_VERSION: $CAPELLA_FORK_VERSION
CAPELLA_FORK_EPOCH: 0

# DENEB
DENEB_FORK_VERSION: $DENEB_FORK_VERSION
DENEB_FORK_EPOCH: $DENEB_FORK_EPOCH

# Electra
ELECTRA_FORK_VERSION: $ELECTRA_FORK_VERSION
ELECTRA_FORK_EPOCH: $ELECTRA_FORK_EPOCH

# EIP7594 - Peerdas
EIP7594_FORK_VERSION: $EIP7594_FORK_VERSION
EIP7594_FORK_EPOCH: $EIP7594_FORK_EPOCH

# Time parameters
# ---------------------------------------------------------------
# 12 seconds
SECONDS_PER_SLOT: $SLOT_DURATION_IN_SECONDS
# 14 (estimate from Eth1 mainnet)
SECONDS_PER_ETH1_BLOCK: $SLOT_DURATION_IN_SECONDS
# 2**8 (= 256) epochs ~27 hours
MIN_VALIDATOR_WITHDRAWABILITY_DELAY: $MIN_VALIDATOR_WITHDRAWABILITY_DELAY
# 2**8 (= 256) epochs ~27 hours
SHARD_COMMITTEE_PERIOD: $SHARD_COMMITTEE_PERIOD
# 2**11 (= 2,048) Eth1 blocks ~8 hours
ETH1_FOLLOW_DISTANCE: $ETH1_FOLLOW_DISTANCE

SLOTS_PER_EPOCH: $SLOTS_PER_EPOCH

# Validator cycle
# ---------------------------------------------------------------
# 2**2 (= 4)
INACTIVITY_SCORE_BIAS: 4
# 2**4 (= 16)
INACTIVITY_SCORE_RECOVERY_RATE: 16
# 2**4 * 10**9 (= 16,000,000,000) Gwei
EJECTION_BALANCE: $EJECTION_BALANCE
# 2**2 (= 4)
MIN_PER_EPOCH_CHURN_LIMIT: 4
# 2**16 (= 65,536)
CHURN_LIMIT_QUOTIENT: $CHURN_LIMIT_QUOTIENT
# [New in Deneb:EIP7514] 2**3 (= 8)
MAX_PER_EPOCH_ACTIVATION_CHURN_LIMIT: $MAX_PER_EPOCH_ACTIVATION_CHURN_LIMIT

# Fork choice
# ---------------------------------------------------------------
# 40%
PROPOSER_SCORE_BOOST: 40
# 20%
REORG_HEAD_WEIGHT_THRESHOLD: 20
# 160%
REORG_PARENT_WEIGHT_THRESHOLD: 160
# 2 epochs
REORG_MAX_EPOCHS_SINCE_FINALIZATION: 2

# Deposit contract
# ---------------------------------------------------------------
DEPOSIT_CHAIN_ID: $CHAIN_ID
DEPOSIT_NETWORK_ID: $CHAIN_ID
DEPOSIT_CONTRACT_ADDRESS: $DEPOSIT_CONTRACT_ADDRESS

# Networking
# ---------------------------------------------------------------
# 10 * 2**20 (= 10485760, 10 MiB)
GOSSIP_MAX_SIZE: 10485760
# 2**10 (= 1024)
MAX_REQUEST_BLOCKS: 1024
# 2**8 (= 256)
EPOCHS_PER_SUBNET_SUBSCRIPTION: 256
# MIN_VALIDATOR_WITHDRAWABILITY_DELAY + CHURN_LIMIT_QUOTIENT // 2 (= 33024, ~5 months)
MIN_EPOCHS_FOR_BLOCK_REQUESTS: 33024
# 10 * 2**20 (=10485760, 10 MiB)
MAX_CHUNK_SIZE: 10485760
# 5s
TTFB_TIMEOUT: 5
# 10s
RESP_TIMEOUT: 10
ATTESTATION_PROPAGATION_SLOT_RANGE: 32
# 500ms
MAXIMUM_GOSSIP_CLOCK_DISPARITY: 500
MESSAGE_DOMAIN_INVALID_SNAPPY: 0x00000000
MESSAGE_DOMAIN_VALID_SNAPPY: 0x01000000
# 2 subnets per node
SUBNETS_PER_NODE: 2
# 2**8 (= 64)
ATTESTATION_SUBNET_COUNT: 64
ATTESTATION_SUBNET_EXTRA_BITS: 0
# ceillog2(ATTESTATION_SUBNET_COUNT) + ATTESTATION_SUBNET_EXTRA_BITS
ATTESTATION_SUBNET_PREFIX_BITS: 6

# Deneb
# 2**7 (=128)
MAX_REQUEST_BLOCKS_DENEB: 128
# MAX_REQUEST_BLOCKS_DENEB * MAX_BLOBS_PER_BLOCK
MAX_REQUEST_BLOB_SIDECARS: 768
# 2**12 (= 4096 epochs, ~18 days)
MIN_EPOCHS_FOR_BLOB_SIDECARS_REQUESTS: 4096
# 6
BLOB_SIDECAR_SUBNET_COUNT: 6

# Whisk
# Epoch(2**8)
WHISK_EPOCHS_PER_SHUFFLING_PHASE: 256
# Epoch(2)
WHISK_PROPOSER_SELECTION_GAP: 2

# EIP7594
NUMBER_OF_COLUMNS: 128
MAX_CELLS_IN_EXTENDED_MATRIX: 768
DATA_COLUMN_SIDECAR_SUBNET_COUNT: $DATA_COLUMN_SIDECAR_SUBNET_COUNT
MAX_REQUEST_DATA_COLUMN_SIDECARS: 16384
SAMPLES_PER_SLOT: $SAMPLES_PER_SLOT
CUSTODY_REQUIREMENT: $CUSTODY_REQUIREMENT
TARGET_NUMBER_OF_PEERS: $TARGET_NUMBER_OF_PEERS

# [New in Electra:EIP7251]
MIN_PER_EPOCH_CHURN_LIMIT_ELECTRA: 128000000000 # 2**7 * 10**9 (= 128,000,000,000)
MAX_PER_EPOCH_ACTIVATION_EXIT_CHURN_LIMIT: 256000000000 # 2**8 * 10**9 (= 256,000,000,000)
`

var mnemonics = `
- mnemonic: "${EL_AND_CL_MNEMONIC}"  # a 24 word BIP 39 mnemonic
  count: $NUMBER_OF_VALIDATORS
`

type ExitCodeStrategy struct {
	expectedExitCode int
	timeout          time.Duration
	pollInterval     time.Duration
}

// NewExitCodeStrategy initializes a new ExitCodeStrategy with default settings.
// It sets the expected exit code to 0, a timeout of 2 minutes, and a poll interval of 2 seconds.
// This function is useful for configuring container readiness checks based on exit codes.
func NewExitCodeStrategy() *ExitCodeStrategy {
	return &ExitCodeStrategy{
		expectedExitCode: 0,
		timeout:          2 * time.Minute,
		pollInterval:     2 * time.Second,
	}
}

// WithTimeout sets the timeout duration for the ExitCodeStrategy.
// It allows users to specify how long to wait for a process to exit with the expected code.
func (w *ExitCodeStrategy) WithTimeout(timeout time.Duration) *ExitCodeStrategy {
	w.timeout = timeout
	return w
}

// WithExitCode sets the expected exit code for the strategy and returns the updated instance.
// This function is useful for configuring container requests to wait for specific exit conditions.
func (w *ExitCodeStrategy) WithExitCode(exitCode int) *ExitCodeStrategy {
	w.expectedExitCode = exitCode
	return w
}

// WithPollInterval sets the interval for polling the exit code and returns the updated strategy.
// This function is useful for configuring how frequently to check the exit code during container execution.
func (w *ExitCodeStrategy) WithPollInterval(pollInterval time.Duration) *ExitCodeStrategy {
	w.pollInterval = pollInterval
	return w
}

// WaitUntilReady implements Strategy.WaitUntilReady
func (w *ExitCodeStrategy) WaitUntilReady(ctx context.Context, target tcwait.StrategyTarget) (err error) {

	ctx, cancel := context.WithTimeout(ctx, w.timeout)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(100 * time.Millisecond):
			state, err := target.State(ctx)
			if err != nil {
				return err
			}

			if state.ExitCode != w.expectedExitCode {
				time.Sleep(w.pollInterval)
				continue
			}
			return nil
		}
	}
}

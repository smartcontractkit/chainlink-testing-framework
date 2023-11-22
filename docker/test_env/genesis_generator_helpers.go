package test_env

import (
	"bytes"
	"context"
	"html/template"
	"time"

	tcwait "github.com/testcontainers/testcontainers-go/wait"
)

func generateEnvValues(config *EthereumChainConfig) (string, error) {
	// GenesisTimestamp needs to be exported in order to be used in the template
	// but I don't want to expose it in config struct, user should not set it manually
	data := struct {
		EthereumChainConfig
		GenesisTimestamp int
	}{
		EthereumChainConfig: *config,
		GenesisTimestamp:    config.genesisTimestamp,
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
export DENEB_FORK_EPOCH="500"
export ELECTRA_FORK_VERSION="0x60000038"
export ELECTRA_FORK_EPOCH=""
export WITHDRAWAL_TYPE="0x00"
export WITHDRAWAL_ADDRESS=0xf97e180c050e5Ab072211Ad2C213Eb5AEE4DF134
export BEACON_STATIC_ENR="enr:-Iq4QJk4WqRkjsX5c2CXtOra6HnxN-BMXnWhmhEQO9Bn9iABTJGdjUOurM7Btj1ouKaFkvTRoju5vz2GPmVON2dffQKGAX53x8JigmlkgnY0gmlwhLKAlv6Jc2VjcDI1NmsxoQK6S-Cii_KmfFdUJL2TANL3ksaKUnNXvTCv1tLwXs0QgIN1ZHCCIyk"
export GENESIS_TIMESTAMP={{.GenesisTimestamp}}
export GENESIS_DELAY={{.GenesisDelay}}
export MAX_CHURN=8
export EJECTION_BALANCE=16000000000
export SLOTS_PER_EPOCH={{.SlotsPerEpoch}}
export PREMINE_ADDRS={{ if .AddressesToFund }}'
{{- $lastIndex := decrement (len .AddressesToFund) }}
{{- range $i, $addr := .AddressesToFund }}
  "{{ $addr }}": 1000000000ETH
{{- end }}'
{{ else }}{}
{{ end }}
`

var elGenesisConfig = `
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
slot_duration_in_seconds: ${SLOT_DURATION_IN_SECONDS}
deneb_fork_epoch: ${DENEB_FORK_EPOCH}
`

var clGenesisConfig = `
# Extends the mainnet preset
PRESET_BASE: 'mainnet'
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

# Time parameters
# ---------------------------------------------------------------
# 12 seconds
SECONDS_PER_SLOT: $SLOT_DURATION_IN_SECONDS
# 14 (estimate from Eth1 mainnet)
SECONDS_PER_ETH1_BLOCK: $SLOT_DURATION_IN_SECONDS
# 2**0 (= 1) epochs ~1 hours
MIN_VALIDATOR_WITHDRAWABILITY_DELAY: 1
# 2**8 (= 256) epochs ~27 hours
SHARD_COMMITTEE_PERIOD: 1
# 2**11 (= 2,048) Eth1 blocks ~8 hours
ETH1_FOLLOW_DISTANCE: 12

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
CHURN_LIMIT_QUOTIENT: 65536
# [New in Deneb:EIP7514] 2**3 (= 8)
MAX_PER_EPOCH_ACTIVATION_CHURN_LIMIT: $MAX_CHURN

# Fork choice
# ---------------------------------------------------------------
# 40%
PROPOSER_SCORE_BOOST: 40

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
# uint64(6)
MAX_BLOBS_PER_BLOCK: 6
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

func NewExitCodeStrategy() *ExitCodeStrategy {
	return &ExitCodeStrategy{
		expectedExitCode: 0,
		timeout:          2 * time.Minute,
		pollInterval:     2 * time.Second,
	}
}

func (w *ExitCodeStrategy) WithTimeout(timeout time.Duration) *ExitCodeStrategy {
	w.timeout = timeout
	return w
}

func (w *ExitCodeStrategy) WithExitCode(exitCode int) *ExitCodeStrategy {
	w.expectedExitCode = exitCode
	return w
}

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
			} else {
				return nil
			}
		}
	}
}

package seth

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"time"

	"math/big"
	"sync"

	"github.com/avast/retry-go"
	"github.com/ethereum/go-ethereum/common"
	"go.uber.org/ratelimit"
)

const (
	ErrKeySyncTimeout = "key sync timeout, consider increasing key_sync_timeout in config (seth.toml or ClientBuilder), or increasing the number of keys"
	ErrKeySync        = "failed to sync the key"
	ErrNonce          = "failed to get nonce"
	TimeoutKeyNum     = -80001
)

// NonceManager tracks nonce for each address
type NonceManager struct {
	*sync.Mutex
	cfg         *NonceManagerCfg
	rl          ratelimit.Limiter
	Client      *Client
	SyncTimeout time.Duration
	SyncedKeys  chan *KeyNonce
	Addresses   []common.Address
	PrivateKeys []*ecdsa.PrivateKey
	Nonces      map[common.Address]int64
}

type KeyNonce struct {
	KeyNum int
	Nonce  uint64
}

func validateNonceManagerConfig(nonceManagerCfg *NonceManagerCfg) error {
	if nonceManagerCfg.KeySyncRateLimitSec <= 0 {
		return fmt.Errorf("key_sync_rate_limit_sec must be positive (current: %d). "+
			"This controls how many sync attempts per second are allowed. "+
			"Set it in the 'nonce_manager' section of config (seth.toml or ClientBuilder)",
			nonceManagerCfg.KeySyncRateLimitSec)
	}
	if nonceManagerCfg.KeySyncTimeout == nil || nonceManagerCfg.KeySyncTimeout.Duration() <= 0 {
		return fmt.Errorf("key_sync_timeout must be positive (current: %v). "+
			"This is how long to wait for a key to sync before timing out. "+
			"Set it in the 'nonce_manager' section of config (seth.toml or ClientBuilder)",
			nonceManagerCfg.KeySyncTimeout)
	}
	if nonceManagerCfg.KeySyncRetries <= 0 {
		return fmt.Errorf("key_sync_retries must be positive (current: %d). "+
			"This is how many times to retry syncing a key before giving up. "+
			"Set it in the 'nonce_manager' section of config (seth.toml or ClientBuilder)",
			nonceManagerCfg.KeySyncRetries)
	}

	return nil
}

// NewNonceManager creates a new nonce manager that tracks nonce for each address
func NewNonceManager(cfg *Config, addrs []common.Address, privKeys []*ecdsa.PrivateKey) (*NonceManager, error) {
	if cfg == nil {
		return nil, fmt.Errorf("seth configuration is nil. Cannot create nonce manager without valid configuration.\n" +
			"This usually means you're trying to create a nonce manager before initializing Seth.\n" +
			"Solutions:\n" +
			"  1. Use NewClient() or NewClientWithConfig() to create a Seth client first\n" +
			"  2. If using ClientBuilder, ensure you call Build() before accessing the nonce manager\n" +
			"  3. Check that your configuration file (seth.toml) exists and is valid")
	}
	if cfg.NonceManager == nil {
		return nil, fmt.Errorf("nonce manager configuration is nil. " +
			"Add a [nonce_manager] section to your config (seth.toml) or use ClientBuilder with:\n" +
			"  - key_sync_rate_limit_per_sec\n" +
			"  - key_sync_timeout\n" +
			"  - key_sync_retries\n" +
			"  - key_sync_retry_delay")
	}
	if cfgErr := validateNonceManagerConfig(cfg.NonceManager); cfgErr != nil {
		return nil, fmt.Errorf("nonce manager configuration validation failed: %w", cfgErr)
	}

	nonces := make(map[common.Address]int64)
	for _, addr := range addrs {
		nonces[addr] = 0
	}
	return &NonceManager{
		Mutex:       &sync.Mutex{},
		cfg:         cfg.NonceManager,
		rl:          ratelimit.New(cfg.NonceManager.KeySyncRateLimitSec, ratelimit.WithoutSlack),
		Nonces:      nonces,
		Addresses:   addrs,
		PrivateKeys: privKeys,
		SyncedKeys:  make(chan *KeyNonce, len(addrs)),
	}, nil
}

// UpdateNonces syncs nonces for addresses
func (m *NonceManager) UpdateNonces() error {
	L.Debug().Interface("Addrs", m.Addresses).Msg("Updating nonces for addresses")
	m.Lock()
	defer m.Unlock()
	for addr := range m.Nonces {
		nonce, err := m.Client.Client.NonceAt(context.Background(), addr, nil)
		if err != nil {
			return fmt.Errorf("failed to updated nonces for address '%s': %w", addr, err)
		}
		m.Nonces[addr] = mustSafeInt64(nonce)
	}
	L.Debug().Interface("Nonces", m.Nonces).Msg("Updated nonces for addresses")
	m.SyncedKeys = make(chan *KeyNonce, len(m.Addresses))
	for keyNum, addr := range m.Addresses[1:] {
		m.SyncedKeys <- &KeyNonce{
			KeyNum: keyNum + 1,
			Nonce:  mustSafeUint64(m.Nonces[addr]),
		}
	}
	return nil
}

// NextNonce returns new nonce for addr
// this method is external for module testing, but you should not use it
// since handling nonces on the client is unpredictable
func (m *NonceManager) NextNonce(addr common.Address) *big.Int {
	m.Lock()
	defer m.Unlock()
	nextNonce := big.NewInt(m.Nonces[addr])
	m.Nonces[addr]++
	return nextNonce
}

func (m *NonceManager) anySyncedKey() int {
	ctx, cancel := context.WithTimeout(context.Background(), m.cfg.KeySyncTimeout.Duration())
	defer cancel()
	select {
	case <-ctx.Done():
		m.Lock()
		defer m.Unlock()
		timeoutErr := fmt.Errorf("key synchronization timed out after %s. "+
			"This means the nonce couldn't be synchronized before the timeout.\n"+
			"Solutions:\n"+
			"  1. Increase 'key_sync_timeout' in config (seth.toml or ClientBuilder) - current: %s\n"+
			"  2. Reduce 'key_sync_rate_limit_per_sec' to allow faster sync attempts\n"+
			"  3. Add more keys with 'ephemeral_addresses_number'\n"+
			"  4. Check RPC node performance and connectivity",
			m.cfg.KeySyncTimeout.Duration(), m.cfg.KeySyncTimeout.Duration())
		L.Error().Msg(timeoutErr.Error())
		m.Client.Errors = append(m.Client.Errors, timeoutErr)
		return TimeoutKeyNum //so that it's pretty unique number of invalid key
	case keyData := <-m.SyncedKeys:
		L.Trace().
			Interface("KeyNum", keyData.KeyNum).
			Uint64("Nonce", keyData.Nonce).
			Interface("Address", m.Addresses[keyData.KeyNum]).
			Msg("Key selected")
		go func() {
			err := retry.Do(
				func() error {
					m.rl.Take()
					L.Trace().
						Interface("KeyNum", keyData.KeyNum).
						Interface("Address", m.Addresses[keyData.KeyNum]).
						Msg("Key is syncing")
					nonce, err := m.Client.Client.NonceAt(context.Background(), m.Addresses[keyData.KeyNum], nil)
					if err != nil {
						return fmt.Errorf("failed to get nonce for address %s (key #%d): %w\n"+
							"This usually indicates:\n"+
							"  1. RPC node connection issues\n"+
							"  2. Network congestion or high latency\n"+
							"  3. Address doesn't exist on the network\n"+
							"Consider increasing key_sync_timeout in your config",
							m.Addresses[keyData.KeyNum].Hex(), keyData.KeyNum, err)
					}
					if nonce == keyData.Nonce+1 {
						L.Trace().
							Interface("KeyNum", keyData.KeyNum).
							Uint64("Nonce", nonce).
							Interface("Address", m.Addresses[keyData.KeyNum]).
							Msg("Key synced")
						m.SyncedKeys <- &KeyNonce{
							KeyNum: keyData.KeyNum,
							Nonce:  nonce,
						}
						return nil
					}

					L.Trace().
						Interface("KeyNum", keyData.KeyNum).
						Uint64("Nonce", nonce).
						Int("Expected nonce", mustSafeInt(keyData.Nonce+1)).
						Interface("Address", m.Addresses[keyData.KeyNum]).
						Msg("Key NOT synced")

					return fmt.Errorf("key #%d (address: %s) sync failed. "+
						"Expected nonce %d, but got %d. "+
						"This indicates the transaction hasn't been mined yet",
						keyData.KeyNum, m.Addresses[keyData.KeyNum].Hex(),
						keyData.Nonce+1, nonce)
				},
				retry.Attempts(m.cfg.KeySyncRetries),
				retry.Delay(m.cfg.KeySyncRetryDelay.Duration()),
			)
			if err != nil {
				syncErr := fmt.Errorf("failed to sync key #%d after %d retries: %w",
					keyData.KeyNum, m.cfg.KeySyncRetries, err)
				m.Client.Errors = append(m.Client.Errors, syncErr)
			}
		}()
		return keyData.KeyNum
	}
}

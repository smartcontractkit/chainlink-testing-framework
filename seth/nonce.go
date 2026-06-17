package seth

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"time"

	"math/big"
	"sync"

	"github.com/avast/retry-go"
	"github.com/ethereum/go-ethereum/common"
	"go.uber.org/ratelimit"
)

var (
	ErrKeySync        = errors.New("failed to sync the key")
	ErrKeySyncTimeout = errors.New("key sync timeout, consider increasing key_sync_timeout in config (seth.toml or ClientBuilder), or increasing the number of keys")
	ErrNonce          = errors.New("failed to get nonce")
)

const (
	TimeoutKeyNum = -80001
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
			return err
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
			"  4. Check RPC node performance and connectivity: %w",
			m.cfg.KeySyncTimeout.Duration(), m.cfg.KeySyncTimeout.Duration(), ErrKeySyncTimeout)
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
			addr := m.Addresses[keyData.KeyNum]
			rpcTimeout := max(
				// Use retry delay as RPC timeout
				m.cfg.KeySyncRetryDelay.Duration(), 5*time.Second)

			// Track the last known nonce across retries for recovery
			var lastKnownNonce uint64
			hasValidNonce := false

			err := retry.Do(
				func() error {
					m.rl.Take()
					L.Trace().
						Interface("KeyNum", keyData.KeyNum).
						Interface("Address", addr).
						Msg("Key is syncing")

					rpcCtx, rpcCancel := context.WithTimeout(context.Background(), rpcTimeout)
					defer rpcCancel()

					// Check both pending and latest nonce to determine if key is available
					pendingNonce, pendingErr := m.Client.Client.PendingNonceAt(rpcCtx, addr)
					if pendingErr != nil {
						return fmt.Errorf("failed to get pending nonce for address %s (key #%d): %w\n"+
							"This usually indicates:\n"+
							"  1. RPC node connection issues\n"+
							"  2. Network congestion or high latency\n"+
							"  3. Address doesn't exist on the network\n"+
							"Consider increasing key_sync_timeout in your config",
							m.Addresses[keyData.KeyNum].Hex(), keyData.KeyNum,
							fmt.Errorf("%w: %w", ErrNonce, pendingErr))
					}
					latestNonce, latestErr := m.Client.Client.NonceAt(rpcCtx, addr, nil)
					if latestErr != nil {
						return fmt.Errorf("failed to get latest nonce for address %s (key #%d): %w\n"+
							"This usually indicates:\n"+
							"  1. RPC node connection issues\n"+
							"  2. Network congestion or high latency\n"+
							"  3. Address doesn't exist on the network\n"+
							"Consider increasing key_sync_timeout in your config",
							m.Addresses[keyData.KeyNum].Hex(), keyData.KeyNum,
							fmt.Errorf("%w: %w", ErrNonce, latestErr))
					}

					// Store for potential recovery use
					lastKnownNonce = latestNonce
					hasValidNonce = true

					// Key is synced if there's no pending transaction (pending == latest)
					// OR if the nonce has incremented from what we expected
					if pendingNonce == latestNonce || latestNonce >= keyData.Nonce+1 {
						L.Trace().
							Interface("KeyNum", keyData.KeyNum).
							Uint64("LatestNonce", latestNonce).
							Uint64("PendingNonce", pendingNonce).
							Interface("Address", addr).
							Msg("Key synced")
						m.SyncedKeys <- &KeyNonce{
							KeyNum: keyData.KeyNum,
							Nonce:  latestNonce,
						}
						return nil
					}

					L.Trace().
						Interface("KeyNum", keyData.KeyNum).
						Uint64("LatestNonce", latestNonce).
						Uint64("PendingNonce", pendingNonce).
						Uint64("ExpectedNonce", keyData.Nonce+1).
						Interface("Address", addr).
						Msg("Key NOT synced - has pending transaction")

					return fmt.Errorf("key #%d (address: %s) sync failed. "+
						"Expected nonce %d, but got %d. "+
						"This indicates the transaction hasn't been mined yet: %w",
						keyData.KeyNum, m.Addresses[keyData.KeyNum].Hex(),
						keyData.Nonce+1, latestNonce, ErrKeySync)
				},
				retry.Attempts(m.cfg.KeySyncRetries),
				retry.Delay(m.cfg.KeySyncRetryDelay.Duration()),
			)
			if err != nil {
				m.Client.Errors = append(m.Client.Errors, ErrKeySync)

				// NEVER leak the key - always return it to the pool
				var nonceToUse uint64
				if hasValidNonce {
					nonceToUse = lastKnownNonce
					L.Warn().
						Interface("KeyNum", keyData.KeyNum).
						Uint64("Nonce", nonceToUse).
						Interface("Address", addr).
						Msg("Key sync failed, returning key to pool with last known nonce")
				} else {
					// Fall back to the original nonce when key was checked out
					nonceToUse = keyData.Nonce
					L.Warn().
						Interface("KeyNum", keyData.KeyNum).
						Uint64("Nonce", nonceToUse).
						Interface("Address", addr).
						Msg("Key sync failed, returning key to pool with original nonce")
				}
				m.SyncedKeys <- &KeyNonce{
					KeyNum: keyData.KeyNum,
					Nonce:  nonceToUse,
				}
			}
		}()
		return keyData.KeyNum
	}
}

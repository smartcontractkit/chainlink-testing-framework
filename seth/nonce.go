package seth

import (
	"context"
	"crypto/ecdsa"
	"time"

	"math/big"
	"sync"

	"github.com/avast/retry-go"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"go.uber.org/ratelimit"
)

const (
	ErrKeySyncTimeout = "key sync timeout, consider increasing key_sync_timeout in seth.toml, or increasing the number of keys"
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
		return errors.New("key_sync_rate_limit_sec should be positive")
	}
	if nonceManagerCfg.KeySyncTimeout == nil || nonceManagerCfg.KeySyncTimeout.Duration() <= 0 {
		return errors.New("key_sync_timeout should be positive")
	}
	if nonceManagerCfg.KeySyncRetries <= 0 {
		return errors.New("key_sync_retries should be positive")
	}

	return nil
}

// NewNonceManager creates a new nonce manager that tracks nonce for each address
func NewNonceManager(cfg *Config, addrs []common.Address, privKeys []*ecdsa.PrivateKey) (*NonceManager, error) {
	if cfg == nil {
		return nil, errors.New(ErrSethConfigIsNil)
	}
	if cfg.NonceManager == nil {
		return nil, errors.New(ErrNonceManagerConfigIsNil)
	}
	if cfgErr := validateNonceManagerConfig(cfg.NonceManager); cfgErr != nil {
		return nil, errors.Wrap(cfgErr, "failed to validate nonce manager config")
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
		L.Error().Msg(ErrKeySyncTimeout)
		m.Client.Errors = append(m.Client.Errors, errors.New(ErrKeySync))
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
						return errors.New(ErrNonce)
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

					return errors.New(ErrKeySync)
				},
				retry.Attempts(m.cfg.KeySyncRetries),
				retry.Delay(m.cfg.KeySyncRetryDelay.Duration()),
			)
			if err != nil {
				m.Client.Errors = append(m.Client.Errors, errors.New(ErrKeySync))
			}
		}()
		return keyData.KeyNum
	}
}

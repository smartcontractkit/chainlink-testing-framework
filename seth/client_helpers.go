package seth

import (
	"crypto/ecdsa"
	"errors"
	"github.com/ethereum/go-ethereum/common"
)

// MustGetRootKeyAddress returns the root key address from the client configuration. If no addresses are found, it panics.
// Root key address is the first address in the list of addresses.
func (m *Client) MustGetRootKeyAddress() common.Address {
	if len(m.Addresses) == 0 {
		panic("no addresses found in the client configuration")
	}
	return m.Addresses[0]
}

// GetRootKeyAddress returns the root key address from the client configuration. If no addresses are found, it returns an error.
// Root key address is the first address in the list of addresses.
func (m *Client) GetRootKeyAddress() (common.Address, error) {
	if len(m.Addresses) == 0 {
		return common.Address{}, errors.New("no addresses found in the client configuration")
	}
	return m.Addresses[0], nil
}

// MustGetRootPrivateKey returns the private key of root key/address from the client configuration. If no private keys are found, it panics.
// Root private key is the first private key in the list of private keys.
func (m *Client) MustGetRootPrivateKey() *ecdsa.PrivateKey {
	if len(m.PrivateKeys) == 0 {
		panic("no private keys found in the client configuration")
	}
	return m.PrivateKeys[0]
}

// GetRootPrivateKey returns the private key of root key/address from the client configuration. If no private keys are found, it returns an error.
// Root private key is the first private key in the list of private keys.
func (m *Client) GetRootPrivateKey() (*ecdsa.PrivateKey, error) {
	if len(m.PrivateKeys) == 0 {
		return nil, errors.New("no private keys found in the client configuration")
	}
	return m.PrivateKeys[0], nil
}

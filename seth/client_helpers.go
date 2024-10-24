package seth

import (
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/common"
)

// MustGetRootKeyAddress returns the root key address from the client configuration. If no addresses are found, it panics.
// Root key address is the first address in the list of addresses.
func (m *Client) MustGetRootKeyAddress() common.Address {
	if err := m.validateAddressesKeyNum(0); err != nil {
		panic(err)
	}
	return m.Addresses[0]
}

// GetRootKeyAddress returns the root key address from the client configuration. If no addresses are found, it returns an error.
// Root key address is the first address in the list of addresses.
func (m *Client) GetRootKeyAddress() (common.Address, error) {
	if err := m.validateAddressesKeyNum(0); err != nil {
		return common.Address{}, err
	}
	return m.Addresses[0], nil
}

// MustGetRootPrivateKey returns the private key of root key/address from the client configuration. If no private keys are found, it panics.
// Root private key is the first private key in the list of private keys.
func (m *Client) MustGetRootPrivateKey() *ecdsa.PrivateKey {
	if err := m.validatePrivateKeysKeyNum(0); err != nil {
		panic(err)
	}
	return m.PrivateKeys[0]
}

// GetRootPrivateKey returns the private key of root key/address from the client configuration. If no private keys are found, it returns an error.
// Root private key is the first private key in the list of private keys.
func (m *Client) GetRootPrivateKey() (*ecdsa.PrivateKey, error) {
	if err := m.validatePrivateKeysKeyNum(0); err != nil {
		return nil, err
	}
	return m.PrivateKeys[0], nil
}

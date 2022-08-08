package client

import (
	"time"

	"github.com/smartcontractkit/chainlink/integration-tests/client"
)

type NodeKeysBundle struct {
	OCR2Key client.OCR2Key
	PeerID  string
	TXKey   client.TxKey
	P2PKeys client.P2PKeys
}

type Node struct {
	ID        int32     `json:"ID"`
	Name      string    `json:"Name"`
	ChainID   string    `json:"ChainID"`
	URL       string    `json:"URL"`
	CreatedAt time.Time `json:"CreatedAt"`
	UpdatedAt time.Time `json:"UpdatedAt"`
}

func CreateNodeKeysBundle(nodes []*client.Chainlink, chainName string) ([]NodeKeysBundle, error) {
	nkb := make([]NodeKeysBundle, 0)
	for _, n := range nodes {
		p2pkeys, err := n.MustReadP2PKeys()
		if err != nil {
			return nil, err
		}

		peerID := p2pkeys.Data[0].Attributes.PeerID
		txKey, _, err := n.CreateTxKey(chainName)
		if err != nil {
			return nil, err
		}

		ocrKey, _, err := n.CreateOCR2Key(chainName)
		if err != nil {
			return nil, err
		}
		nkb = append(nkb, NodeKeysBundle{
			PeerID:  peerID,
			OCR2Key: *ocrKey,
			TXKey:   *txKey,
			P2PKeys: *p2pkeys,
		})

	}

	return nkb, nil
}

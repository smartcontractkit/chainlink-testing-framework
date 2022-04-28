package actions

import (
	"bytes"
	"encoding/hex"
	"sort"
	"strings"
	"time"

	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/smartcontractkit/integrations-framework/client"
	"github.com/smartcontractkit/integrations-framework/contracts"
	"github.com/smartcontractkit/libocr/offchainreporting2/confighelper"
	"github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	"golang.org/x/crypto/curve25519"
	"gopkg.in/guregu/null.v4"
)

type NodeKeysBundle struct {
	PeerID  string
	OCR2Key *client.OCR2Key
	TXKey   *client.TxKey
}

// OCR2 keys are in format OCR2<key_type>_<network>_<key>
func stripKeyPrefix(key string) string {
	chunks := strings.Split(key, "_")
	if len(chunks) == 3 {
		return chunks[2]
	}
	return key
}

func CreateEVMChainAndNode(chainID int, chainName string, httpURL string, wsURL string, nodes []client.Chainlink) error {
	ci := *utils.NewBigI(int64(chainID))
	for _, n := range nodes {
		// FIXME: bug #38295, chain is polling the nodes immediately and always returns an error
		_ = n.CreateEVMChain(client.CreateEVMChainRequest{
			ID: ci,
			Config: client.ChainCfg{
				BlockHistoryEstimatorBlockDelay:       null.IntFrom(1),
				BlockHistoryEstimatorBlockHistorySize: null.IntFrom(12),
				EvmEIP1559DynamicFees:                 null.BoolFrom(false),
				MinIncomingConfirmations:              null.IntFrom(10),
			}},
		)
		if err := n.CreateEVMNode(client.NewEVMNode{
			Name:       chainName,
			EVMChainID: ci,
			HTTPURL:    null.NewString(httpURL, true),
			WSURL:      null.NewString(wsURL, true),
		}); err != nil {
			return err
		}
		if err := n.UpdateEVMChain(client.UpdateEVMChainRequest{ID: ci.String(), Enabled: true}); err != nil {
			return err
		}
	}
	return nil
}

func CreateNodeKeysBundle(nodes []client.Chainlink, chainName string) ([]NodeKeysBundle, error) {
	nkb := make([]NodeKeysBundle, 0)
	for _, n := range nodes {
		p2pkeys, err := n.ReadP2PKeys()
		if err != nil {
			return nil, err
		}

		peerID := p2pkeys.Data[0].Attributes.PeerID
		txKey, err := n.CreateTxKey(chainName)
		if err != nil {
			return nil, err
		}
		ocrKey, err := n.CreateOCR2Key(chainName)
		if err != nil {
			return nil, err
		}
		nkb = append(nkb, NodeKeysBundle{
			PeerID:  peerID,
			OCR2Key: ocrKey,
			TXKey:   txKey,
		})
	}
	return nkb, nil
}

func createOracleIdentities(nkb []NodeKeysBundle) ([]confighelper.OracleIdentityExtra, error) {
	oracleIdentities := make([]confighelper.OracleIdentityExtra, 0)
	for _, nodeKeys := range nkb {
		offChainPubKeyRaw, err := hex.DecodeString(stripKeyPrefix(nodeKeys.OCR2Key.Data.Attributes.OffChainPublicKey))
		if err != nil {
			return nil, err
		}
		var offChainPubKey types.OffchainPublicKey
		copy(offChainPubKey[:], offChainPubKeyRaw)
		onChainPubKey, err := hex.DecodeString(stripKeyPrefix(nodeKeys.OCR2Key.Data.Attributes.OnChainPublicKey))
		if err != nil {
			return nil, err
		}
		cfgPubKeyTemp, err := hex.DecodeString(stripKeyPrefix(nodeKeys.OCR2Key.Data.Attributes.ConfigPublicKey))
		if err != nil {
			return nil, err
		}
		cfgPubKeyBytes := [curve25519.PointSize]byte{}
		copy(cfgPubKeyBytes[:], cfgPubKeyTemp)
		oracleIdentities = append(oracleIdentities, confighelper.OracleIdentityExtra{
			OracleIdentity: confighelper.OracleIdentity{
				OffchainPublicKey: offChainPubKey,
				OnchainPublicKey:  onChainPubKey,
				PeerID:            nodeKeys.PeerID,
				TransmitAccount:   types.Account(nodeKeys.TXKey.Data.Attributes.PublicKey),
			},
			ConfigEncryptionPublicKey: cfgPubKeyBytes,
		})
	}
	// program sorts oracles (need to pre-sort to allow correct onchainConfig generation)
	sort.Slice(oracleIdentities, func(i, j int) bool {
		return bytes.Compare(oracleIdentities[i].OracleIdentity.OnchainPublicKey, oracleIdentities[j].OracleIdentity.OnchainPublicKey) < 0
	})
	return oracleIdentities, nil
}

// OffChainConfigParamsFromNodes creates contracts.OffChainAggregatorV2Config
func OffChainConfigParamsFromNodes(nodes []client.Chainlink, nkb []NodeKeysBundle) (contracts.OffChainAggregatorV2Config, error) {
	oi, err := createOracleIdentities(nkb)
	if err != nil {
		return contracts.OffChainAggregatorV2Config{}, err
	}
	s := make([]int, 0)
	for range nodes {
		s = append(s, 1)
	}
	faultyNodes := 0
	if len(nodes) > 1 {
		faultyNodes = len(nodes)/3 - 1
	}
	if faultyNodes == 0 {
		faultyNodes = 1
	}
	return contracts.OffChainAggregatorV2Config{
		DeltaProgress: 2 * time.Second,
		DeltaResend:   5 * time.Second,
		DeltaRound:    1 * time.Second,
		DeltaGrace:    500 * time.Millisecond,
		DeltaStage:    10 * time.Second,
		RMax:          3,
		S:             s,
		Oracles:       oi,
		ReportingPluginConfig: median.OffchainConfig{
			AlphaReportPPB: uint64(0),
			AlphaAcceptPPB: uint64(0),
		}.Encode(),
		MaxDurationQuery:                        0,
		MaxDurationObservation:                  500 * time.Millisecond,
		MaxDurationReport:                       500 * time.Millisecond,
		MaxDurationShouldAcceptFinalizedReport:  500 * time.Millisecond,
		MaxDurationShouldTransmitAcceptedReport: 500 * time.Millisecond,
		F:                                       faultyNodes,
		OnchainConfig:                           []byte{},
	}, nil
}

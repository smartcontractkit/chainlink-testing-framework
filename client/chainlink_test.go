package client

//revive:disable:defer
import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"net/http/httptest"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var spec = `{
  "initiators": [
    {
      "type": "runlog"
    }
  ],
  "tasks": [
    {
      "type": "httpget"
    },
    {
      "type": "jsonparse"
    },
    {
      "type": "multiply"
    },
    {
      "type": "ethuint256"
    },
    {
      "type": "ethtx"
    }
  ]
}`

var _ = Describe("Chainlink @unit", func() {

	// Mocks the creation, read, delete cycle for any job type
	It("can Create, Read, and Delete jobs", func() {
		server := mockedServer(func(rw http.ResponseWriter, req *http.Request) {
			switch req.Method {
			case http.MethodPost:
				Expect(req.URL.Path).Should(Or(Equal("/v2/jobs"), Equal("/sessions")))
				if req.URL.Path == "/sessions" {
					writeCookie(rw)
				} else {
					writeResponse(rw, http.StatusOK, Job{
						Data: JobData{
							ID: "1",
						},
					})
				}
			case http.MethodGet:
				switch req.URL.Path {
				case "/v2/jobs":
					writeResponse(rw, http.StatusOK, ResponseSlice{
						Data: []map[string]interface{}{},
					})
				default:
					Expect("/v2/jobs/1").Should(Equal(req.URL.Path))
					writeResponse(rw, http.StatusOK, Response{
						Data: map[string]interface{}{},
					})
				}
			case http.MethodDelete:
				Expect("/v2/jobs/1").Should(Equal(req.URL.Path))
				writeResponse(rw, http.StatusNoContent, nil)
			}
		})
		defer server.Close()

		c, err := newDefaultClient(server.URL)
		Expect(err).ShouldNot(HaveOccurred())
		c.SetClient(server.Client())

		s, err := c.CreateJobRaw("schemaVersion = 1")
		Expect(err).ShouldNot(HaveOccurred())

		_, err = c.ReadJob(s.Data.ID)
		Expect(err).ShouldNot(HaveOccurred())

		err = c.DeleteJob(s.Data.ID)
		Expect(err).ShouldNot(HaveOccurred())
	})

	// Mocks the creation, read, delete cycle for any job spec
	It("can Create, Read, and Delete job specs", func() {
		specID := "c142042149f64911bb4698fb08572040"

		server := mockedServer(func(rw http.ResponseWriter, req *http.Request) {
			switch req.Method {
			case http.MethodPost:
				Expect(req.URL.Path).Should(Or(Equal("/v2/specs"), Equal("/sessions")))
				if req.URL.Path == "/sessions" {
					writeCookie(rw)
				} else {
					writeResponse(rw, http.StatusOK, Spec{
						Data: SpecData{ID: specID},
					})
				}
			case http.MethodGet:
				Expect(fmt.Sprintf("/v2/specs/%s", specID)).Should(Equal(req.URL.Path))
				writeResponse(rw, http.StatusOK, Response{
					Data: map[string]interface{}{},
				})
			case http.MethodDelete:
				Expect(fmt.Sprintf("/v2/specs/%s", specID)).Should(Equal(req.URL.Path))
				writeResponse(rw, http.StatusNoContent, nil)
			}
		})
		defer server.Close()

		c, err := newDefaultClient(server.URL)
		Expect(err).ShouldNot(HaveOccurred())
		c.SetClient(server.Client())

		s, err := c.CreateSpec(spec)
		Expect(err).ShouldNot(HaveOccurred())

		_, err = c.ReadSpec(s.Data.ID)
		Expect(err).ShouldNot(HaveOccurred())

		err = c.DeleteSpec(s.Data.ID)
		Expect(err).ShouldNot(HaveOccurred())
	})

	// Mocks the creation, read, delete cycle for Chainlink bridges
	It("can Create, Read, and Delete bridges", func() {
		bta := BridgeTypeAttributes{
			Name: "example",
			URL:  "https://example.com",
		}

		server := mockedServer(func(rw http.ResponseWriter, req *http.Request) {
			switch req.Method {
			case http.MethodPost:
				Expect(req.URL.Path).Should(Or(Equal("/v2/bridge_types"), Equal("/sessions")))
				if req.URL.Path == "/sessions" {
					writeCookie(rw)
				} else {
					writeResponse(rw, http.StatusOK, nil)
				}
			case http.MethodGet:
				Expect("/v2/bridge_types/example").Should(Equal(req.URL.Path))
				writeResponse(rw, http.StatusOK, BridgeType{
					Data: BridgeTypeData{
						Attributes: bta,
					},
				})
			case http.MethodDelete:
				Expect("/v2/bridge_types/example").Should(Equal(req.URL.Path))
				writeResponse(rw, http.StatusOK, nil)
			}
		})
		defer server.Close()

		c, err := newDefaultClient(server.URL)
		Expect(err).ShouldNot(HaveOccurred())
		c.SetClient(server.Client())

		err = c.CreateBridge(&bta)
		Expect(err).ShouldNot(HaveOccurred())

		bt, err := c.ReadBridge(bta.Name)
		Expect(err).ShouldNot(HaveOccurred())

		Expect(bt.Data.Attributes.Name).Should(Equal(bta.Name))
		Expect(bt.Data.Attributes.URL).Should(Equal(bta.URL))

		err = c.DeleteBridge(bta.Name)
		Expect(err).ShouldNot(HaveOccurred())
	})

	// Mocks the creation, read, delete cycle for OCR keys
	It("can Create, Read, and Delete OCR keys", func() {
		ocrKeyData := OCRKeyData{
			ID: "1",
			Attributes: OCRKeyAttributes{
				ConfigPublicKey:       "someNon3sens3",
				OffChainPublicKey:     "mor3Non3sens3",
				OnChainSigningAddress: "thisActuallyMak3sS3ns3",
			},
		}

		server := mockedServer(func(rw http.ResponseWriter, req *http.Request) {
			switch req.Method {
			case http.MethodPost:
				Expect(req.URL.Path).Should(Or(Equal("/v2/keys/ocr"), Equal("/sessions")))
				if req.URL.Path == "/sessions" {
					writeCookie(rw)
				} else {
					writeResponse(rw, http.StatusOK, OCRKey{ocrKeyData})
				}
			case http.MethodGet:
				Expect("/v2/keys/ocr").Should(Equal(req.URL.Path))
				writeResponse(rw, http.StatusOK, OCRKeys{
					Data: []OCRKeyData{ocrKeyData},
				})
			case http.MethodDelete:
				Expect("/v2/keys/ocr/1").Should(Equal(req.URL.Path))
				writeResponse(rw, http.StatusOK, nil)
			}
		})
		defer server.Close()

		c, err := newDefaultClient(server.URL)
		Expect(err).ShouldNot(HaveOccurred())
		c.SetClient(server.Client())

		receivedKey, err := c.CreateOCRKey()
		Expect(err).ShouldNot(HaveOccurred())

		keys, err := c.ReadOCRKeys()
		Expect(err).ShouldNot(HaveOccurred())
		Expect(keys.Data).Should(ContainElement(receivedKey.Data))

		err = c.DeleteOCRKey("1")
		Expect(err).ShouldNot(HaveOccurred())
	})

	// Mocks the creation, read, delete cycle for OCR keys
	It("can Create, Read, and Delete OCR2 keys", func() {
		for _, chain := range []string{"evm", "solana"} {
			ocrKeyData := OCR2KeyData{
				ID: "1",
				Attributes: OCR2KeyAttributes{
					ChainType:         chain,
					ConfigPublicKey:   "someNon3sens3",
					OffChainPublicKey: "mor3Non3sens3",
					OnChainPublicKey:  "thisActuallyMak3sS3ns3",
				},
			}
			server := mockedServer(func(rw http.ResponseWriter, req *http.Request) {
				switch req.Method {
				case http.MethodPost:
					Expect(req.URL.Path).Should(Or(Equal(fmt.Sprintf("/v2/keys/ocr2/%s", chain)), Equal("/sessions")))
					if req.URL.Path == "/sessions" {
						writeCookie(rw)
					} else {
						writeResponse(rw, http.StatusOK, OCR2Key{ocrKeyData})
					}
				case http.MethodGet:
					Expect("/v2/keys/ocr2").Should(Equal(req.URL.Path))
					writeResponse(rw, http.StatusOK, OCR2Keys{
						Data: []OCR2KeyData{ocrKeyData},
					})
				case http.MethodDelete:
					Expect("/v2/keys/ocr2/1").Should(Equal(req.URL.Path))
					writeResponse(rw, http.StatusOK, nil)
				}
			})
			defer server.Close()

			c, err := newDefaultClient(server.URL)
			Expect(err).ShouldNot(HaveOccurred())
			c.SetClient(server.Client())

			receivedKey, err := c.CreateOCR2Key(chain)
			Expect(err).ShouldNot(HaveOccurred())

			keys, err := c.ReadOCR2Keys()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(keys.Data).Should(ContainElement(receivedKey.Data))

			err = c.DeleteOCR2Key("1")
			Expect(err).ShouldNot(HaveOccurred())
		}
	})

	// Mocks the creation, read, delete cycle for P2P keys
	It("can Create, Read, and Delete P2P keys", func() {
		p2pKeyData := P2PKeyData{
			P2PKeyAttributes{
				ID:        1,
				PeerID:    "someNon3sens3",
				PublicKey: "mor3Non3sens3",
			},
		}

		server := mockedServer(func(rw http.ResponseWriter, req *http.Request) {
			switch req.Method {
			case http.MethodPost:
				Expect(req.URL.Path).Should(Or(Equal("/v2/keys/p2p"), Equal("/sessions")))
				if req.URL.Path == "/sessions" {
					writeCookie(rw)
				} else {
					writeResponse(rw, http.StatusOK, P2PKey{p2pKeyData})
				}
			case http.MethodGet:
				Expect("/v2/keys/p2p").Should(Equal(req.URL.Path))
				writeResponse(rw, http.StatusOK, P2PKeys{
					Data: []P2PKeyData{p2pKeyData},
				})
			case http.MethodDelete:
				Expect("/v2/keys/p2p/1").Should(Equal(req.URL.Path))
				writeResponse(rw, http.StatusOK, nil)
			}
		})
		defer server.Close()

		c, err := newDefaultClient(server.URL)
		Expect(err).ShouldNot(HaveOccurred())
		c.SetClient(server.Client())

		receivedKey, err := c.CreateP2PKey()
		Expect(err).ShouldNot(HaveOccurred())

		keys, err := c.ReadP2PKeys()
		Expect(err).ShouldNot(HaveOccurred())
		Expect(keys.Data).Should(ContainElement(receivedKey.Data))

		err = c.DeleteP2PKey(1)
		Expect(err).ShouldNot(HaveOccurred())
	})

	// Mocks the creation, read, delete cycle for ETH keys
	It("can Create, Read, and Delete ETH keys", func() {
		ethKeyData := ETHKeyData{
			Attributes: ETHKeyAttributes{
				Address: "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
			},
		}

		server := mockedServer(func(rw http.ResponseWriter, req *http.Request) {
			switch req.Method {
			case http.MethodPost:
				Expect(req.URL.Path).Should(Or(Equal("/v2/keys/eth"), Equal("/sessions")))
				if req.URL.Path == "/sessions" {
					writeCookie(rw)
				} else {
					writeResponse(rw, http.StatusOK, ETHKey{ethKeyData})
				}
			case http.MethodGet:
				Expect("/v2/keys/eth").Should(Equal(req.URL.Path))
				writeResponse(rw, http.StatusOK, ETHKeys{
					Data: []ETHKeyData{ethKeyData},
				})
			case http.MethodDelete:
				Expect("/v2/keys/eth/1").Should(Equal(req.URL.Path))
				writeResponse(rw, http.StatusOK, nil)
			}
		})
		defer server.Close()

		c, err := newDefaultClient(server.URL)
		Expect(err).ShouldNot(HaveOccurred())
		c.SetClient(server.Client())

		receivedKeys, err := c.ReadETHKeys()
		Expect(err).ShouldNot(HaveOccurred())
		Expect(receivedKeys.Data).Should(ContainElement(ethKeyData))
	})

	// Mocks the creation, read, delete cycle for Tx keys
	It("can Create, Read, and Delete Tx keys", func() {
		for _, chain := range []string{"solana"} {
			txKeyData := TxKeyData{
				Type: "encryptedKeyPlacholder",
				ID:   "someTestKeyID",
				Attributes: TxKeyAttributes{
					PublicKey: "aRandomTestPublicKeyForArbitraryChain",
				},
			}

			server := mockedServer(func(rw http.ResponseWriter, req *http.Request) {
				endpoint := fmt.Sprintf("/v2/keys/%s", chain)
				switch req.Method {
				case http.MethodPost:
					Expect(req.URL.Path).Should(Or(Equal(endpoint), Equal("/sessions")))
					if req.URL.Path == "/sessions" {
						writeCookie(rw)
					} else {
						writeResponse(rw, http.StatusOK, TxKey{txKeyData})
					}
				case http.MethodGet:
					Expect(endpoint).Should(Equal(req.URL.Path))
					writeResponse(rw, http.StatusOK, TxKeys{
						Data: []TxKeyData{txKeyData},
					})
				case http.MethodDelete:
					Expect(endpoint + "/1").Should(Equal(req.URL.Path))
					writeResponse(rw, http.StatusOK, nil)
				}
			})
			defer server.Close()

			c, err := newDefaultClient(server.URL)
			Expect(err).ShouldNot(HaveOccurred())
			c.SetClient(server.Client())

			receivedKey, err := c.CreateTxKey(chain)
			Expect(err).ShouldNot(HaveOccurred())

			keys, err := c.ReadTxKeys(chain)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(keys.Data).Should(ContainElement(receivedKey.Data))

			err = c.DeleteTxKey(chain, "1")
			Expect(err).ShouldNot(HaveOccurred())
		}
	})

	// Mocks the reading transactions and attempted transactions
	It("can read Transactions and Transaction Attempts from the chainlink node", func() {
		mockTxString := `{
	"data": [
		{
			"type": "transactions",
			"id": "0xd694f4a84b3aa8f1fae2443c2444760306eed5d575a7eb22eb64101511a3a5c0",
			"attributes": {
				"state": "confirmed",
				"data": "0x0",
				"from": "0x0",
				"gasLimit": "2650000",
				"gasPrice": "10000001603",
				"hash": "0xd694f4a84b3aa8f1fae2443c2444760306eed5d575a7eb22eb64101511a3a5c0",
				"rawHex": "0x0",
				"nonce": "1",
				"sentAt": "199",
				"to": "0x610178da211fef7d417bc0e6fed39f05609ad788",
				"value": "0.000000000000000000",
				"evmChainID": "1337"
			}
		}
	],
	"meta": {
		"count": 14
	}
}`
		mockTxData := TransactionsData{}
		err := json.Unmarshal([]byte(mockTxString), &mockTxData)
		Expect(err).ShouldNot(HaveOccurred())

		server := mockedServer(func(rw http.ResponseWriter, req *http.Request) {
			actualEndpoint := "/v2/transactions"
			attemptsEndpoint := "/v2/tx_attempts"
			switch req.Method {
			case http.MethodGet:
				Expect(req.URL.Path).Should(Or(Equal(actualEndpoint), Equal(attemptsEndpoint)))
				writeResponse(rw, http.StatusOK, mockTxData)
			case http.MethodPost:
				Expect(req.URL.Path).Should(Equal("/sessions"))
				writeCookie(rw)
			}
		})
		defer server.Close()

		c, err := newDefaultClient(server.URL)
		Expect(err).ShouldNot(HaveOccurred())
		c.SetClient(server.Client())

		_, err = c.ReadTransactionAttempts()
		Expect(err).ShouldNot(HaveOccurred())

		_, err = c.ReadTransactions()
		Expect(err).ShouldNot(HaveOccurred())
	})

	// Mocks the reading transactions and attempted transactions
	It("can send ETH transactions", func() {
		mockTxString := `{"data": {
			"type": "transactions",
			"id": "",
			"attributes": {
				"state": "in_progress",
				"data": "",
				"from": "",
				"gasLimit": "21000",
				"gasPrice": "",
				"hash": "",
				"rawHex": "",
				"nonce": "1",
				"sentAt": "199",
				"to": "0x610178da211fef7d417bc0e6fed39f05609ad788",
				"value": "0.000000000000000000",
				"evmChainID": "1337"
			}
		}}`
		mockTxData := SingleTransactionDataWrapper{}
		err := json.Unmarshal([]byte(mockTxString), &mockTxData)
		Expect(err).ShouldNot(HaveOccurred())

		server := mockedServer(func(rw http.ResponseWriter, req *http.Request) {
			switch req.Method {
			case http.MethodPost:
				if req.URL.Path == "/sessions" {
					writeCookie(rw)
				} else {
					writeResponse(rw, http.StatusOK, mockTxData)
				}
			}
		})
		defer server.Close()

		c, err := newDefaultClient(server.URL)
		Expect(err).ShouldNot(HaveOccurred())
		c.SetClient(server.Client())

		_, err = c.SendNativeToken(big.NewInt(1), "0x123", "0x420")
		Expect(err).ShouldNot(HaveOccurred())
	})

	// Mocks the creation, read cycle for CSA keys
	It("can Create and Read CSA keys", func() {
		csaKeyData := CSAKeyData{
			"csaKeys",
			"id",
			CSAKeyAttributes{
				PublicKey: "mor3Non3sens3",
				Version:   1,
			},
		}

		server := mockedServer(func(rw http.ResponseWriter, req *http.Request) {
			switch req.Method {
			case http.MethodPost:
				Expect(req.URL.Path).Should(Or(Equal("/v2/keys/csa"), Equal("/sessions")))
				if req.URL.Path == "/sessions" {
					writeCookie(rw)
				} else {
					writeResponse(rw, http.StatusOK, CSAKey{csaKeyData})
				}
			case http.MethodGet:
				Expect("/v2/keys/csa").Should(Equal(req.URL.Path))
				writeResponse(rw, http.StatusOK, CSAKeys{
					Data: []CSAKeyData{csaKeyData},
				})
			}
		})
		defer server.Close()

		c, err := newDefaultClient(server.URL)
		Expect(err).ShouldNot(HaveOccurred())
		c.SetClient(server.Client())

		receivedKey, err := c.CreateCSAKey()
		Expect(err).ShouldNot(HaveOccurred())

		keys, err := c.ReadCSAKeys()
		Expect(err).ShouldNot(HaveOccurred())
		Expect(keys.Data).Should(ContainElement(receivedKey.Data))
	})

	// Mocks the creation, read, delete cycle for Chainlink EIs
	It("can Create, Read, and Delete external initiators", func() {
		eia := EIAttributes{
			Name: "example",
			URL:  "https://example.com",
		}

		server := mockedServer(func(rw http.ResponseWriter, req *http.Request) {
			switch req.Method {
			case http.MethodPost:
				Expect(req.URL.Path).Should(Or(Equal("/v2/external_initiators"), Equal("/sessions")))
				if req.URL.Path == "/sessions" {
					writeCookie(rw)
				} else {
					writeResponse(rw, http.StatusCreated, EIKeyCreate{
						Data: EIKey{
							Attributes: eia,
						},
					})
				}
			case http.MethodGet:
				Expect("/v2/external_initiators").Should(Equal(req.URL.Path))
				writeResponse(rw, http.StatusOK, EIKeys{
					Data: []EIKey{
						{
							Attributes: eia,
						},
					},
				})
			case http.MethodDelete:
				Expect("/v2/external_initiators/example").Should(Equal(req.URL.Path))
				writeResponse(rw, http.StatusNoContent, nil)
			}
		})
		defer server.Close()

		c, err := newDefaultClient(server.URL)
		Expect(err).ShouldNot(HaveOccurred())
		c.SetClient(server.Client())

		ei, err := c.CreateEI(&eia)
		Expect(err).ShouldNot(HaveOccurred())

		Expect(ei.Data.Attributes.Name).Should(Equal(eia.Name))
		Expect(ei.Data.Attributes.URL).Should(Equal(eia.URL))

		eis, err := c.ReadEIs()
		Expect(err).ShouldNot(HaveOccurred())

		Expect(eis.Data[0].Attributes.Name).Should(Equal(eia.Name))
		Expect(eis.Data[0].Attributes.URL).Should(Equal(eia.URL))

		err = c.DeleteEI(eia.Name)
		Expect(err).ShouldNot(HaveOccurred())
	})

	// Mocks the creation, read cycle for chains
	It("can create chains", func() {
		terraAttr := TerraChainAttributes{
			ChainID: "chainId",
		}
		solAttr := SolanaChainAttributes{
			ChainID: "chainId",
		}

		server := mockedServer(func(rw http.ResponseWriter, req *http.Request) {
			switch req.Method {
			case http.MethodPost:
				switch req.URL.Path {
				case "/sessions":
					writeCookie(rw)
				case "/v2/chains/terra":
					writeResponse(rw, http.StatusCreated, TerraChainCreate{
						Data: TerraChain{
							Attributes: terraAttr,
						},
					})
				case "/v2/chains/solana":
					writeResponse(rw, http.StatusCreated, SolanaChainCreate{
						Data: SolanaChain{
							Attributes: solAttr,
						},
					})
				default:
					// error if unknown path
					Expect(errors.New("unknown path")).ShouldNot(HaveOccurred())
				}
			}
		})
		defer server.Close()

		c, err := newDefaultClient(server.URL)
		Expect(err).ShouldNot(HaveOccurred())
		c.SetClient(server.Client())

		resTerra, err := c.CreateTerraChain(&terraAttr)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(resTerra.Data.Attributes.ChainID).Should(Equal(terraAttr.ChainID))

		resSol, err := c.CreateSolanaChain(&solAttr)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(resSol.Data.Attributes.ChainID).Should(Equal(solAttr.ChainID))
	})

	// Mocks the creation, read cycle for nodes
	It("can create nodes", func() {
		terraAttr := TerraNodeAttributes{
			Name:          "name",
			TerraChainID:  "chainid",
			TendermintURL: "http://tendermint.com",
		}
		solAttr := SolanaNodeAttributes{
			Name:          "name",
			SolanaChainID: "chainid",
			SolanaURL:     "http://solana.com",
		}

		server := mockedServer(func(rw http.ResponseWriter, req *http.Request) {
			switch req.Method {
			case http.MethodPost:
				switch req.URL.Path {
				case "/sessions":
					writeCookie(rw)
				case "/v2/nodes/terra":
					writeResponse(rw, http.StatusOK, TerraNodeCreate{
						Data: TerraNode{
							Attributes: terraAttr,
						},
					})
				case "/v2/nodes/solana":
					writeResponse(rw, http.StatusOK, SolanaNodeCreate{
						Data: SolanaNode{
							Attributes: solAttr,
						},
					})
				default:
					// error if unknown path
					Expect(errors.New("unknown path")).ShouldNot(HaveOccurred())
				}
			}
		})
		defer server.Close()

		c, err := newDefaultClient(server.URL)
		Expect(err).ShouldNot(HaveOccurred())
		c.SetClient(server.Client())

		resTerra, err := c.CreateTerraNode(&terraAttr)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(resTerra.Data.Attributes.Name).Should(Equal(terraAttr.Name))
		Expect(resTerra.Data.Attributes.TerraChainID).Should(Equal(terraAttr.TerraChainID))
		Expect(resTerra.Data.Attributes.TendermintURL).Should(Equal(terraAttr.TendermintURL))

		resSol, err := c.CreateSolanaNode(&solAttr)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(resSol.Data.Attributes.Name).Should(Equal(solAttr.Name))
		Expect(resSol.Data.Attributes.SolanaChainID).Should(Equal(solAttr.SolanaChainID))
		Expect(resSol.Data.Attributes.SolanaURL).Should(Equal(solAttr.SolanaURL))
	})
})

func newDefaultClient(url string) (Chainlink, error) {
	cl, err := NewChainlink(&ChainlinkConfig{
		Email:    "admin@node.local",
		Password: "twochains",
		URL:      url,
	}, nil)
	return cl, err
}

func mockedServer(handlerFunc http.HandlerFunc) *httptest.Server {
	return httptest.NewServer(handlerFunc)
}

func writeCookie(rw http.ResponseWriter) {
	cookie := http.Cookie{
		Name:    "clsession",
		Value:   "something",
		Expires: time.Now().Add(time.Minute * 5),
	}
	http.SetCookie(rw, &cookie)
}

func writeResponse(rw http.ResponseWriter, statusCode int, obj interface{}) {
	rw.WriteHeader(statusCode)
	if obj == nil {
		return
	}
	b, err := json.Marshal(obj)
	Expect(err).ShouldNot(HaveOccurred())
	_, err = rw.Write(b)
	Expect(err).ShouldNot(HaveOccurred())
}

package client

//revive:disable:defer
import (
	"encoding/json"
	"fmt"
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

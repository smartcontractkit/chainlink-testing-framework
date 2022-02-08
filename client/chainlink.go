package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"strings"

	"github.com/rs/zerolog/log"
)

// OneLINK representation of a single LINK token
var OneLINK = big.NewFloat(1e18)

// Chainlink interface that enables interactions with a chainlink node
type Chainlink interface {
	URL() string
	CreateJob(spec JobSpec) (*Job, error)
	CreateJobRaw(spec string) (*Job, error)
	ReadJobs() (*ResponseSlice, error)
	ReadJob(id string) (*Response, error)
	DeleteJob(id string) error

	CreateSpec(spec string) (*Spec, error)
	ReadSpec(id string) (*Response, error)
	DeleteSpec(id string) error

	CreateBridge(bta *BridgeTypeAttributes) error
	ReadBridge(name string) (*BridgeType, error)
	DeleteBridge(name string) error

	ReadRunsByJob(jobID string) (*JobRunsResponse, error)

	CreateOCRKey() (*OCRKey, error)
	ReadOCRKeys() (*OCRKeys, error)
	DeleteOCRKey(id string) error

	CreateOCR2Key(chain string) (*OCR2Key, error)
	ReadOCR2Keys() (*OCR2Keys, error)
	DeleteOCR2Key(id string) error

	CreateP2PKey() (*P2PKey, error)
	ReadP2PKeys() (*P2PKeys, error)
	DeleteP2PKey(id int) error

	ReadETHKeys() (*ETHKeys, error)
	PrimaryEthAddress() (string, error)

	CreateTxKey(chain string) (*TxKey, error)
	ReadTxKeys(chain string) (*TxKeys, error)
	DeleteTxKey(chain, id string) error

	ReadTransactionAttempts() (*TransactionsData, error)
	ReadTransactions() (*TransactionsData, error)
	SendNativeToken(amount *big.Int, fromAddress, toAddress string) (interface{}, error)

	CreateVRFKey() (*VRFKey, error)
	ReadVRFKeys() (*VRFKeys, error)

	CreateCSAKey() (*CSAKey, error)
	ReadCSAKeys() (*CSAKeys, error)

	CreateEI(eia *EIAttributes) (*EIKeyCreate, error)
	ReadEIs() (*EIKeys, error)
	DeleteEI(name string) error

	CreateTerraChain(node *TerraChainAttributes) (*TerraChainCreate, error)
	CreateTerraNode(node *TerraNodeAttributes) (*TerraNodeCreate, error)

	RemoteIP() string
	SetSessionCookie() error

	SetPageSize(size int)

	// SetClient is used for testing
	SetClient(client *http.Client)
}

type chainlink struct {
	*BasicHTTPClient
	Config            *ChainlinkConfig
	pageSize          int
	primaryEthAddress string
}

// NewChainlink creates a new chainlink model using a provided config
func NewChainlink(c *ChainlinkConfig, httpClient *http.Client) (Chainlink, error) {
	cl := &chainlink{
		Config:          c,
		BasicHTTPClient: NewBasicHTTPClient(httpClient, c.URL),
		pageSize:        25,
	}
	return cl, cl.SetSessionCookie()
}

// URL chainlink instance http url
func (c *chainlink) URL() string {
	return c.Config.URL
}

// CreateJobRaw creates a Chainlink job based on the provided spec string
func (c *chainlink) CreateJobRaw(spec string) (*Job, error) {
	job := &Job{}
	log.Info().Str("Node URL", c.Config.URL).Str("Job Body", spec).Msg("Creating Job")
	_, err := c.do(http.MethodPost, "/v2/jobs", &JobForm{
		TOML: spec,
	}, &job, http.StatusOK)
	return job, err
}

// CreateJob creates a Chainlink job based on the provided spec struct
func (c *chainlink) CreateJob(spec JobSpec) (*Job, error) {
	job := &Job{}
	specString, err := spec.String()
	if err != nil {
		return nil, err
	}
	log.Info().Str("Node URL", c.Config.URL).Str("Type", spec.Type()).Msg("Creating Job")
	_, err = c.do(http.MethodPost, "/v2/jobs", &JobForm{
		TOML: specString,
	}, &job, http.StatusOK)
	return job, err
}

// ReadJobs reads all jobs from the Chainlink node
func (c *chainlink) ReadJobs() (*ResponseSlice, error) {
	specObj := &ResponseSlice{}
	log.Info().Str("Node URL", c.Config.URL).Msg("Getting Jobs")
	_, err := c.do(http.MethodGet, "/v2/jobs", nil, specObj, http.StatusOK)
	return specObj, err
}

// ReadJob reads a job with the provided ID from the Chainlink node
func (c *chainlink) ReadJob(id string) (*Response, error) {
	specObj := &Response{}
	log.Info().Str("Node URL", c.Config.URL).Str("ID", id).Msg("Reading Job")
	_, err := c.do(http.MethodGet, fmt.Sprintf("/v2/jobs/%s", id), nil, specObj, http.StatusOK)
	return specObj, err
}

// DeleteJob deletes a job with a provided ID from the Chainlink node
func (c *chainlink) DeleteJob(id string) error {
	log.Info().Str("Node URL", c.Config.URL).Str("ID", id).Msg("Deleting Job")
	_, err := c.do(http.MethodDelete, fmt.Sprintf("/v2/jobs/%s", id), nil, nil, http.StatusNoContent)
	return err
}

// CreateSpec creates a job spec on the Chainlink node
func (c *chainlink) CreateSpec(spec string) (*Spec, error) {
	s := &Spec{}
	r := strings.NewReplacer("\n", "", " ", "", "\\", "") // Makes it more compact and readable for logging
	log.Info().Str("Node URL", c.Config.URL).Str("Spec", r.Replace(spec)).Msg("Creating Spec")
	_, err := c.doRaw(http.MethodPost, "/v2/specs", []byte(spec), s, http.StatusOK)
	return s, err
}

// ReadSpec reads a job spec with the provided ID on the Chainlink node
func (c *chainlink) ReadSpec(id string) (*Response, error) {
	specObj := &Response{}
	log.Info().Str("Node URL", c.Config.URL).Str("ID", id).Msg("Reading Spec")
	_, err := c.do(http.MethodGet, fmt.Sprintf("/v2/specs/%s", id), nil, specObj, http.StatusOK)
	return specObj, err
}

// ReadRunsByJob reads all runs for a job
func (c *chainlink) ReadRunsByJob(jobID string) (*JobRunsResponse, error) {
	runsObj := &JobRunsResponse{}
	log.Debug().Str("Node URL", c.Config.URL).Str("JobID", jobID).Msg("Reading runs for a job")
	_, err := c.do(http.MethodGet, fmt.Sprintf("/v2/jobs/%s/runs", jobID), nil, runsObj, http.StatusOK)
	return runsObj, err
}

// DeleteSpec deletes a job spec with the provided ID from the Chainlink node
func (c *chainlink) DeleteSpec(id string) error {
	log.Info().Str("Node URL", c.Config.URL).Str("ID", id).Msg("Deleting Spec")
	_, err := c.do(http.MethodDelete, fmt.Sprintf("/v2/specs/%s", id), nil, nil, http.StatusNoContent)
	return err
}

// CreateBridge creates a bridge on the Chainlink node based on the provided attributes
func (c *chainlink) CreateBridge(bta *BridgeTypeAttributes) error {
	log.Info().Str("Node URL", c.Config.URL).Str("Name", bta.Name).Msg("Creating Bridge")
	_, err := c.do(http.MethodPost, "/v2/bridge_types", bta, nil, http.StatusOK)
	return err
}

// ReadBridge reads a bridge from the Chainlink node based on the provided name
func (c *chainlink) ReadBridge(name string) (*BridgeType, error) {
	bt := BridgeType{}
	log.Info().Str("Node URL", c.Config.URL).Str("Name", name).Msg("Reading Bridge")
	_, err := c.do(http.MethodGet, fmt.Sprintf("/v2/bridge_types/%s", name), nil, &bt, http.StatusOK)
	return &bt, err
}

// DeleteBridge deletes a bridge on the Chainlink node based on the provided name
func (c *chainlink) DeleteBridge(name string) error {
	log.Info().Str("Node URL", c.Config.URL).Str("Name", name).Msg("Deleting Bridge")
	_, err := c.do(http.MethodDelete, fmt.Sprintf("/v2/bridge_types/%s", name), nil, nil, http.StatusOK)
	return err
}

// CreateOCRKey creates an OCRKey on the Chainlink node
func (c *chainlink) CreateOCRKey() (*OCRKey, error) {
	ocrKey := &OCRKey{}
	log.Info().Str("Node URL", c.Config.URL).Msg("Creating OCR Key")
	_, err := c.do(http.MethodPost, "/v2/keys/ocr", nil, ocrKey, http.StatusOK)
	return ocrKey, err
}

// ReadOCRKeys reads all OCRKeys from the Chainlink node
func (c *chainlink) ReadOCRKeys() (*OCRKeys, error) {
	ocrKeys := &OCRKeys{}
	log.Info().Str("Node URL", c.Config.URL).Msg("Reading OCR Keys")
	_, err := c.do(http.MethodGet, "/v2/keys/ocr", nil, ocrKeys, http.StatusOK)
	for index := range ocrKeys.Data {
		ocrKeys.Data[index].Attributes.ConfigPublicKey = strings.TrimPrefix(
			ocrKeys.Data[index].Attributes.ConfigPublicKey, "ocrcfg_")
		ocrKeys.Data[index].Attributes.OffChainPublicKey = strings.TrimPrefix(
			ocrKeys.Data[index].Attributes.OffChainPublicKey, "ocroff_")
		ocrKeys.Data[index].Attributes.OnChainSigningAddress = strings.TrimPrefix(
			ocrKeys.Data[index].Attributes.OnChainSigningAddress, "ocrsad_")
	}
	return ocrKeys, err
}

// DeleteOCRKey deletes an OCRKey based on the provided ID
func (c *chainlink) DeleteOCRKey(id string) error {
	log.Info().Str("Node URL", c.Config.URL).Str("ID", id).Msg("Deleting OCR Key")
	_, err := c.do(http.MethodDelete, fmt.Sprintf("/v2/keys/ocr/%s", id), nil, nil, http.StatusOK)
	return err
}

// CreateOCR2Key creates an OCR2Key on the Chainlink node
func (c *chainlink) CreateOCR2Key(chain string) (*OCR2Key, error) {
	ocr2Key := &OCR2Key{}
	log.Info().Str("Node URL", c.Config.URL).Msg("Creating OCR2 Key")
	_, err := c.do(http.MethodPost, fmt.Sprintf("/v2/keys/ocr2/%s", chain), nil, ocr2Key, http.StatusOK)
	return ocr2Key, err
}

// ReadOCR2Keys reads all OCR2Keys from the Chainlink node
func (c *chainlink) ReadOCR2Keys() (*OCR2Keys, error) {
	ocr2Keys := &OCR2Keys{}
	log.Info().Str("Node URL", c.Config.URL).Msg("Reading OCR2 Keys")
	_, err := c.do(http.MethodGet, "/v2/keys/ocr2", nil, ocr2Keys, http.StatusOK)
	return ocr2Keys, err
}

// DeleteOCR2Key deletes an OCR2Key based on the provided ID
func (c *chainlink) DeleteOCR2Key(id string) error {
	log.Info().Str("Node URL", c.Config.URL).Str("ID", id).Msg("Deleting OCR2 Key")
	_, err := c.do(http.MethodDelete, fmt.Sprintf("/v2/keys/ocr2/%s", id), nil, nil, http.StatusOK)
	return err
}

// CreateP2PKey creates an P2PKey on the Chainlink node
func (c *chainlink) CreateP2PKey() (*P2PKey, error) {
	p2pKey := &P2PKey{}
	log.Info().Str("Node URL", c.Config.URL).Msg("Creating P2P Key")
	_, err := c.do(http.MethodPost, "/v2/keys/p2p", nil, p2pKey, http.StatusOK)
	return p2pKey, err
}

// ReadP2PKeys reads all P2PKeys from the Chainlink node
func (c *chainlink) ReadP2PKeys() (*P2PKeys, error) {
	p2pKeys := &P2PKeys{}
	log.Info().Str("Node URL", c.Config.URL).Msg("Reading P2P Keys")
	_, err := c.do(http.MethodGet, "/v2/keys/p2p", nil, p2pKeys, http.StatusOK)
	if len(p2pKeys.Data) == 0 {
		err = fmt.Errorf("Found no P2P Keys on the chainlink node. Node URL: %s", c.Config.URL)
		log.Err(err).Msg("Error getting P2P keys")
		return nil, err
	}
	for index := range p2pKeys.Data {
		p2pKeys.Data[index].Attributes.PeerID = strings.TrimPrefix(p2pKeys.Data[index].Attributes.PeerID, "p2p_")
	}
	return p2pKeys, err
}

// DeleteP2PKey deletes a P2PKey on the Chainlink node based on the provided ID
func (c *chainlink) DeleteP2PKey(id int) error {
	log.Info().Str("Node URL", c.Config.URL).Int("ID", id).Msg("Deleting P2P Key")
	_, err := c.do(http.MethodDelete, fmt.Sprintf("/v2/keys/p2p/%d", id), nil, nil, http.StatusOK)
	return err
}

// ReadETHKeys reads all ETH keys from the Chainlink node
func (c *chainlink) ReadETHKeys() (*ETHKeys, error) {
	ethKeys := &ETHKeys{}
	log.Info().Str("Node URL", c.Config.URL).Msg("Reading ETH Keys")
	_, err := c.do(http.MethodGet, "/v2/keys/eth", nil, ethKeys, http.StatusOK)
	if len(ethKeys.Data) == 0 {
		log.Warn().Str("Node URL", c.Config.URL).Msg("Found no ETH Keys on the node")
	}
	return ethKeys, err
}

// CreateOCR2Key creates an OCR2Key on the Chainlink node
func (c *chainlink) CreateTxKey(chain string) (*TxKey, error) {
	txKey := &TxKey{}
	log.Info().Str("Node URL", c.Config.URL).Msg("Creating Tx Key")
	_, err := c.do(http.MethodPost, fmt.Sprintf("/v2/keys/%s", chain), nil, txKey, http.StatusOK)
	return txKey, err
}

// ReadOCR2Keys reads all OCR2Keys from the Chainlink node
func (c *chainlink) ReadTxKeys(chain string) (*TxKeys, error) {
	txKeys := &TxKeys{}
	log.Info().Str("Node URL", c.Config.URL).Msg("Reading Tx Keys")
	_, err := c.do(http.MethodGet, fmt.Sprintf("/v2/keys/%s", chain), nil, txKeys, http.StatusOK)
	return txKeys, err
}

// DeleteOCR2Key deletes an OCR2Key based on the provided ID
func (c *chainlink) DeleteTxKey(chain string, id string) error {
	log.Info().Str("Node URL", c.Config.URL).Str("ID", id).Msg("Deleting Tx Key")
	_, err := c.do(http.MethodDelete, fmt.Sprintf("/v2/keys/%s/%s", chain, id), nil, nil, http.StatusOK)
	return err
}

// ReadTransactionAttempts reads all transaction attempts on the chainlink node
func (c *chainlink) ReadTransactionAttempts() (*TransactionsData, error) {
	txsData := &TransactionsData{}
	log.Info().Str("Node URL", c.Config.URL).Msg("Reading Transaction Attempts")
	_, err := c.do(http.MethodGet, "/v2/tx_attempts", nil, txsData, http.StatusOK)
	return txsData, err
}

// ReadTransactions reads all transactions made by the chainlink node
func (c *chainlink) ReadTransactions() (*TransactionsData, error) {
	txsData := &TransactionsData{}
	log.Info().Str("Node URL", c.Config.URL).Msg("Reading Transactions")
	_, err := c.do(http.MethodGet, "/v2/transactions", nil, txsData, http.StatusOK)
	return txsData, err
}

// SendNativeToken sends native token (ETH usually) of a specified amount from one of its addresses to the target address
func (c *chainlink) SendNativeToken(amount *big.Int, fromAddress, toAddress string) (interface{}, error) {
	request := SendEtherRequest{
		DestinationAddress: toAddress,
		FromAddress:        fromAddress,
		Amount:             amount.String(),
		AllowHigherAmounts: true,
	}
	var ret interface{}
	log.Info().
		Str("Node URL", c.Config.URL).
		Str("From", fromAddress).
		Str("To", toAddress).
		Int64("Amount", amount.Int64()).
		Msg("Sending Native Token")
	_, err := c.do(http.MethodPost, "/v2/transfers", request, ret, http.StatusOK)
	return ret, err
}

// ReadVRFKeys reads all VRF keys from the Chainlink node
func (c *chainlink) ReadVRFKeys() (*VRFKeys, error) {
	vrfKeys := &VRFKeys{}
	log.Info().Str("Node URL", c.Config.URL).Msg("Reading VRF Keys")
	_, err := c.do(http.MethodGet, "/v2/keys/vrf", nil, vrfKeys, http.StatusOK)
	if len(vrfKeys.Data) == 0 {
		log.Warn().Str("Node URL", c.Config.URL).Msg("Found no VRF Keys on the node")
	}
	return vrfKeys, err
}

// CreateVRFKey creates a VRF key on the Chainlink node
func (c *chainlink) CreateVRFKey() (*VRFKey, error) {
	vrfKey := &VRFKey{}
	log.Info().Str("Node URL", c.Config.URL).Msg("Creating VRF Key")
	_, err := c.do(http.MethodPost, "/v2/keys/vrf", nil, vrfKey, http.StatusOK)
	return vrfKey, err
}

// CreateCSAKey creates a CSA key on the Chainlink node, only 1 CSA key per noe
func (c *chainlink) CreateCSAKey() (*CSAKey, error) {
	csaKey := &CSAKey{}
	log.Info().Str("Node URL", c.Config.URL).Msg("Creating CSA Key")
	_, err := c.do(http.MethodPost, "/v2/keys/csa", nil, csaKey, http.StatusOK)
	return csaKey, err
}

// ReadCSAKeys reads CSA keys from the Chainlink node
func (c *chainlink) ReadCSAKeys() (*CSAKeys, error) {
	csaKeys := &CSAKeys{}
	log.Info().Str("Node URL", c.Config.URL).Msg("Reading CSA Keys")
	_, err := c.do(http.MethodGet, "/v2/keys/csa", nil, csaKeys, http.StatusOK)
	if len(csaKeys.Data) == 0 {
		log.Warn().Str("Node URL", c.Config.URL).Msg("Found no CSA Keys on the node")
	}
	return csaKeys, err
}

// PrimaryEthAddress returns the primary ETH address for the chainlink node
func (c *chainlink) PrimaryEthAddress() (string, error) {
	if c.primaryEthAddress == "" {
		ethKeys, err := c.ReadETHKeys()
		if err != nil {
			return "", err
		}
		c.primaryEthAddress = ethKeys.Data[0].Attributes.Address
	}
	return c.primaryEthAddress, nil
}

// CreateEI creates an EI on the Chainlink node based on the provided attributes and returns the respective secrets
func (c *chainlink) CreateEI(eia *EIAttributes) (*EIKeyCreate, error) {
	ei := EIKeyCreate{}
	log.Info().Str("Node URL", c.Config.URL).Str("Name", eia.Name).Msg("Creating External Initiator")
	_, err := c.do(http.MethodPost, "/v2/external_initiators", eia, &ei, http.StatusCreated)
	return &ei, err
}

// ReadEIs reads all of the configured EIs from the chainlink node
func (c *chainlink) ReadEIs() (*EIKeys, error) {
	ei := EIKeys{}
	log.Info().Str("Node URL", c.Config.URL).Msg("Reading EI Keys")
	_, err := c.do(http.MethodGet, "/v2/external_initiators", nil, &ei, http.StatusOK)
	return &ei, err
}

// DeleteEI deletes an external initiator in the Chainlink node based on the provided name
func (c *chainlink) DeleteEI(name string) error {
	log.Info().Str("Node URL", c.Config.URL).Str("Name", name).Msg("Deleting EI")
	_, err := c.do(http.MethodDelete, fmt.Sprintf("/v2/external_initiators/%s", name), nil, nil, http.StatusNoContent)
	return err
}

// CreateTerraChain creates a terra chain
func (c *chainlink) CreateTerraChain(chain *TerraChainAttributes) (*TerraChainCreate, error) {
	response := TerraChainCreate{}
	log.Info().Str("Node URL", c.Config.URL).Str("Chain ID", chain.ChainID).Msg("Creating Terra Chain")
	_, err := c.do(http.MethodPost, "/v2/chains/terra", chain, &response, http.StatusCreated)
	return &response, err
}

// CreateTerraNode creates a terra node
func (c *chainlink) CreateTerraNode(node *TerraNodeAttributes) (*TerraNodeCreate, error) {
	response := TerraNodeCreate{}
	log.Info().Str("Node URL", c.Config.URL).Str("Name", node.Name).Msg("Creating Terra Node")
	_, err := c.do(http.MethodPost, "/v2/nodes/terra", node, &response, http.StatusOK)
	return &response, err
}

// RemoteIP retrieves the inter-cluster IP of the chainlink node, for use with inter-node communications
func (c *chainlink) RemoteIP() string {
	return c.Config.RemoteIP
}

// SetSessionCookie authenticates against the Chainlink node and stores the cookie in client state
func (c *chainlink) SetSessionCookie() error {
	session := &Session{Email: c.Config.Email, Password: c.Config.Password}
	b, err := json.Marshal(session)
	if err != nil {
		return err
	}
	resp, err := http.Post(
		fmt.Sprintf("%s/sessions", c.Config.URL),
		"application/json",
		bytes.NewReader(b),
	)
	if err != nil {
		return err
	}
	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf(
			"error while reading response: %v\nURL: %s\nresponse received: %s",
			err,
			c.Config.URL,
			string(b),
		)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf(
			"status code of %d was returned when trying to get a session\nURL: %s\nresponse received: %s",
			resp.StatusCode,
			c.Config.URL,
			b,
		)
	}
	if len(resp.Cookies()) == 0 {
		return fmt.Errorf("no cookie was returned after getting a session")
	}
	c.BasicHTTPClient.Cookies = resp.Cookies()

	sessionFound := false
	for _, cookie := range resp.Cookies() {
		if cookie.Name == "clsession" {
			sessionFound = true
		}
	}
	if !sessionFound {
		return fmt.Errorf("chainlink: session cookie wasn't returned on login")
	}
	return nil
}

// SetClient overrides the http client, used for mocking out the Chainlink server for unit testing
func (c *chainlink) SetClient(client *http.Client) {
	c.HttpClient = client
}

// SetPageSize globally sets the page
func (c *chainlink) SetPageSize(size int) {
	c.pageSize = size
}

func (c *chainlink) doRaw(
	method,
	endpoint string,
	body []byte, obj interface{},
	expectedStatusCode int,
) (*http.Response, error) {
	client := c.HttpClient

	req, err := http.NewRequest(
		method,
		fmt.Sprintf("%s%s", c.Config.URL, endpoint),
		bytes.NewBuffer(body),
	)
	if err != nil {
		return nil, err
	}
	for _, cookie := range c.Cookies {
		req.AddCookie(cookie)
	}

	q := req.URL.Query()
	q.Add("size", fmt.Sprint(c.pageSize))
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return resp, err
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf(
			"error while reading response: %v\nURL: %s\nresponse received: %s",
			err,
			c.Config.URL,
			string(b),
		)
	}
	if resp.StatusCode == http.StatusNotFound {
		return resp, ErrNotFound
	} else if resp.StatusCode == http.StatusUnprocessableEntity {
		return resp, ErrUnprocessableEntity
	} else if resp.StatusCode != expectedStatusCode {
		return resp, fmt.Errorf(
			"unexpected response code, got %d, expected %d\nURL: %s\nresponse received: %s",
			resp.StatusCode,
			expectedStatusCode,
			c.Config.URL,
			string(b),
		)
	}

	if obj == nil {
		return resp, err
	}
	err = json.Unmarshal(b, &obj)
	if err != nil {
		return nil, fmt.Errorf(
			"error while unmarshaling response: %v\nURL: %s\nresponse received: %s",
			err,
			c.Config.URL,
			string(b),
		)
	}
	return resp, err
}

func (c *chainlink) do(
	method,
	endpoint string,
	body interface{},
	obj interface{},
	expectedStatusCode int,
) (*http.Response, error) {
	b, err := json.Marshal(body)
	if body != nil && err != nil {
		return nil, err
	}
	return c.doRaw(method, endpoint, b, obj, expectedStatusCode)
}

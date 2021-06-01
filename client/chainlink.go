package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"integrations-framework/contracts/ethereum"
	"io/ioutil"
	"math/big"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
)

var ErrNotFound = errors.New("unexpected response code, got 404")
var ErrUnprocessableEntity = errors.New("unexpected response code, got 422")

// Chainlink interface that enables interactions with a chainlink node
type Chainlink interface {
	// Fund sends specified currencies to the contract
	Fund(fromWallet BlockchainWallet, ethAmount *big.Int, linkAmount *big.Int) error

	CreateJob(spec string) (*Job, error)
	ReadJob(id string) error
	DeleteJob(id string) error

	CreateSpec(spec string) (*Spec, error)
	ReadSpec(id string) (*Response, error)
	DeleteSpec(id string) error

	CreateBridge(bta *BridgeTypeAttributes) error
	ReadBridge(name string) (*BridgeType, error)
	DeleteBridge(name string) error

	CreateOCRKey() (*OCRKey, error)
	ReadOCRKeys() (*OCRKeys, error)
	DeleteOCRKey(id string) error

	CreateP2PKey() (*P2PKey, error)
	ReadP2PKeys() (*P2PKeys, error)
	DeleteP2PKey(id int) error

	ReadETHKeys() (*ETHKeys, error)

	SetSessionCookie() error

	// Used for testing
	SetClient(client *http.Client)
}

type chainlink struct {
	EthClient  *EthereumClient
	HttpClient *http.Client
	Config     *ChainlinkConfig
	Cookies    []*http.Cookie
}

// NewChainlink creates a new chainlink model using a provided config
func NewChainlink(c *ChainlinkConfig, ethClient *EthereumClient) (Chainlink, error) {
	cl := &chainlink{Config: c, EthClient: ethClient}
	return cl, cl.SetSessionCookie()
}

// CreateTemplateNodes lauches 5 chainlink nodes in a default config for testing
func CreateTemplateNodes(ethClient *EthereumClient, linkAddress string) ([]Chainlink, error) {
	err := CleanTemplateNodes()
	if err != nil {
		return nil, err
	}
	urlBase := "http://localhost:"
	email := "notreal@fakeemail.ch"
	pass := "twochains"
	ports := []string{"6711", "6722", "6733", "6744", "6755"}

	p, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		return nil, err
	}
	projectDir := strings.TrimSpace(string(p))

	log.Info().Str("CMD", "docker-compose -f docker-compose.yml up").Msg("Running command")
	cmd := exec.Command("docker-compose", "-f", "docker-compose.yml", "up")
	cmd.Dir = string(projectDir) + "/tools/chainlink_nodes"
	ethUrl := "ETH_URL=" + ethClient.Network.URL()
	if ethClient.Network.ID() == EthereumHardhatID {
		ethUrl = "ETH_URL=ws://host.docker.internal:8545"
	}
	chainId := "ETH_CHAIN_ID=" + ethClient.Network.ChainID().String()
	la := "LINK_CONTRACT_ADDRESS=" + linkAddress
	cmd.Env = []string{ethUrl, chainId, la}
	var e bytes.Buffer
	cmd.Stderr = &e
	err = cmd.Start()
	if err != nil {
		return nil, err
	}

	// Wait for Docker Compose to be up and healthy
	for _, port := range ports {
		resp, err := http.Get(urlBase + port)
		for start := time.Now(); time.Since(start) < 2*time.Minute; time.Sleep(time.Second * 3) {
			resp, err = http.Get(urlBase + port)
			if err == nil && resp.StatusCode == 200 {
				break
			}
		}
		if err != nil {
			log.Error().Str("ERROR", e.String()).Msg("STDERR")
			if resp != nil {
				log.Info().Int("Status Code", resp.StatusCode).Msg("Monitor Response")
			} else {
				log.Error().Msg("Monitor Response NIL")
			}
			_ = CleanTemplateNodes()
			return nil, err
		}
		log.Info().Str("URL", urlBase+port).Msg("Chainlink Node Healthy")
	}

	var cls []Chainlink
	for _, port := range ports {
		c := &ChainlinkConfig{
			URL:      urlBase + port,
			Email:    email,
			Password: pass,
		}
		cl, err := NewChainlink(c, ethClient)
		if err != nil {
			return nil, err
		}
		cl.SetClient(http.DefaultClient)
		cls = append(cls, cl)
	}

	return cls, err
}

// CleanTemplateNodes cleans the default setup for chainlink nodes
func CleanTemplateNodes() error {
	p, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		return err
	}
	projectDir := strings.TrimSpace(string(p))
	cmd := exec.Command("docker-compose", "-f", "./docker-compose.yml", "down", "-v", "--remove-orphans")
	cmd.Dir = string(projectDir) + "/tools/chainlink_nodes"

	log.Info().Str("CMD", "docker-compose -f ./docker-compose.yml down -v --remove-orphans").Msg("Running command")
	return cmd.Run()
}

// CreateJob creates a Chainlink job based on the provided spec string
func (c *chainlink) CreateJob(spec string) (*Job, error) {
	job := &Job{}
	log.Info().Str("Chainlink Node URL", c.Config.URL).Msg("Creating Job")
	_, err := c.do(http.MethodPost, "/v2/jobs", &JobForm{
		TOML: spec,
	}, &job, http.StatusOK)
	return job, err
}

// Fund sends specified currencies to the contract
func (c *chainlink) Fund(fromWallet BlockchainWallet, ethAmount, linkAmount *big.Int) error {
	ethKeys, err := c.ReadETHKeys()
	if err != nil {
		return err
	}
	toAddress := ethKeys.Data[0].Attributes.Address
	// Send ETH if not 0
	if ethAmount != nil && big.NewInt(0).Cmp(ethAmount) != 0 {
		log.Info().
			Str("Token", "ETH").
			Str("From", fromWallet.Address()).
			Str("To", toAddress).
			Str("Amount", ethAmount.String()).
			Str("Node URL", c.Config.URL).
			Msg("Funding Chainlink Node")
		_, err := c.EthClient.SendTransaction(fromWallet, common.HexToAddress(toAddress), ethAmount, nil)
		if err != nil {
			return err
		}
	}

	// Send LINK if not 0
	if linkAmount != nil && big.NewInt(0).Cmp(linkAmount) != 0 {
		// Prepare data field for token tx
		log.Info().
			Str("Token", "LINK").
			Str("From", fromWallet.Address()).
			Str("To", toAddress).
			Str("Amount", linkAmount.String()).
			Str("Node URL", c.Config.URL).
			Msg("Funding Chainlink Node")
		linkAddress := common.HexToAddress(c.EthClient.Network.Config().LinkTokenAddress)
		linkInstance, err := ethereum.NewLinkToken(linkAddress, c.EthClient.Client)
		if err != nil {
			return err
		}
		opts, err := c.EthClient.TransactionOpts(fromWallet, common.HexToAddress(toAddress), nil, nil)
		if err != nil {
			return err
		}
		tx, err := linkInstance.Transfer(opts, common.HexToAddress(toAddress), linkAmount)
		if err != nil {
			return err
		}

		err = c.EthClient.WaitForTransaction(tx.Hash())
		if err != nil {
			return err
		}
	}
	return nil
}

// ReadJob reads a job with the provided ID from the Chainlink node
func (c *chainlink) ReadJob(id string) error {
	_, err := c.do(http.MethodGet, fmt.Sprintf("/v2/jobs/%s", id), nil, nil, http.StatusOK)
	return err
}

// DeleteJob deletes a job with a provided ID from the Chainlink node
func (c *chainlink) DeleteJob(id string) error {
	_, err := c.do(http.MethodDelete, fmt.Sprintf("/v2/jobs/%s", id), nil, nil, http.StatusNoContent)
	return err
}

// CreateSpec creates a job spec on the Chainlink node
func (c *chainlink) CreateSpec(spec string) (*Spec, error) {
	s := &Spec{}
	_, err := c.doRaw(http.MethodPost, "/v2/specs", []byte(spec), s, http.StatusOK)
	return s, err
}

// ReadSpec reads a job spec with the provided ID on the Chainlink node
func (c *chainlink) ReadSpec(id string) (*Response, error) {
	specObj := &Response{}
	_, err := c.do(http.MethodGet, fmt.Sprintf("/v2/specs/%s", id), nil, specObj, http.StatusOK)
	return specObj, err
}

// DeleteSpec deletes a job spec with the provided ID from the Chainlink node
func (c *chainlink) DeleteSpec(id string) error {
	_, err := c.do(http.MethodDelete, fmt.Sprintf("/v2/specs/%s", id), nil, nil, http.StatusNoContent)
	return err
}

// CreateBridge creates a bridge on the Chainlink node based on the provided attributes
func (c *chainlink) CreateBridge(bta *BridgeTypeAttributes) error {
	_, err := c.do(http.MethodPost, "/v2/bridge_types", bta, nil, http.StatusOK)
	return err
}

// ReadBridge reads a bridge from the Chainlink node based on the provided name
func (c *chainlink) ReadBridge(name string) (*BridgeType, error) {
	bt := BridgeType{}
	_, err := c.do(http.MethodGet, fmt.Sprintf("/v2/bridge_types/%s", name), nil, &bt, http.StatusOK)
	return &bt, err
}

// DeleteBridge deletes a bridge on the Chainlink node based on the provided name
func (c *chainlink) DeleteBridge(name string) error {
	_, err := c.do(http.MethodDelete, fmt.Sprintf("/v2/bridge_types/%s", name), nil, nil, http.StatusOK)
	return err
}

// CreateOCRKey creates an OCRKey on the Chainlink node
func (c *chainlink) CreateOCRKey() (*OCRKey, error) {
	ocrKey := &OCRKey{}
	_, err := c.do(http.MethodPost, "/v2/keys/ocr", nil, ocrKey, http.StatusOK)
	return ocrKey, err
}

// ReadOCRKeys reads all OCRKeys from the Chainlink node
func (c *chainlink) ReadOCRKeys() (*OCRKeys, error) {
	ocrKeys := &OCRKeys{}
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
	_, err := c.do(http.MethodDelete, fmt.Sprintf("/v2/keys/ocr/%s", id), nil, nil, http.StatusOK)
	return err
}

// CreateP2PKey creates an P2PKey on the Chainlink node
func (c *chainlink) CreateP2PKey() (*P2PKey, error) {
	p2pKey := &P2PKey{}
	_, err := c.do(http.MethodPost, "/v2/keys/p2p", nil, p2pKey, http.StatusOK)
	return p2pKey, err
}

// ReadP2PKeys reads all P2PKeys from the Chainlink node
func (c *chainlink) ReadP2PKeys() (*P2PKeys, error) {
	p2pKeys := &P2PKeys{}
	_, err := c.do(http.MethodGet, "/v2/keys/p2p", nil, p2pKeys, http.StatusOK)
	for index := range p2pKeys.Data {
		p2pKeys.Data[index].Attributes.PeerID = strings.TrimPrefix(p2pKeys.Data[index].Attributes.PeerID, "p2p_")
	}
	return p2pKeys, err
}

// DeleteP2PKey deletes a P2PKey on the Chainlink node based on the provided ID
func (c *chainlink) DeleteP2PKey(id int) error {
	_, err := c.do(http.MethodDelete, fmt.Sprintf("/v2/keys/p2p/%d", id), nil, nil, http.StatusOK)
	return err
}

// ReadETHKeys reads all ETH keys from the Chainlink node
func (c *chainlink) ReadETHKeys() (*ETHKeys, error) {
	ethKeys := &ETHKeys{}
	_, err := c.do(http.MethodGet, "/v2/keys/eth", nil, ethKeys, http.StatusOK)
	return ethKeys, err
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
	c.Cookies = resp.Cookies()

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
			"unexpected response code, got %d, expected 200\nURL: %s\nresponse received: %s",
			resp.StatusCode,
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
	log.Info().
		Str("Method", method).
		Str("Endpoint", endpoint).
		Str("URL", c.Config.URL).
		Str("Body", string(b)).
		Msg("Calling to Chainlink node")
	return c.doRaw(method, endpoint, b, obj, expectedStatusCode)
}

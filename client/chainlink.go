package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

var ErrNotFound = errors.New("unexpected response code, got 404")
var ErrUnprocessableEntity = errors.New("unexpected response code, got 422")

// Chainlink interface that enables interactions with a chainlink node
type Chainlink interface {
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
	Client  *http.Client
	Config  *ChainlinkConfig
	Cookies []*http.Cookie
}

// NewChainlink creates a new chainlink model using a provided config
func NewChainlink(c *ChainlinkConfig) Chainlink {
	cl := &chainlink{Config: c}
	return cl
}

// CreateTemplateNodes lauches 5 chainlink nodes in a default config for testing
func CreateTemplateNodes(network BlockchainNetwork, linkAddress string) ([]Chainlink, error) {
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
	cmd.Dir = string(projectDir) + "/chainlink_nodes"
	ethUrl := "ETH_URL=" + network.URL()
	if network.ID() == string(EthereumHardhatID) {
		ethUrl = "ETH_URL=ws://host.docker.internal:8545"
	}
	chainId := "ETH_CHAIN_ID=" + network.ChainID().String()
	la := "LINK_CONTRACT_ADDRESS=" + linkAddress
	cmd.Env = []string{ethUrl, chainId, la}
	var e bytes.Buffer
	cmd.Stderr = &e
	cmd.Start()

	// Wait for Docker Compose to be up and healthy
	resp, err := http.Get(urlBase + "6711")
	for start := time.Now(); time.Since(start) < 2*time.Minute; time.Sleep(time.Second * 20) {
		resp, err = http.Get(urlBase + "6711")
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
		CleanTemplateNodes()
		return nil, err
	}

	var cls []Chainlink
	for _, port := range ports {
		c := &ChainlinkConfig{
			URL:      urlBase + port,
			Email:    email,
			Password: pass,
		}
		cls = append(cls, NewChainlink(c))
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
	cmd.Dir = string(projectDir) + "/chainlink_nodes"

	log.Info().Str("CMD", "docker-compose -f ./docker-compose.yml down -v --remove-orphans").Msg("Running command")
	return cmd.Run()
}

// CreateJob creates a Chainlink job based on the provided spec string
func (c *chainlink) CreateJob(spec string) (*Job, error) {
	job := &Job{}
	_, err := c.do(http.MethodPost, "/v2/jobs", &JobForm{
		TOML: spec,
	}, &job, http.StatusOK)
	return job, err
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
	c.Client = client
}

func (c *chainlink) doRaw(
	method,
	endpoint string,
	body []byte, obj interface{},
	expectedStatusCode int,
) (*http.Response, error) {
	client := c.Client

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
	return c.doRaw(method, endpoint, b, obj, expectedStatusCode)
}

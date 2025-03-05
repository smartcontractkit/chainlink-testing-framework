// atlas_client.go
package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/logging"
)

// AtlasClient defines the structure of the Atlas service client.
type AtlasClient struct {
	BaseURL    string
	RestClient *resty.Client
	Logger     logging.Logger
}

// TransactionHash defines a single transaction hash.
type TransactionHash struct {
	MessageID string `json:"messageId"`
}

// TransactionResponse defines the response format from Atlas for transaction details.
type TransactionResponse struct {
	TransactionHash []TransactionHash `json:"transactionHash"`
}

// NewAtlasClient creates a new Atlas client instance.
func NewAtlasClient(baseURL string) *AtlasClient {
	logging.Init()
	logger := logging.GetLogger(nil, "ATLAS_CLIENT_LOG_LEVEL")
	logger.Info().
		Str("BaseURL", baseURL).
		Msg("Initializing Atlas Client")

	isDebug := os.Getenv("RESTY_DEBUG") == "true"
	restyClient := resty.New().SetDebug(isDebug)

	return &AtlasClient{
		BaseURL:    baseURL,
		RestClient: restyClient,
		Logger:     logger,
	}
}

// GetTransactionDetails retrieves transaction details using the provided msgIdOrTxnHash.
func (ac *AtlasClient) GetTransactionDetails(msgIdOrTxnHash string) (*TransactionResponse, error) {
	endpoint := fmt.Sprintf("%s/atlas/search?msgIdOrTxnHash=%s", ac.BaseURL, msgIdOrTxnHash)
	ac.Logger.Info().
		Str("msgIdOrTxnHash", msgIdOrTxnHash).
		Msg("Sending request to fetch transaction details")

	resp, err := ac.RestClient.R().
		SetHeader("Content-Type", "application/json").
		Get(endpoint)
	if err != nil {
		ac.Logger.Error().Err(err).Msg("Failed to send request")
		return nil, err
	}

	if resp.StatusCode() != 200 {
		ac.Logger.Error().Int("statusCode", resp.StatusCode()).Msg("Received non-OK status")
		return nil, errors.New("failed to retrieve transaction details")
	}

	var transactionResponse TransactionResponse
	ac.Logger.Debug().Bytes("Response", resp.Body())
	if err := json.Unmarshal(resp.Body(), &transactionResponse); err != nil {
		ac.Logger.Error().Err(err).Msg("Failed to unmarshal response")
		return nil, err
	}

	ac.Logger.Info().
		Str("transactionHash", msgIdOrTxnHash).
		Msg("Successfully retrieved transaction details")
	return &transactionResponse, nil
}

// CustomTime wraps time.Time to support custom unmarshalling.
type CustomTime struct {
	time.Time
}

// UnmarshalJSON parses a timestamp string into a CustomTime.
func (ct *CustomTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	if s == "null" || s == "" {
		ct.Time = time.Time{}
		return nil
	}
	t, err := time.Parse("2006-01-02T15:04:05", s)
	if err != nil {
		return err
	}
	ct.Time = t
	return nil
}

// TransactionDetails represents detailed transaction data returned by Atlas.
type TransactionDetails struct {
	MessageId              *string     `json:"messageId"`
	State                  *int        `json:"state"`
	Votes                  *int        `json:"votes"`
	SourceNetworkName      *string     `json:"sourceNetworkName"`
	DestNetworkName        *string     `json:"destNetworkName"`
	CommitBlockTimestamp   *CustomTime `json:"commitBlockTimestamp"`
	Root                   *string     `json:"root"`
	SendFinalized          *CustomTime `json:"sendFinalized"`
	CommitStore            *string     `json:"commitStore"`
	Origin                 *string     `json:"origin"`
	SequenceNumber         *int        `json:"sequenceNumber"`
	Sender                 *string     `json:"sender"`
	Receiver               *string     `json:"receiver"`
	SourceChainId          *string     `json:"sourceChainId"`
	DestChainId            *string     `json:"destChainId"`
	RouterAddress          *string     `json:"routerAddress"`
	OnrampAddress          *string     `json:"onrampAddress"`
	OfframpAddress         *string     `json:"offrampAddress"`
	DestRouterAddress      *string     `json:"destRouterAddress"`
	SendTransactionHash    *string     `json:"sendTransactionHash"`
	SendTimestamp          *CustomTime `json:"sendTimestamp"`
	SendBlock              *int        `json:"sendBlock"`
	SendLogIndex           *int        `json:"sendLogIndex"`
	Min                    *string     `json:"min"`
	Max                    *string     `json:"max"`
	CommitTransactionHash  *string     `json:"commitTransactionHash"`
	CommitBlockNumber      *int        `json:"commitBlockNumber"`
	CommitLogIndex         *int        `json:"commitLogIndex"`
	Arm                    *string     `json:"arm"`
	BlessTransactionHash   *string     `json:"blessTransactionHash"`
	BlessBlockNumber       *int        `json:"blessBlockNumber"`
	BlessBlockTimestamp    *CustomTime `json:"blessBlockTimestamp"`
	BlessLogIndex          *int        `json:"blessLogIndex"`
	ReceiptTransactionHash *string     `json:"receiptTransactionHash"`
	ReceiptTimestamp       *CustomTime `json:"receiptTimestamp"`
	ReceiptBlock           *int        `json:"receiptBlock"`
	ReceiptLogIndex        *int        `json:"receiptLogIndex"`
	ReceiptFinalized       *CustomTime `json:"receiptFinalized"`
	Data                   *string     `json:"data"`
	Strict                 *bool       `json:"strict"`
	Nonce                  *int        `json:"nonce"`
	FeeToken               *string     `json:"feeToken"`
	GasLimit               *string     `json:"gasLimit"`
	FeeTokenAmount         *string     `json:"feeTokenAmount"`
	TokenAmounts           *[]string   `json:"tokenAmounts"`
}

// GetMessageDetails fetches detailed transaction info using the message endpoint.
func (ac *AtlasClient) GetMessageDetails(messageID string) (*TransactionDetails, error) {
	endpoint := fmt.Sprintf("%s/atlas/message/%s", ac.BaseURL, messageID)
	ac.Logger.Info().
		Str("messageID", messageID).
		Msg("Sending request to fetch message details")

	resp, err := ac.RestClient.R().
		SetHeader("Content-Type", "application/json").
		Get(endpoint)
	if err != nil {
		ac.Logger.Error().Err(err).Msg("Failed to send request for message details")
		return nil, err
	}

	if resp.StatusCode() != 200 {
		ac.Logger.Error().Int("statusCode", resp.StatusCode()).Msg("Received non-OK status for message details")
		return nil, fmt.Errorf("failed to retrieve message details, status code: %d", resp.StatusCode())
	}

	body := resp.Body()
	if string(body) == "Message not found" {
		ac.Logger.Warn().
			Str("messageID", messageID).
			Msg("Message not found in Atlas")
		return nil, errors.New("message not found")
	}

	var details TransactionDetails
	ac.Logger.Debug().Bytes("Response", body)
	if err := json.Unmarshal(body, &details); err != nil {
		ac.Logger.Error().Err(err).Msg("Failed to unmarshal message details")
		return nil, err
	}

	return &details, nil
}

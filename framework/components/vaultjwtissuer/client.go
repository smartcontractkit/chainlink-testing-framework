package vaultjwtissuer

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const DefaultAudience = "https://vault.test.chain.link"

// ClaimVaultSecretManagementEnabled matches the Vault authorizer claim gate in
// Chainlink. Test JWTs should set this so downstream Vault auth behaves the
// same way as production JWT-authenticated requests.
const ClaimVaultSecretManagementEnabled = "urn:chainlink:claim_vault_secret_management_enabled" // #nosec G101 -- static JWT claim name, not credentials

type TokenClaims struct {
	OrgID         string
	WorkflowOwner string
	RequestDigest string
	Issuer        string
	Audience      string
	Subject       string
	JWTID         string
	KeyID         string
	IssuedAt      time.Time
	ExpiresAt     time.Time
	ExtraClaims   map[string]any
}

type Client struct {
	localURL     string
	dockerURL    string
	privateKey   *rsa.PrivateKey
	defaultKeyID string
}

func NewClient(localURL, dockerURL string) (*Client, error) {
	privateKey, err := parseDefaultJWTSigningKey()
	if err != nil {
		return nil, err
	}

	return &Client{
		localURL:     NormalizeIssuerURL(localURL),
		dockerURL:    NormalizeIssuerURL(dockerURL),
		privateKey:   privateKey,
		defaultKeyID: DefaultJWTIssuerKeyID,
	}, nil
}

func NewClientFromOutput(out *Output) (*Client, error) {
	if out == nil {
		return nil, errors.New("vault JWT issuer output is nil")
	}

	return NewClient(out.LocalHTTPURL, out.DockerHTTPURL)
}

func (c *Client) LocalIssuerURL() string {
	if c == nil {
		return ""
	}
	return c.localURL
}

func (c *Client) DockerIssuerURL() string {
	if c == nil {
		return ""
	}
	return c.dockerURL
}

func (c *Client) MintToken(claims TokenClaims) (string, error) {
	if c == nil || c.privateKey == nil {
		return "", errors.New("JWT issuer signing key is not configured")
	}
	if claims.KeyID == "" {
		claims.KeyID = c.defaultKeyID
	}
	if claims.Issuer == "" {
		claims.Issuer = c.LocalIssuerURL()
	}
	if claims.Audience == "" {
		claims.Audience = DefaultAudience
	}

	return signToken(c.privateKey, claims)
}

func NormalizeIssuerURL(raw string) string {
	if raw == "" || strings.HasSuffix(raw, "/") {
		return raw
	}
	return raw + "/"
}

func signToken(privateKey *rsa.PrivateKey, claims TokenClaims) (string, error) {
	if privateKey == nil {
		return "", errors.New("private key is required")
	}
	if claims.KeyID == "" {
		return "", errors.New("kid is required")
	}
	if claims.Issuer == "" {
		return "", errors.New("issuer is required")
	}
	if claims.OrgID == "" {
		return "", errors.New("org_id is required")
	}
	if claims.RequestDigest == "" {
		return "", errors.New("request_digest is required")
	}

	now := time.Now().UTC()
	if claims.IssuedAt.IsZero() {
		claims.IssuedAt = now
	}
	if claims.ExpiresAt.IsZero() {
		claims.ExpiresAt = claims.IssuedAt.Add(5 * time.Minute)
	}
	if claims.Subject == "" {
		claims.Subject = claims.OrgID
	}
	if claims.Audience == "" {
		claims.Audience = DefaultAudience
	}

	authorizationDetails := []map[string]string{{
		"type":  "request_digest",
		"value": claims.RequestDigest,
	}}
	if claims.WorkflowOwner != "" {
		authorizationDetails = append(authorizationDetails, map[string]string{
			"type":  "workflow_owner",
			"value": claims.WorkflowOwner,
		})
	}

	tokenClaims := jwt.MapClaims{
		"iss":                                claims.Issuer,
		"aud":                                claims.Audience,
		"sub":                                claims.Subject,
		"iat":                                jwt.NewNumericDate(claims.IssuedAt),
		"exp":                                jwt.NewNumericDate(claims.ExpiresAt),
		"org_id":                             claims.OrgID,
		ClaimVaultSecretManagementEnabled:    "true",
		"authorization_details":              authorizationDetails,
	}
	if claims.JWTID != "" {
		tokenClaims["jti"] = claims.JWTID
	}
	for key, value := range claims.ExtraClaims {
		tokenClaims[key] = value
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, tokenClaims)
	token.Header["kid"] = claims.KeyID

	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign JWT: %w", err)
	}

	return tokenString, nil
}

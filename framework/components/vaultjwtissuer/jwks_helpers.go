package vaultjwtissuer

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"strings"
)

const (
	DefaultJWTIssuerKeyID   = "vault-jwt-test-key"
	defaultJWTPrivateKeyPEM = `-----BEGIN PRIVATE KEY-----
MIIEvwIBADANBgkqhkiG9w0BAQEFAASCBKkwggSlAgEAAoIBAQDYhEVPZ8YdC3Va
DGZ2hWPt+VYptOt0heTulBOwBW0ESavpfvokLYGFu+bLkGhIw365nCFw0eulLZYN
tD4nzq7F5Swtb2iIaDK19PBVNcukU/CY6j44KC1eomyaOvPXKWKwcc7qxjy9bIyA
TyOmOlxNxcNRSjL2SOApFkzb8M/RymHlMT/RY5ubytvjcbQgn2gy19U7HuNLYW1P
gviAAMY635u0A+HAxXx83lQSz9gy08/uBarmKAd2OadCA8cNiTSYyfUS6m1pycA7
j8ZHY75xL4hm+p2PJd9V1x3Z4S1TpZDIj+YAG/v4ZHB1vLTLoPIgwLEqwGRRWijl
sbdUZRd9AgMBAAECggEAGCiWFTiWheof43bLvgC/OC/gedHajctc0nQKSFMqqVZR
DMIixgOf1pyzMVaBFFFf4/T0VELQAMO34PqSDt4EaUdbaQxrxQCfW+cjI9bXTJQj
HeTRIXH2Mf98j67xQzo2bUqdlFufLmGcwbpS13rejrz4wKq/SfSyslLvK4FQpu8x
5J9ntn2wdgeUQCm62FyuNPxFMBldcovnwf9bbojTjMAatWfyF++W8OAcRqZCab1H
1WNPyhBqG5vDVMtgBdTkwZHqI01B+ozMnBLuEhsLVzvQWE79ZouWtU76GIeFlr0n
bC/3uWq9LBo1kEbLIPucxYA14ytWfpQwUvy1k11s4QKBgQD4dz2fVYSVb6hn0Pon
EQtunruNB7F2JlobY2s3C7aBKs+l48J16whKFcqHUA6NpuSvyUhFTqIpxM0LXdar
6nWu4Yw0kbqACJOHXuG71VhfkUgRJMOZoC/V0RKudoTwWDzFgNXvYF3bqtpmQDW7
2dUrSJ+jMOU7eCzXOdHDTFGhbQKBgQDfFQT/NACHapIn5w6c1Dha6fy7t1Z6A2zw
bUUzAh5C1kZ8yeDrkVfr5Ys+Y7Am/tfFteXO2XRSGH5yqq9YHVr0RihavqX72FGT
YY2rmyht+JjnZ3y+vOG5LXePR9tilvGei3jH0lTRPdwKpa6feHKry9MBx5xmqKqQ
xKRmyXaUUQKBgQCcOp3MqgEL1YGWhZhFKDp/+98B9mxnVgYiYojvu7Wt0jVuoZ+M
dZRowPrvyi7ccqwou+9tZNwiV1R2aTKqNmp44+k8xMT37GyXGdnmOWev77HY1b0H
w+lQEH4mpO9CELlllnTuZzGdBfj9gjJHQ9j9tlRqUDxTAGVxjzGOE1bgoQKBgQCu
DxmCAlIzVqzJY5hcN53tGcrvsKJRu2CBy9CFdy6jWctPzLipNROT5Nubh27HTmqP
QlkX50XCVIg88f60UttH44HTJBQgh+1GgIRolDycaa7sRyvnKzs4IEi8TAXaTAok
eZB44Rz60jhhOlsg5HscnoF6TwQyeYH0SOo5pRHXsQKBgQCY/pua7PceD5ZQ4lae
Pi5E9LzPjoeFegVgAP7bRUeC21nzLZlKYOcRCV2WkGLsz60bZm+7VEyFZmrrFoTE
58G0eCLCUq3Dj+NPfIvXNWwSuUAdDspWOBSCyENP+y+jLzIa2OtCj+KJe6Oe28pf
CcSeCJqr6aLeDRPcuD7yUat1OA==
-----END PRIVATE KEY-----`
)

type jwtWebKey struct {
	Kid string `json:"kid"`
	Alg string `json:"alg"`
	Kty string `json:"kty"`
	Use string `json:"use"`
	N   string `json:"n"`
	E   string `json:"e"`
}

func parseDefaultJWTSigningKey() (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(defaultJWTPrivateKeyPEM))
	if block == nil {
		return nil, errors.New("failed to decode default JWT signing key PEM")
	}
	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse default JWT signing key: %w", err)
	}
	rsaKey, ok := privateKey.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("default JWT signing key is not RSA")
	}
	return rsaKey, nil
}

func rsaPublicKeyToJWK(keyID string, publicKey *rsa.PublicKey) jwtWebKey {
	return jwtWebKey{
		Kid: keyID,
		Alg: "RS256",
		Kty: "RSA",
		Use: "sig",
		N:   base64.RawURLEncoding.EncodeToString(publicKey.N.Bytes()),
		E:   base64.RawURLEncoding.EncodeToString(big.NewInt(int64(publicKey.E)).Bytes()),
	}
}

func requestBaseURL(r *http.Request) string {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	if forwardedProto := strings.TrimSpace(r.Header.Get("X-Forwarded-Proto")); forwardedProto != "" {
		scheme = forwardedProto
	}
	host := strings.TrimSpace(r.Host)
	if forwardedHost := strings.TrimSpace(r.Header.Get("X-Forwarded-Host")); forwardedHost != "" {
		host = forwardedHost
	}
	return fmt.Sprintf("%s://%s", scheme, host)
}

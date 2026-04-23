package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// PKIClient interacts with Vault's PKI secrets engine to issue certificates.
type PKIClient struct {
	address string
	token   string
	mount   string
	client  *http.Client
}

// Certificate holds the issued certificate data returned by Vault.
type Certificate struct {
	SerialNumber string
	Certificate  string
	PrivateKey   string
	CA           string
	Expiration   time.Time
}

// IssueCertRequest defines the parameters for issuing a certificate.
type IssueCertRequest struct {
	Role       string
	CommonName string
	TTL        string
	AltNames   []string
}

// NewPKIClient creates a new PKIClient. mount defaults to "pki".
func NewPKIClient(address, token, mount string) (*PKIClient, error) {
	if address == "" {
		return nil, fmt.Errorf("pki: vault address is required")
	}
	if token == "" {
		return nil, fmt.Errorf("pki: vault token is required")
	}
	if mount == "" {
		mount = "pki"
	}
	return &PKIClient{
		address: strings.TrimRight(address, "/"),
		token:   token,
		mount:   mount,
		client:  &http.Client{Timeout: 15 * time.Second},
	}, nil
}

// IssueCertificate requests Vault to issue a certificate for the given role.
func (p *PKIClient) IssueCertificate(req IssueCertRequest) (*Certificate, error) {
	if req.Role == "" {
		return nil, fmt.Errorf("pki: role is required")
	}
	if req.CommonName == "" {
		return nil, fmt.Errorf("pki: common_name is required")
	}

	body := map[string]interface{}{
		"common_name": req.CommonName,
	}
	if req.TTL != "" {
		body["ttl"] = req.TTL
	}
	if len(req.AltNames) > 0 {
		body["alt_names"] = strings.Join(req.AltNames, ",")
	}

	url := fmt.Sprintf("%s/v1/%s/issue/%s", p.address, p.mount, req.Role)
	payload, err := jsonBody(body)
	if err != nil {
		return nil, fmt.Errorf("pki: encoding request: %w", err)
	}

	httpReq, err := http.NewRequest(http.MethodPost, url, payload)
	if err != nil {
		return nil, fmt.Errorf("pki: creating request: %w", err)
	}
	httpReq.Header.Set("X-Vault-Token", p.token)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("pki: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("pki: unexpected status %d for role %s", resp.StatusCode, req.Role)
	}

	var result struct {
		Data struct {
			SerialNumber string `json:"serial_number"`
			Certificate  string `json:"certificate"`
			PrivateKey   string `json:"private_key"`
			IssuingCA    string `json:"issuing_ca"`
			Expiration   int64  `json:"expiration"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("pki: decoding response: %w", err)
	}

	return &Certificate{
		SerialNumber: result.Data.SerialNumber,
		Certificate:  result.Data.Certificate,
		PrivateKey:   result.Data.PrivateKey,
		CA:           result.Data.IssuingCA,
		Expiration:   time.Unix(result.Data.Expiration, 0),
	}, nil
}

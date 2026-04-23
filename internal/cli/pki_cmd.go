package cli

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/example/vaultpull/internal/vault"
)

// PKIFlags holds parsed CLI flags for the pki subcommand.
type PKIFlags struct {
	Address    string
	Token      string
	Mount      string
	Role       string
	CommonName string
	TTL        string
	AltNames   string
}

// ParsePKIFlags parses CLI arguments for the pki issue subcommand.
func ParsePKIFlags(args []string) (*PKIFlags, error) {
	fs := flag.NewFlagSet("pki", flag.ContinueOnError)

	address := fs.String("address", os.Getenv("VAULT_ADDR"), "Vault address")
	token := fs.String("token", os.Getenv("VAULT_TOKEN"), "Vault token")
	mount := fs.String("mount", "pki", "PKI secrets engine mount path")
	role := fs.String("role", "", "PKI role name (required)")
	cn := fs.String("common-name", "", "Certificate common name (required)")
	ttl := fs.String("ttl", "", "Certificate TTL (e.g. 24h)")
	altNames := fs.String("alt-names", "", "Comma-separated SAN alt names")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}
	if *role == "" {
		return nil, fmt.Errorf("pki: --role is required")
	}
	if *cn == "" {
		return nil, fmt.Errorf("pki: --common-name is required")
	}

	return &PKIFlags{
		Address:    *address,
		Token:      *token,
		Mount:      *mount,
		Role:       *role,
		CommonName: *cn,
		TTL:        *ttl,
		AltNames:   *altNames,
	}, nil
}

// RunPKI issues a certificate using the Vault PKI engine and prints the result.
func RunPKI(f *PKIFlags) error {
	c, err := vault.NewPKIClient(f.Address, f.Token, f.Mount)
	if err != nil {
		return fmt.Errorf("pki: %w", err)
	}

	var altNames []string
	if f.AltNames != "" {
		altNames = strings.Split(f.AltNames, ",")
	}

	cert, err := c.IssueCertificate(vault.IssueCertRequest{
		Role:       f.Role,
		CommonName: f.CommonName,
		TTL:        f.TTL,
		AltNames:   altNames,
	})
	if err != nil {
		return fmt.Errorf("pki: issuing certificate: %w", err)
	}

	fmt.Printf("Serial:      %s\n", cert.SerialNumber)
	fmt.Printf("Expiration:  %s\n", cert.Expiration.Format("2006-01-02T15:04:05Z"))
	fmt.Printf("Certificate:\n%s\n", cert.Certificate)
	fmt.Printf("Private Key:\n%s\n", cert.PrivateKey)
	fmt.Printf("CA:\n%s\n", cert.CA)
	return nil
}

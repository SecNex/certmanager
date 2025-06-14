package certificate

import (
	"fmt"

	"github.com/go-acme/lego/certificate"
	"github.com/go-acme/lego/challenge/dns01"
	"github.com/go-acme/lego/challenge/http01"
	"github.com/google/uuid"
	"github.com/secnex/certmanager/common/account"
	"github.com/secnex/certmanager/models"
	"github.com/secnex/certmanager/store"
)

type ChallengeType string

const (
	ChallengeTypeHTTP ChallengeType = "http"
	ChallengeTypeDNS  ChallengeType = "dns"
)

type CertificateConfig struct {
	ChallengeType ChallengeType
	DNSProvider   string // Provider name for DNS challenge
}

type Certificate struct {
	models.Certificate
	Account *account.Account
	Cert    *certificate.Resource
	Config  *CertificateConfig
}

func NewCertificate(domains []string, account *account.Account, store *store.Storage) (*Certificate, error) {
	return NewCertificateWithConfig(domains, account, store, &CertificateConfig{
		ChallengeType: ChallengeTypeHTTP,
	})
}

func NewCertificateWithConfig(domains []string, account *account.Account, store *store.Storage, config *CertificateConfig) (*Certificate, error) {
	cert := &Certificate{
		Certificate: models.Certificate{
			ID:        uuid.New(),
			Domains:   domains,
			AccountID: account.ID,
		},
		Account: account,
		Config:  config,
	}

	err := cert.RequestNewCertificate()
	if err != nil {
		return nil, err
	}

	err = cert.Save(store)
	if err != nil {
		return nil, err
	}

	return cert, nil
}

func (c *Certificate) RequestNewCertificate() error {
	switch c.Config.ChallengeType {
	case ChallengeTypeHTTP:
		err := c.setupHTTPChallenge()
		if err != nil {
			return fmt.Errorf("failed to setup HTTP challenge: %w", err)
		}
	case ChallengeTypeDNS:
		err := c.setupDNSChallenge()
		if err != nil {
			return fmt.Errorf("failed to setup DNS challenge: %w", err)
		}
	default:
		return fmt.Errorf("unsupported challenge type: %s", c.Config.ChallengeType)
	}

	request := certificate.ObtainRequest{
		Domains: c.Domains,
		Bundle:  true,
	}

	certs, err := c.Account.Client.Certificate.Obtain(request)
	if err != nil {
		return err
	}

	c.Cert = certs
	return nil
}

func (c *Certificate) setupHTTPChallenge() error {
	provider := http01.NewProviderServer("", "80")
	return c.Account.Client.Challenge.SetHTTP01Provider(provider)
}

func (c *Certificate) setupDNSChallenge() error {
	// Use manual DNS provider - requires user to manually set DNS records
	provider, err := dns01.NewDNSProviderManual()
	if err != nil {
		return fmt.Errorf("failed to create manual DNS provider: %w", err)
	}
	return c.Account.Client.Challenge.SetDNS01Provider(provider)
}

func (c *Certificate) Save(store *store.Storage) error {
	return store.SaveCertificate(c.ID.String(), c.Cert)
}

// GetSupportedChallengeTypes returns a list of supported challenge types
func GetSupportedChallengeTypes() []string {
	return []string{"http", "dns"}
}

// ValidateDNSProviderConfig validates DNS provider configuration
func ValidateDNSProviderConfig(provider string) error {
	// For manual DNS, no specific environment variables are required
	// The user will need to manually set DNS records as prompted
	if provider == "manual" {
		return nil
	}

	// For future automated DNS providers, add validation here
	return fmt.Errorf("DNS provider '%s' not supported yet. Use 'manual' for manual DNS record setup", provider)
}

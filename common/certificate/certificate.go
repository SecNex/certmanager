package certificate

import (
	"log"

	"github.com/go-acme/lego/certificate"
	"github.com/go-acme/lego/challenge/http01"
	"github.com/google/uuid"
	"github.com/secnex/certmanager/common/account"
	"github.com/secnex/certmanager/models"
	"github.com/secnex/certmanager/store"
)

type Certificate struct {
	models.Certificate
	Account *account.Account
	Cert    *certificate.Resource
}

func NewCertificate(domains []string, account *account.Account, store *store.Storage) (*Certificate, error) {
	cert := &Certificate{
		Certificate: models.Certificate{
			ID:        uuid.New(),
			Domains:   domains,
			AccountID: account.ID,
		},
		Account: account,
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
	err := c.Account.Client.Challenge.SetHTTP01Provider(http01.NewProviderServer("", "80"))
	if err != nil {
		log.Fatal(err)
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

func (c *Certificate) Save(store *store.Storage) error {
	return store.SaveCertificate(c.ID.String(), c.Cert)
}

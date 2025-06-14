package manager

import (
	"github.com/secnex/certmanager/common/account"
	"github.com/secnex/certmanager/common/certificate"
	"github.com/secnex/certmanager/database"
	"github.com/secnex/certmanager/store"
)

type Manager struct {
	Database database.Database
	Storage  *store.Storage
}

func NewManager(database *database.Database, storage *store.Storage) *Manager {
	return &Manager{
		Database: *database,
		Storage:  storage,
	}
}

func (m *Manager) NewAccount(email string) (*account.Account, error) {
	return account.NewAccount(email, &m.Database, m.Storage)
}

func (m *Manager) NewCertificate(domains []string, account *account.Account) (*certificate.Certificate, error) {
	return certificate.NewCertificate(domains, account, m.Storage)
}

package certmanager

import (
	"github.com/secnex/certmanager/database"
	"github.com/secnex/certmanager/manager"
	"github.com/secnex/certmanager/server"
	"github.com/secnex/certmanager/store"
)

type CertManager struct {
	Manager *manager.Manager
}

func NewCertManager(database *database.Database, storage *store.Storage) *CertManager {
	return &CertManager{
		Manager: manager.NewManager(database, storage),
	}
}

func (m *CertManager) RunServer() {
	apiServer := server.NewApiServer(nil, nil, m.Manager)
	apiServer.Start()
}

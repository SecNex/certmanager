package certmanager

import "github.com/secnex/certmanager/database"

type Manager struct {
	Database database.Database
}

func NewManager(database *database.Database) *Manager {
	return &Manager{
		Database: *database,
	}
}

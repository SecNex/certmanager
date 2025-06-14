package models

import "github.com/google/uuid"

type CertificateType string

const (
	CertificateTypeLetsEncrypt CertificateType = "letsencrypt"
	CertificateTypeCustom      CertificateType = "custom"
)

type Certificate struct {
	ID        uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	AccountID uuid.UUID `gorm:"not null" json:"account_id"`

	Domains []string `gorm:"not null" json:"domains"`

	Account Account `gorm:"foreignKey:AccountID" json:"account"`
}

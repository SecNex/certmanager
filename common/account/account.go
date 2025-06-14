package account

import (
	"crypto"
	"crypto/rsa"
	"log"

	"github.com/go-acme/lego/certcrypto"
	"github.com/go-acme/lego/lego"
	"github.com/go-acme/lego/registration"
	"github.com/google/uuid"
	"github.com/secnex/certmanager/database"
	"github.com/secnex/certmanager/models"
	"github.com/secnex/certmanager/store"
)

type Account struct {
	models.Account
	PrivateKey         *crypto.PrivateKey
	LetsEncryptAccount *LetsEncryptAccount
	Client             *lego.Client
}

type LetsEncryptAccount struct {
	Email        string
	Registration *registration.Resource
	Key          *rsa.PrivateKey
}

func (u *LetsEncryptAccount) GetEmail() string {
	return u.Email
}
func (u *LetsEncryptAccount) GetRegistration() *registration.Resource {
	return u.Registration
}
func (u *LetsEncryptAccount) GetPrivateKey() crypto.PrivateKey {
	return u.Key
}

func NewAccount(email string, cnx *database.Database, storage *store.Storage) (*Account, error) {
	account := Account{
		Account: models.Account{
			Email: email,
		},
	}
	// Check if account for this email already exists
	var existingAccount models.Account
	log.Println("Checking if account for this email already exists")
	err := cnx.Database.Where("email = ?", account.Email).First(&existingAccount).Error
	// If the account exists, is not a error, but an empty struct
	if err != nil && err.Error() == "record not found" {
		err = nil
	}

	if existingAccount.ID != uuid.Nil {
		account.ID = existingAccount.ID

		privateKey, err := storage.ReadPrivateKey(account.ID.String())
		if err != nil {
			return nil, err
		}

		account.PrivateKey = &privateKey

		account.LetsEncryptAccount = &LetsEncryptAccount{
			Email: account.Email,
			Key:   privateKey.(*rsa.PrivateKey),
		}

		account.Client, err = account.CreateClient()
		if err != nil {
			return nil, err
		}

		return &account, nil
	}

	log.Println("Account does not exist, creating new account!")
	var newAccount Account
	newAccount.Email = account.Email
	err = cnx.Database.Create(&newAccount).Error
	if err != nil {
		return nil, err
	}
	account.ID = newAccount.ID

	log.Println("Creating private key for account!")
	privateKey, err := account.CreatePrivateKey(storage)
	if err != nil {
		return nil, err
	}

	account.LetsEncryptAccount = &LetsEncryptAccount{
		Email: account.Email,
		Key:   (*privateKey).(*rsa.PrivateKey),
	}

	client, err := account.CreateClient()
	if err != nil {
		return nil, err
	}
	account.Client = client

	return &account, nil
}

func (a *Account) CreateClient() (*lego.Client, error) {
	config := lego.NewConfig(a.LetsEncryptAccount)
	client, err := lego.NewClient(config)
	if err != nil {
		return nil, err
	}
	return client, nil
}
func (a *Account) CreatePrivateKey(storage *store.Storage) (*crypto.PrivateKey, error) {
	accountPrivateKey, err := certcrypto.GeneratePrivateKey(certcrypto.RSA2048)
	if err != nil {
		return nil, err
	}

	err = storage.SavePrivateKey(a.ID.String(), accountPrivateKey)
	if err != nil {
		return nil, err
	}

	a.PrivateKey = &accountPrivateKey

	return &accountPrivateKey, nil
}

func (a *Account) Create(cnx *database.Database, storage *store.Storage) (*Account, error) {
	// Check if account for this email already exists
	var existingAccount models.Account
	log.Println("Checking if account for this email already exists")
	err := cnx.Database.Where("email = ?", a.Email).First(&existingAccount).Error
	if err != nil {
		return nil, err
	}

	if existingAccount.ID != uuid.Nil {
		a.ID = existingAccount.ID
		return a, nil
	}

	return a, nil
}

func (a *Account) GetAccount(cnx *database.Database, id string) (*Account, error) {
	var account Account
	err := cnx.Database.Where("id = ?", id).First(&account).Error
	if err != nil {
		return nil, err
	}
	return &account, nil
}

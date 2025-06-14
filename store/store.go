package store

import (
	"bytes"
	"context"
	"crypto"
	"crypto/x509"
	"encoding/pem"
	"io"

	"github.com/go-acme/lego/certificate"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Storage struct {
	Context context.Context
	Client  *minio.Client
	Bucket  string
}

func NewStorage(endpoint string, accessKey string, secretKey string, bucket string) (*Storage, error) {
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: true,
	})
	if err != nil {
		return nil, err
	}

	return &Storage{Context: context.Background(), Client: minioClient, Bucket: bucket}, nil
}

func (s *Storage) SavePrivateKey(id string, key crypto.PrivateKey) error {
	keyBytes, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		return err
	}
	pemBlock := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: keyBytes,
	}
	pemBytes := pem.EncodeToMemory(pemBlock)
	return s.Save(id, pemBytes)
}

func (s *Storage) Save(key string, data []byte) error {
	_, err := s.Client.PutObject(s.Context, s.Bucket, key, bytes.NewReader(data), int64(len(data)), minio.PutObjectOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) SaveCertificate(id string, cert *certificate.Resource) error {
	return s.Save(id, cert.Certificate)
}

func (s *Storage) ReadPrivateKey(id string) (crypto.PrivateKey, error) {
	data, err := s.Read(id)
	if err != nil {
		return nil, err
	}

	// Decode PEM block first
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, err
	}

	return x509.ParsePKCS8PrivateKey(block.Bytes)
}

func (s *Storage) Read(key string) ([]byte, error) {
	object, err := s.Client.GetObject(s.Context, s.Bucket, key, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	defer object.Close()
	return io.ReadAll(object)
}

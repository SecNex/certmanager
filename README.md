# Certificate Manager

A comprehensive Go-based certificate management system that automates SSL/TLS certificate provisioning, management, and storage using Let's Encrypt ACME protocol.

## 🚀 Features

- **Automated Certificate Management**: Seamless SSL/TLS certificate provisioning through Let's Encrypt
- **Secure Storage**: Private keys and certificates stored securely in MinIO object storage
- **Database Integration**: PostgreSQL database for account and certificate metadata
- **HTTP/HTTPS API Server**: RESTful API with automatic HTTPS certificate provisioning
- **Account Management**: ACME account creation and management
- **Challenge Support**: HTTP-01 challenge provider for domain validation
- **Logging**: Comprehensive HTTP request logging
- **Health Monitoring**: Built-in health check endpoints

## 📋 Prerequisites

- **Go 1.24.3+**
- **PostgreSQL 12+**
- **MinIO or S3-compatible object storage**
- **Domain with DNS control** (for Let's Encrypt certificates)

## 🛠️ Installation

### 1. Clone the Repository

```bash
git clone https://github.com/secnex/certmanager.git
cd certmanager
```

### 2. Install Dependencies

```bash
go mod download
```

### 3. Database Setup

Create a PostgreSQL database:

```sql
CREATE DATABASE cert;
```

### 4. MinIO/Object Storage Setup

Ensure you have access to MinIO or S3-compatible object storage with the following credentials:

- Endpoint URL
- Access Key
- Secret Key
- Bucket name

## ⚙️ Configuration

### Environment Variables

Create a `.env` file or set the following environment variables:

```bash
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_DATABASE=cert

# Domain Configuration (optional)
DOMAIN=yourdomain.com

# MinIO/Storage Configuration (set in code)
STORAGE_ENDPOINT=your-minio-endpoint.com
STORAGE_ACCESS_KEY=your-access-key
STORAGE_SECRET_KEY=your-secret-key
STORAGE_BUCKET=your-bucket-name
```

## 🚀 Usage

### Basic Example

```go
package main

import (
    "log"

    "github.com/secnex/certmanager"
    "github.com/secnex/certmanager/database"
    "github.com/secnex/certmanager/store"
)

func main() {
    // Initialize object storage
    storage, err := store.NewStorage(
        "your-minio-endpoint.com",
        "your-access-key",
        "your-secret-key",
        "your-bucket-name",
    )
    if err != nil {
        log.Fatalf("Failed to create storage: %v", err)
    }

    // Initialize certificate manager
    manager := certmanager.NewCertManager(
        database.NewConnection(
            "localhost",
            5432,
            "postgres",
            "postgres",
            "cert",
        ),
        storage,
    )

    // Start the server
    manager.RunServer()
}
```

### Running with Environment Variables

```go
package main

import (
    "log"

    "github.com/secnex/certmanager"
    "github.com/secnex/certmanager/database"
    "github.com/secnex/certmanager/store"
)

func main() {
    // Initialize with environment variables
    storage, err := store.NewStorage(
        os.Getenv("STORAGE_ENDPOINT"),
        os.Getenv("STORAGE_ACCESS_KEY"),
        os.Getenv("STORAGE_SECRET_KEY"),
        os.Getenv("STORAGE_BUCKET"),
    )
    if err != nil {
        log.Fatalf("Failed to create storage: %v", err)
    }

    manager := certmanager.NewCertManager(
        database.NewConnectionFromEnv(),
        storage,
    )

    manager.RunServer()
}
```

## 🌐 API Endpoints

### Health Check

```
GET /healthz
```

Returns server health status.

**Response:**

```
Status: 200 OK
Body: OK
```

## 📁 Project Structure

```
certmanager/
├── .example/             # Example usage
│   ├── main.go
│   ├── go.mod
│   └── go.sum
├── common/               # Common utilities
│   ├── account/          # Account management
│   └── certificate/      # Certificate management
├── database/             # Database layer
│   └── database.go
├── logger/               # HTTP logging middleware
├── manager/              # Certificate manager core
├── models/               # Data models
│   ├── account.go
│   └── certificate.go
├── server/               # HTTP server
│   └── api.go
├── store/                # Object storage layer
│   └── store.go
├── go.mod
├── go.sum
└── README.md
```

## 🔧 Core Components

### Database Layer

- **GORM Integration**: Object-relational mapping for PostgreSQL
- **Auto Migration**: Automatic database schema management
- **Connection Management**: Robust database connection handling

### Storage Layer

- **MinIO Integration**: Secure object storage for certificates and private keys
- **PEM Encoding**: Proper encoding/decoding of cryptographic materials
- **Bucket Management**: Automatic bucket creation and management

### Certificate Management

- **ACME Protocol**: Full Let's Encrypt integration
- **HTTP-01 Challenge**: Domain validation support
- **RSA Key Generation**: Secure private key generation
- **Certificate Storage**: Encrypted storage of certificates and keys

### API Server

- **Gorilla Mux**: HTTP routing and middleware
- **Auto HTTPS**: Automatic SSL certificate provisioning
- **Health Checks**: Built-in monitoring endpoints
- **Request Logging**: Comprehensive HTTP request logging

## 🔒 Security Features

- **Private Key Security**: All private keys stored encrypted in object storage
- **PEM Formatting**: Industry-standard certificate formatting
- **ACME Compliance**: Full Let's Encrypt protocol compliance
- **Secure Defaults**: Secure configuration defaults throughout

## 🚀 Deployment

### Docker Deployment

Create a `Dockerfile`:

```dockerfile
FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o certmanager .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/certmanager .
CMD ["./certmanager"]
```

### Environment Configuration

```yaml
# docker-compose.yml
version: "3.8"
services:
  certmanager:
    build: .
    environment:
      - DB_HOST=postgres
      - DB_PORT=5432
      - DB_USER=postgres
      - DB_PASSWORD=postgres
      - DB_DATABASE=cert
      - DOMAIN=yourdomain.com
    ports:
      - "80:80"
      - "443:443"
    depends_on:
      - postgres
      - minio

  postgres:
    image: postgres:15
    environment:
      - POSTGRES_DB=cert
      - POSTGRES_USER=postgres
      - POSTGRES_PASSWORD=postgres
    volumes:
      - postgres_data:/var/lib/postgresql/data

  minio:
    image: minio/minio
    command: server /data --console-address ":9001"
    environment:
      - MINIO_ROOT_USER=minioadmin
      - MINIO_ROOT_PASSWORD=minioadmin
    volumes:
      - minio_data:/data
    ports:
      - "9000:9000"
      - "9001:9001"

volumes:
  postgres_data:
  minio_data:
```

## 📊 Monitoring

The application includes built-in health check endpoints:

- **Health Check**: `GET /healthz` - Returns server status
- **Request Logging**: All HTTP requests are logged with timestamps

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🔄 Changelog

### v0.1.0

- Initial release
- Basic certificate management
- Let's Encrypt integration
- MinIO storage support
- PostgreSQL database support
- HTTP/HTTPS API server

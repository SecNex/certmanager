# Certificate Manager

A comprehensive Go-based certificate management system that automates SSL/TLS certificate provisioning, management, and storage using Let's Encrypt ACME protocol.

## 🚀 Features

- **Automated Certificate Management**: Seamless SSL/TLS certificate provisioning through Let's Encrypt
- **Multiple Challenge Types**: HTTP-01 (automatic) and DNS-01 (manual) challenge support
- **Wildcard Certificates**: Support for `*.example.com` certificates via DNS challenge
- **JWT Authentication**: Secure API endpoints with JWT token validation
- **Selective Protection**: Public health endpoints with protected management endpoints
- **Secure Storage**: Private keys and certificates stored securely in MinIO object storage
- **Database Integration**: PostgreSQL database for account and certificate metadata
- **HTTP/HTTPS API Server**: RESTful API with automatic HTTPS certificate provisioning
- **Account Management**: ACME account creation and management
- **Flexible Deployment**: Works with public servers and private networks
- **Comprehensive Logging**: HTTP request logging for all endpoints
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

# Authentication Configuration (optional)
AUTH_ENABLED=true
SECNEX_GATEWAY_PUBLIC_KEY=/path/to/your/public.key

# MinIO/Storage Configuration (set in code)
STORAGE_ENDPOINT=your-minio-endpoint.com
STORAGE_ACCESS_KEY=your-access-key
STORAGE_SECRET_KEY=your-secret-key
STORAGE_BUCKET=your-bucket-name
```

## 🚀 Usage

### Basic Example (HTTP Challenge)

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

### Certificate Creation with DNS Challenge

```go
package main

import (
    "log"

    "github.com/secnex/certmanager/common/account"
    "github.com/secnex/certmanager/common/certificate"
    "github.com/secnex/certmanager/database"
    "github.com/secnex/certmanager/store"
)

func main() {
    // Initialize storage and database
    storage, err := store.NewStorage(
        "your-minio-endpoint.com",
        "your-access-key",
        "your-secret-key",
        "your-bucket-name",
    )
    if err != nil {
        log.Fatalf("Failed to create storage: %v", err)
    }

    db := database.NewConnectionFromEnv()

    // Create or get ACME account
    acc, err := account.NewAccount("your-email@example.com", db, storage)
    if err != nil {
        log.Fatalf("Failed to create account: %v", err)
    }

    // Configure DNS challenge
    config := &certificate.CertificateConfig{
        ChallengeType: certificate.ChallengeTypeDNS,
        DNSProvider:   "manual",
    }

    // Create certificate with DNS challenge
    domains := []string{"example.com", "*.example.com"}
    cert, err := certificate.NewCertificateWithConfig(domains, acc, storage, config)
    if err != nil {
        log.Fatalf("Failed to create certificate: %v", err)
    }

    log.Printf("Certificate created successfully for domains: %v", cert.Domains)
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

## 🔐 Challenge Types

The certificate manager supports two types of ACME challenges for domain validation:

### HTTP-01 Challenge (Default)

**When to use:**

- Standard domain validation
- Server is publicly accessible on port 80
- Simple setup without DNS configuration

**How it works:**

1. Let's Encrypt provides a token
2. Certificate manager serves the token on `http://yourdomain.com/.well-known/acme-challenge/TOKEN`
3. Let's Encrypt validates the token
4. Certificate is issued

**Setup:**

- Ensure port 80 is accessible from the internet
- Domain must point to your server's IP address
- No additional configuration required

### DNS-01 Challenge (Manual)

**When to use:**

- Wildcard certificates (`*.example.com`)
- Server is behind firewall/not publicly accessible
- Internal networks or private deployments

**How it works:**

1. Let's Encrypt provides a DNS challenge token
2. You manually create a TXT record: `_acme-challenge.yourdomain.com`
3. Let's Encrypt validates the DNS record
4. Certificate is issued

**Setup Process:**

1. **Run the certificate creation code** (see example above)
2. **Monitor the logs** - you'll see output like:

   ```
   Please create the following DNS TXT record:
   Name: _acme-challenge.example.com
   Value: abc123def456...

   Waiting for DNS propagation...
   ```

3. **Create the DNS TXT record** in your DNS provider:

   - **Record Type**: TXT
   - **Name**: `_acme-challenge.example.com` (replace with your domain)
   - **Value**: The token provided in the logs
   - **TTL**: 300 seconds (5 minutes) or lower if possible

4. **Wait for DNS propagation** (usually 1-5 minutes)

5. **Verify the DNS record** (optional):

   ```bash
   # Linux/macOS
   dig TXT _acme-challenge.example.com

   # Windows
   nslookup -type=TXT _acme-challenge.example.com
   ```

6. **Continue the process** - the certificate manager will automatically verify and complete the process

**Example DNS Record:**

```
Type: TXT
Name: _acme-challenge.example.com
Value: J7tUF9F8F8F8F8F8F8F8F8F8F8F8F8F8F8F8F8F8F8F8
TTL: 300
```

**For Wildcard Certificates:**

```go
domains := []string{"example.com", "*.example.com"}
// This will require DNS challenge for the wildcard domain
```

## 🔐 Authentication

The certificate manager includes a comprehensive authentication system using JWT tokens. Authentication is selectively applied to protect sensitive endpoints while keeping essential monitoring endpoints publicly accessible.

### Authentication Flow

1. **JWT Token Validation**: All protected endpoints require a valid JWT token
2. **Bearer Token**: Include the token in the `Authorization` header
3. **Public Key Verification**: Tokens are verified using a configurable public key
4. **Selective Protection**: Only sensitive endpoints require authentication

### Configuration

Set the following environment variables to enable authentication:

```bash
# Enable authentication
AUTH_ENABLED=true

# Path to the public key file for JWT verification
SECNEX_GATEWAY_PUBLIC_KEY=/path/to/your/public.key
```

**Note**: If `AUTH_ENABLED` is not set to `true` or `1`, authentication is disabled and all endpoints are publicly accessible.

### Public Key Setup

1. **Generate a key pair** (if you don't have one):

   ```bash
   # Generate private key
   openssl genrsa -out private.key 2048

   # Generate public key
   openssl rsa -in private.key -pubout -out public.key
   ```

2. **Configure the public key path**:
   ```bash
   export SECNEX_GATEWAY_PUBLIC_KEY=/path/to/public.key
   ```

### Making Authenticated Requests

Include the JWT token in the `Authorization` header:

```bash
# Example authenticated request
curl -H "Authorization: Bearer YOUR_JWT_TOKEN" \
     https://your-domain.com/test
```

**JWT Token Structure:**

```json
{
  "alg": "RS256",
  "typ": "JWT"
}
{
  "sub": "user_id",
  "iat": 1640995200,
  "exp": 1640998800
}
```

## 🌐 API Endpoints

### Public Endpoints (No Authentication Required)

#### Health Check

```
GET /healthz
```

Returns server health status. This endpoint is publicly accessible for monitoring purposes.

**Response:**

```
Status: 200 OK
Body: OK
```

**Example:**

```bash
curl https://your-domain.com/healthz
```

### Protected Endpoints (Authentication Required)

All protected endpoints require a valid JWT token in the `Authorization` header.

#### Test Endpoint

```
GET /test
```

A simple test endpoint to verify authentication is working.

**Headers:**

```
Authorization: Bearer YOUR_JWT_TOKEN
```

**Response:**

```
Status: 200 OK
Body: This is a test!
```

**Example:**

```bash
curl -H "Authorization: Bearer YOUR_JWT_TOKEN" \
     https://your-domain.com/test
```

### Authentication Errors

#### 401 Unauthorized

Returned when:

- No `Authorization` header is provided
- Invalid JWT token
- Expired JWT token
- Token signature verification fails

**Response:**

```json
{
	"error": "Unauthorized"
}
```

#### 500 Internal Server Error

Returned when:

- Public key file cannot be read
- Server configuration error

**Response:**

```json
{
	"error": "Internal Server Error"
}
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
- **HTTP-01 Challenge**: Automatic domain validation for public servers
- **DNS-01 Challenge**: Manual DNS validation for wildcard certificates and private networks
- **RSA Key Generation**: Secure private key generation
- **Certificate Storage**: Encrypted storage of certificates and keys
- **Wildcard Support**: `*.example.com` certificates via DNS challenge

### API Server

- **Gorilla Mux**: HTTP routing and middleware
- **JWT Authentication**: Secure endpoint protection with public key verification
- **Selective Middleware**: Public health endpoints with protected management routes
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

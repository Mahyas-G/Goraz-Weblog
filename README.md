# Goraz Weblog Application

A secure **server-side rendered (SSR) weblog application** built with **Go**, **Echo**, and **PostgreSQL**.

The application supports authentication, public and private weblog posts, private post sharing, comments, image uploads, and secure session-based access control.

---

## Live Deployment

🚀 **Production URL**

`https://weblog-app.up.railway.app`

---

# Features

* User registration, login, and logout
* Secure server-side sessions using random tokens
* Public and private weblog posts
* Sharing private posts with selected users
* Personalized feed based on access permissions
* Comments on accessible posts
* Post deletion restricted to owners
* Image upload support:

  * JPEG, PNG, WebP
  * Real file-content validation
  * Maximum size limit
  * Randomized filenames
  * Automatic cleanup of unused files
* CSRF protection
* Secure cookies
* Docker-based deployment support

> Post editing was intentionally excluded because it was outside the project requirements.

---

# Tech Stack

| Technology    | Purpose               |
| ------------- | --------------------- |
| Go 1.22+      | Backend               |
| Echo v4       | Web framework         |
| PostgreSQL 16 | Database              |
| sqlx + pgx    | Database access       |
| goose         | Database migrations   |
| html/template | Server-side rendering |
| bcrypt        | Password hashing      |
| log/slog      | Structured logging    |
| Docker        | Containerization      |
| Railway       | Cloud deployment      |

---

# Architecture

The project follows a layered architecture:

```text
HTTP Request

      ↓

Handler

      ↓

Service

      ↓

Repository

      ↓

PostgreSQL
```

### Layers

**Handler**

* Handles HTTP requests
* Validates input
* Renders HTML templates

**Service**

* Contains business logic
* Implements authorization rules

**Repository**

* Handles database operations

**Model**

* Defines application data structures

---

# Project Structure

```text
weblog/
├── cmd/
│   └── server/
│       └── main.go
│
├── internal/
│   ├── auth/
│   ├── config/
│   ├── db/
│   │   └── migrations/
│   ├── handler/
│   ├── middleware/
│   ├── model/
│   ├── repository/
│   ├── service/
│   ├── upload/
│   └── validation/
│
├── web/
│   ├── templates/
│   └── static/
│
├── Dockerfile
├── docker-compose.yml
├── go.mod
└── README.md
```

---

# Running Locally

## Requirements

Install:

* Go 1.22+
* Docker Desktop
* Docker Compose
* goose migration tool

Install goose:

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

---

## Environment Configuration

Create `.env`:

```env
DATABASE_URL=postgres://postgres:postgres@localhost:5432/weblog?sslmode=disable
PORT=8080
COOKIE_SECURE=false
```

---

# Run with Docker Compose (Recommended)

Start application and PostgreSQL:

```bash
docker compose up --build
```

Run migrations:

```bash
goose -dir internal/db/migrations postgres \
"postgres://postgres:postgres@localhost:5432/weblog?sslmode=disable" up
```

Open:

```text
http://localhost:8080
```

---

# Run Without Docker

Start PostgreSQL manually, then:

```bash
go mod tidy
```

Run migrations:

```bash
goose -dir internal/db/migrations postgres "$DATABASE_URL" up
```

Start application:

```bash
go run ./cmd/server
```

---

# Testing

Run:

```bash
go vet ./...
go test ./...
```

---

# Railway Deployment

The application is prepared for deployment using Docker and Railway.

Deployment architecture:

```text
GitHub Repository

        ↓

Railway Service

        ↓

Docker Container

        ↓

Railway PostgreSQL

        ↓

Production Application
```

## Deployment Steps

### 1. Create Railway Project

Connect the GitHub repository:

```
New Project
    ↓
Deploy from GitHub Repository
```

Railway automatically detects the Dockerfile.

---

### 2. Add PostgreSQL

Create a PostgreSQL service:

```
New
 ↓
Database
 ↓
PostgreSQL
```

---

### 3. Configure Environment Variables

Set:

```env
DATABASE_URL=<Railway PostgreSQL connection string>
COOKIE_SECURE=true
```

---

### 4. Run Database Migrations

Execute:

```bash
goose -dir internal/db/migrations postgres "$DATABASE_URL" up
```

---

### 5. Enable Persistent Upload Storage

Create a Railway Volume and mount it to:

```text
web/static/uploads
```

so uploaded images remain available after redeployment.

---

### 6. Deploy

Push changes to GitHub:

```bash
git push
```

Railway automatically builds and deploys the application.

---

# Security Highlights

* Passwords are hashed using bcrypt
* Sessions use secure random tokens
* Session data is stored server-side
* CSRF protection is enabled
* Secure cookie configuration
* Private posts are protected by authorization checks
* Only owners can delete posts
* Uploaded files are validated by actual content
* Request body limits are enforced

---

# Project Requirements

Developed as the final project for:

**Goraz Final Module - APA Bootcamp**

Implemented requirements:

* Authentication
* Public and private posts
* Private post sharing
* Access control
* Comments
* Image uploads
* Cloud deployment configuration

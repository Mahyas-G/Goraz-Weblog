# Weblog Application

A secure server-side rendered (SSR) weblog application built with **Go**, **Echo**, and **PostgreSQL**.

The application supports user authentication, public and private posts, private post sharing, comments, optional image uploads, and secure session-based access control.

## Features

* User registration, login, and logout
* Secure server-side sessions using opaque random tokens
* Create public or private weblog posts
* Share private posts with selected users by username
* Personalized feed based on user access permissions
* View post details, images, and comments
* Add comments to accessible posts
* Delete posts only by their owners
* Optional image upload with:

  * JPEG, PNG, and WebP support
  * Real file-content validation
  * 5 MB maximum file size
  * Randomized filenames
  * Automatic cleanup of unused files
* CSRF protection for all state-changing requests
* Secure cookies using `HttpOnly`, `SameSite=Lax`, and optional `Secure`
* Automatic cleanup of expired sessions
* Docker support and Fly.io deployment configuration

> Post editing is intentionally not included because it is outside the project requirements.

## Tech Stack

| Technology      | Purpose               |
| --------------- | --------------------- |
| Go 1.22+        | Backend language      |
| Echo v4         | Web framework         |
| PostgreSQL 16   | Database              |
| `sqlx` + `pgx`  | Database access       |
| `goose`         | Database migrations   |
| `html/template` | Server-side rendering |
| `bcrypt`        | Password hashing      |
| `log/slog`      | Structured logging    |
| Docker          | Containerization      |
| Fly.io          | Deployment            |

## Architecture

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

* **Handler:** Processes HTTP requests and renders templates.
* **Service:** Contains business rules, authorization, and application logic.
* **Repository:** Handles database operations and SQL queries.
* **Model:** Defines the application's data structures.

## Project Structure

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
│   ├── logger/
│   ├── middleware/
│   ├── model/
│   ├── repository/
│   ├── service/
│   ├── upload/
│   └── validation/
│
├── web/
│   ├── static/
│   │   ├── css/
│   │   └── uploads/
│   └── templates/
│       ├── auth/
│       ├── layouts/
│       ├── partials/
│       └── weblog/
│
├── .env.example
├── .gitignore
├── Dockerfile
├── docker-compose.yml
├── fly.toml
├── go.mod
├── go.sum
└── README.md
```

## Prerequisites

* Go 1.22 or newer
* PostgreSQL 16
* Docker and Docker Compose (optional)
* `goose` for database migrations

Install Goose:

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

## Environment Variables

Create a local environment file from the example:

```bash
cp .env.example .env
```

Configure the following variables:

| Variable        | Required | Default | Description                    |
| --------------- | -------: | ------- | ------------------------------ |
| `DATABASE_URL`  |      Yes | —       | PostgreSQL connection URL      |
| `PORT`          |       No | `8080`  | Application port               |
| `COOKIE_SECURE` |       No | `false` | Set to `true` when using HTTPS |

Example:

```env
DATABASE_URL=postgres://postgres:password@localhost:5432/weblog?sslmode=disable
PORT=8080
COOKIE_SECURE=false
```

## Run Locally

Install dependencies:

```bash
go mod tidy
```

Load the environment variables:

```bash
export $(grep -v '^#' .env | xargs)
```

Run the database migrations:

```bash
goose -dir internal/db/migrations postgres "$DATABASE_URL" up
```

Start the application:

```bash
go run ./cmd/server
```

Open:

```text
http://localhost:8080
```

## Run with Docker

Start the application and PostgreSQL:

```bash
docker-compose up --build
```

Then run the migrations:

```bash
goose -dir internal/db/migrations postgres \
"postgres://postgres:postgres@localhost:5432/weblog?sslmode=disable" up
```

## Testing

Run the following checks:

```bash
go vet ./...
go test ./...
```

The project currently includes table-driven unit tests for the validation package.

## Manual Test Flow

1. Open `/signup` and create a user account.
2. Log in through `/login`.
3. Create a public or private post from `/weblog/new`.
4. For a private post, enter the usernames allowed to access it.
5. Open a post from the feed to view its content, image, and comments.
6. Add a comment to a post you are allowed to view.
7. Verify that only the post owner can see and use the delete option.
8. Log out using the navigation bar.

## Security Highlights

* Passwords are hashed using `bcrypt`.
* Sessions use cryptographically secure random tokens.
* Session tokens are stored server-side in PostgreSQL.
* CSRF protection is enabled for all POST forms.
* Session cookies use `HttpOnly` and `SameSite=Lax`.
* HTTPS deployments enable secure cookies through `COOKIE_SECURE=true`.
* Private posts are protected by server-side authorization checks.
* Post deletion is restricted to the original author.
* Uploaded images are validated by decoding their actual content, not only by checking file extensions.
* Request body size is globally limited.

## Deployment

The project includes:

* Multi-stage `Dockerfile`
* `docker-compose.yml`
* `fly.toml`
* Persistent storage configuration for uploaded images

For deployment, configure the production database URL and deploy using:

```bash
fly deploy
```

## Project Requirements

This application was developed as the final project for the **Goraz Final Module - APA Bootcamp**.

The project implements the required weblog functionality, including:

* Authentication
* Public and private posts
* Private post sharing
* Access control
* Comments
* Image uploads
* Deployment configuration

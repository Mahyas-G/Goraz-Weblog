# Goraz Weblog Application

A secure **server-side rendered (SSR) weblog application** built with **Go**, **Echo**, and **PostgreSQL**.

The application supports authentication, public and private weblog posts, private post sharing, comments, image uploads, and secure session-based access control.

---

## Live Deployment

🚀 **Production URL**

`https://goraz-weblog.up.railway.app`

---

## Features

### Authentication

* User signup and login
* Secure password hashing using bcrypt
* Database-backed sessions
* Logout functionality

### Weblog

* Create blog posts
* Public and private posts
* Share private posts with selected users
* View feed based on user permissions
* View detailed posts
* Delete posts only by their owners

### Comments

* Authenticated users can comment on accessible posts
* Comments include author information and timestamps

### Image Upload

* Optional image upload for posts
* Supported formats:

  * JPEG
  * PNG
  * WebP
* File validation and cleanup handling

### Security

* CSRF protection
* Secure session handling
* Input validation
* Ownership and access control checks
* Protected file upload validation

---

## Tech Stack

| Component        | Technology            |
| ---------------- | --------------------- |
| Language         | Go 1.22+              |
| Framework        | Echo v4               |
| Database         | PostgreSQL            |
| Database Driver  | pgx + sqlx            |
| Migration Tool   | goose                 |
| Templates        | html/template (SSR)   |
| Authentication   | Cookie-based sessions |
| Password Hashing | bcrypt                |
| Deployment       | Docker                |

---

## Project Structure

```
.
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

## Run Locally

### Requirements

* Go 1.22+
* PostgreSQL
* goose (for migrations)

### 1. Clone repository

```bash
git clone <repository-url>
cd Goraz-Weblog-main
```

### 2. Configure environment

Create a `.env` file:

```env
DATABASE_URL=postgres://user:password@localhost:5432/weblog
PORT=8080
COOKIE_SECURE=false
```

### 3. Run migrations

```bash
goose -dir internal/db/migrations postgres "$DATABASE_URL" up
```

### 4. Start application

```bash
go run ./cmd/server
```

Application will be available at:

```
http://localhost:8080
```

---

## Docker Deployment

Run the application with PostgreSQL:

```bash
docker compose up --build
```

The application is configured for container-based deployment.

---

## Testing

Run tests:

```bash
go test ./...
```

Run static analysis:

```bash
go vet ./...
```

---

## Deployment

The application can be deployed using platforms such as:

* Railway
* Render
* Fly.io

Deployment requires:

* PostgreSQL database
* Environment variables
* Persistent storage for uploaded images

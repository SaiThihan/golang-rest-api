# go-basic

REST API built with Go, chi, and PostgreSQL.

## Quick Start

cp .env.example .env  # edit DATABASE_URL
docker-compose up -d
go run .

## Endpoints

| Method | Path | Auth |
|--------|------|------|
| GET | `/health` | No |
| POST | `/register` | No |
| POST | `/auth/token` | No |
| GET | `/posts` | No |
| GET | `/posts/{id}` | Bearer |
| POST | `/posts` | Bearer |
| PUT | `/posts/{id}` | Bearer |
| DELETE | `/posts/{id}` | Bearer |

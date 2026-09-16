# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Personal portfolio website built with Go, Gin web framework, and PostgreSQL. Serves server-rendered HTML pages for home, projects, blog, and contact sections.

## Build & Run Commands

```bash
# Install dependencies
go mod tidy

# Run the server (starts on :8080)
go run main.go

# Build binary
go build -o portfolio .
```

## Database Setup

Requires a running PostgreSQL instance. Tables are auto-created on startup via `database.Connect()`.

Connection configured via environment variables (defaults in parentheses):
- `DB_HOST` (localhost), `DB_PORT` (5432), `DB_USER` (portfolio), `DB_PASSWORD` (portfolio), `DB_NAME` (portfolio)

## Architecture

**Routing**: All routes defined in `main.go` using Gin. Templates loaded via `r.LoadHTMLGlob("templates/*")`.

**Layers**:
- `handlers/` — Gin handler functions, one file per page area. Each handler queries the DB directly using `database.DB` and renders an HTML template.
- `models/` — Plain Go structs (`BlogPost`, `Project`, `ContactMessage`) with no methods or ORM. Fields map 1:1 to database columns.
- `database/` — Single `db.go` that opens the PostgreSQL connection, exposes a global `DB *sql.DB`, and runs `CREATE TABLE IF NOT EXISTS` on startup.
- `templates/` — Go `html/template` files. Every page template calls `{{ template "layout" . }}` and defines a `{{ define "content" }}` block. Shared chrome (nav, footer) lives in `layout.html`.
- `static/` — CSS served at `/static`.

**Key pattern**: There is no service/repository layer. Handlers use raw `database/sql` queries directly against `database.DB`. Keep this flat structure unless the project grows significantly.

**Template data**: Handlers pass `gin.H` maps with a `"title"` key (used by layout) plus page-specific data (`"posts"`, `"projects"`, `"post"`, `"error"`, `"success"`).

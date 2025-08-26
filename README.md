# OpenB2B

## Overview

**OpenB2B** is an open-source, multi-tenant ordering and invoicing platform designed for small wholesalers to manage products, process orders, and track invoices — all in one place.

It’s ideal for businesses like book wholesalers, restaurant supply vendors, event merchandise suppliers, or generic B2B product distributors. OpenB2B ensures accurate order totals, organized product management, and seamless invoice generation with multi-user and role-based access per organization.

### ✨ Features

- ✅ **Order Management** — Track orders, order items, and their statuses with relational integrity.  
- 📦 **Product Management** — Add/edit/delete products with SKU, category, price, inventory, and images.  
- 🧾 **Invoice Generation** — Auto-generate PDF invoices, download or email to customers.  
- 🔒 **Multi-Tenancy & Roles** — Org-level accounts with Owner, Admin, Sales, and Viewer roles.  
- 📊 **Basic Reporting** — View total sales, top products, and orders per customer.  
- 🔍 **Customer Management** — Store customer details and track order history.  

---

## API Documentation

Detailed API reference and interactive Swagger UI can be found in [api.md](./docs/api.md).

--

## Technology Stack

- **Frontend:** React + TypeScript  
- **Build Tool:** Vite  
- **Styling:** TailwindCSS  
- **Authentication:** Clerk (handles signup, login, OAuth, and sessions)  
- **Routing:** React Router  
- **State Management:** React Context API  
- **Backend:** Golang, Chi  
- **Database:** PostgreSQL, GORM  

---

## Prerequisites

### Frontend

- Node.js (v22.0.0 or higher, recommended v22.11.0)  
- npm (v11.0.0 or higher, recommended v11.2.0) or yarn  

### Backend

- Golang (v1.23.4 or higher)  
- Chi (v5.2.1 or higher)  
- PostgreSQL (v15.4 or higher)  
- Redis (v7.2.0 or higher, optional for caching/session)  

---

## Authentication

OpenB2B uses [Clerk](https://clerk.com) for authentication and session management. Clerk handles **sign-up, login, OAuth, and session handling**, while the backend manages user and organization data via webhooks and custom session claims.

For full setup instructions, see the [Authentication Documentation](docs/authentication.md).

---

## Installation

### Clone the repository

```bash
git clone https://github.com/yourusername/openb2b.git
cd openb2b
```

### Install frontend dependencies

```bash
cd frontend && pnpm install
```

### Install backend dependencies

```bash
cd api && pnpm install
```


### Configure environment variables

Create a `.env` file in the **frontend** directory:

```env
VITE_API_URL=your_api_url
VITE_OAUTH_CLIENT_ID=your_oauth_client_id
```

Create a `.env` file in the **backend** directory:

```env
PORT=
APP_ENV=
DB_HOST=
DB_PORT=
DB_NAME=
DB_USER=
DB_PASSWORD=
CLERK_SECRET_KEY=
CLERK_WEBHOOK_SIGNING_SECRET=

#goose
GOOSE_DRIVER=
GOOSE_DBSTRING=
GOOSE_MIGRATION_DIR=
```

### Start development servers

**Frontend:**

```bash
cd frontend
npm run dev
# or
yarn dev
```

**Backend:**

```bash
cd backend
go run main.go
```

---

## 🐦 Database Migrations Setup

OpenB2B uses **Goose** to manage database migrations via SQL files.

### Prerequisites

* Go 1.18+
* PostgreSQL running locally or remotely
* set environment variables 

```env
GOOSE_DRIVER=postgres
GOOSE_DBSTRING=postgres://admin:admin@localhost:5432/openb2b
GOOSE_MIGRATION_DIR=./db/migrations
```

### Install Goose CLI

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
export PATH="$PATH:$(go env GOPATH)/bin"
```

### Migration Commands

```bash
# Create a new migration
make migrate-new name=create_users_table

# Run migrations
make migrate-up

# Rollback migrations
make migrate-down

# Check migration status
make migrate-status
```

---

## Support and Feedback

For support, feature requests, or feedback, contact:

* **Email:** [support@openb2b.com](mailto:support@openb2b.com)
* **Phone:** +234-8067177670
* **Website:** [https://openb2b.com/support](https://openb2b.com/support)

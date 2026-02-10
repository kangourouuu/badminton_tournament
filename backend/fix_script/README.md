# Database Migration Fix Script

This script is designed to run database migrations in environments where standard Go tools fail due to SSL/TLS certificate issues (e.g., "x509: certificate signed by unknown authority").

## Prerequisites

- Node.js (v18+)
- npm

## Setup

1.  Navigate to this directory:

    ```bash
    cd backend/fix_script
    ```

2.  Install dependencies:
    ```bash
    npm install
    ```

## Usage

You must provide the `DATABASE_URL` environment variable.

### Option 1: If you have the URL

```bash
export DATABASE_URL="postgres://user:pass@host:5432/dbname?sslmode=require"
node migrate.js
```

### Option 2: Inline

```bash
DATABASE_URL="postgres://..." node migrate.js
```

## What it does

It connects to the Postgres database using Node.js's `pg` driver with `rejectUnauthorized: false` to bypass strict SSL verification, and executes the necessary SQL `ALTER TABLE` statements to add missing columns (`pool` for teams/groups, `sets_detail` for matches).

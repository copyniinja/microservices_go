# Database Migrations (Postgres)

We use **[golang-migrate](https://github.com/golang-migrate/migrate)** to manage database schema and seed data for the auth service.

---

## Step 1: Install `migrate` CLI

> If not already installed

```bash
# Linux
curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.0/migrate.linux-amd64.tar.gz | tar xvz
sudo mv migrate /usr/local/bin/
```

Verify installation:

```bash
migrate -version
```

---

## Step 2: Run Migrations

From your project root:

```bash
migrate -path ./migrations -database "postgres://<DB_USER>:<DB_PASSWORD>@<DB_HOST>:<DB_PORT>/<DB_NAME>?sslmode=disable" up
```

### Parameters

- `-path` → folder where migration files are located
- `-database` → your Postgres connection string
- `up` → applies all pending migrations

After running, the database will have the **`users` table** and initial seeded users.

---

## Step 3: Rollback Migrations (Optional)

To rollback the last applied migration:

```bash
migrate -path ./migrations -database "postgres://<DB_USER>:<DB_PASSWORD>@<DB_HOST>:<DB_PORT>/<DB_NAME>?sslmode=disable" down
```

> Rolls back **one migration at a time**. Repeat `down` to step back multiple migrations.

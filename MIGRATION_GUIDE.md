# Database Migration Guide

Sistem WMS menggunakan **golang-migrate** untuk manajemen database migrations, mirip dengan Laravel's Artisan migrate command.

## 🚀 Quick Start

### 1. **Setup Environment**

```bash
cp .env.example .env
# Edit .env dengan PostgreSQL credentials
```

### 2. **Install Dependencies**

```bash
go mod download
```

### 3. **Run Migrations**

```bash
go run main.go migrate
```

---

## 📋 Available Commands

### **Run All Pending Migrations**

```bash
go run main.go migrate
```

Menjalankan semua migration files yang belum di-apply ke database.

**Output:**

```
🔄 Running migrations...
✅ All migrations completed successfully
```

### **Rollback Last Migration**

```bash
go run main.go migrate:rollback
```

Membatalkan migration terakhir yang sudah di-apply.

**Output:**

```
⏮️  Rolling back last migration...
✅ Migration rolled back successfully
```

### **Check Migration Status**

```bash
go run main.go migrate:status
```

Menampilkan versi migration terakhir yang di-apply dan status database.

**Output:**

```
📊 Migration Status: Version 9 (clean)
```

### **Fresh Migrations** (⚠️ Dangerous!)

```bash
go run main.go migrate:fresh
```

**Drops semua tables** dan menjalankan kembali semua migrations dari awal.

**Use Case**: Local development environment reset.

**Output:**

```
🔄 Refreshing database (this will drop all tables)...
Are you sure? (yes/no): yes
✅ All tables dropped successfully
✅ Database refreshed and all migrations applied
```

### **Seed Database** (Coming Soon)

```bash
go run main.go migrate:seed
```

Will populate database dengan sample data.

### **Start Server** (Default)

```bash
go run main.go
go run main.go serve
```

Menjalankan HTTP server (otomatis menjalankan pending migrations).

---

## 📁 Migration Files Structure

Migration files tersimpan di `db/migrations/` dengan naming convention:

```
db/migrations/
├── 001_create_users_table.up.sql
├── 001_create_users_table.down.sql
├── 002_create_taxes_table.up.sql
├── 002_create_taxes_table.down.sql
├── 003_create_categories_stickers_classes_buyers_tables.up.sql
├── 003_create_categories_stickers_classes_buyers_tables.down.sql
├── 004_create_stores_crews_racks_tables.up.sql
├── 004_create_stores_crews_racks_tables.down.sql
├── 005_create_products_tables.up.sql
├── 005_create_products_tables.down.sql
├── 006_create_cargos_bags_tables.up.sql
├── 006_create_cargos_bags_tables.down.sql
├── 007_create_orders_tables.up.sql
├── 007_create_orders_tables.down.sql
├── 008_create_store_transfers_user_class_logs_tables.up.sql
├── 008_create_store_transfers_user_class_logs_tables.down.sql
├── 009_create_slow_moving_tables.up.sql
└── 009_create_slow_moving_tables.down.sql
```

### **Naming Convention**

- **Prefix**: `NNN_` (3 digits, auto-incrementing)
- **Description**: Snake_case description
- **Suffix**: Either `.up.sql` (apply) or `.down.sql` (rollback)

**Example**:

- `010_add_indexes_to_orders.up.sql` - Apply migration
- `010_add_indexes_to_orders.down.sql` - Rollback migration

---

## ✅ Current Migrations

| #   | Migration                             | Status                                  |
| --- | ------------------------------------- | --------------------------------------- |
| 001 | Create Users Table                    | ✅ Initial users setup                  |
| 002 | Create Taxes Table                    | ✅ Tax management                       |
| 003 | Categories, Stickers, Classes, Buyers | ✅ Product & customer classification    |
| 004 | Stores, Crews, Racks                  | ✅ Warehouse infrastructure             |
| 005 | Products (Document, Master, Log)      | ✅ Product catalog management           |
| 006 | Cargos & Bags                         | ✅ Inventory containers                 |
| 007 | Orders (Order, Item, Cargo)           | ✅ Order management                     |
| 008 | Store Transfers & User Class Logs     | ✅ Store operations & customer tracking |
| 009 | Slow Moving Items                     | ✅ Inventory analytics                  |

---

## 🛠️ Creating New Migrations

### **Step 1: Create Migration Files**

Generate files dengan naming convention:

```bash
# .up.sql file (apply)
db/migrations/010_add_vendor_table.up.sql

# .down.sql file (rollback)
db/migrations/010_add_vendor_table.down.sql
```

### **Step 2: Write SQL**

**010_add_vendor_table.up.sql:**

```sql
CREATE TABLE IF NOT EXISTS vendors (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255),
    phone VARCHAR(20),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_vendors_email ON vendors(email);
```

**010_add_vendor_table.down.sql:**

```sql
DROP TABLE IF EXISTS vendors CASCADE;
```

### **Step 3: Run Migration**

```bash
go run main.go migrate
```

---

## 🔄 Migration Workflow Examples

### **Scenario 1: Development - Add New Column**

```bash
# Create new migration
# db/migrations/011_add_notes_to_orders.up.sql
ALTER TABLE orders ADD COLUMN notes TEXT;

# Apply
go run main.go migrate

# Back to previous state
go run main.go migrate:rollback
```

### **Scenario 2: Production - Fresh Setup**

```bash
# Initial setup on production server
go run main.go migrate

# This will apply migrations 001-009 in order
```

### **Scenario 3: Local Development Reset**

```bash
# Full reset (drops all data)
go run main.go migrate:fresh

# Or programmatically:
go run main.go migrate:rollback  # Repeat this 9 times to rollback all
go run main.go migrate           # Apply all again
```

---

## 📊 Database Versioning

Golang-migrate menggunakan special table `schema_migrations` untuk track:

- **version**: Migration version number
- **dirty**: Boolean flag (true = migration incomplete)

**View version table:**

```sql
SELECT * FROM schema_migrations;
```

**Output:**

```
 version | dirty
---------+-------
       9 | false
```

---

## ⚠️ Important Notes

### **Do NOT:**

- ❌ Modify migration files after applying (breaking change!)
- ❌ Manually run SQL `CREATE TABLE` without migration files
- ❌ Skip migrations by manually changing `schema_migrations` table

### **DO:**

- ✅ Always create `.up.sql` AND `.down.sql` files
- ✅ Test migrations locally before production
- ✅ Use incremental version numbers (001, 002, 003...)
- ✅ Write reversible migrations (clean up in .down.sql)

---

## 🐛 Troubleshooting

### **Problem: "dirty" Migration Status**

```
📊 Migration Status: Version 5 (DIRTY - migration incomplete!)
```

**Solution:**

```bash
# Check if migration got stuck
# Then manually fix issue and re-run
go run main.go migrate
```

### **Problem: Migration File Not Found**

```
Error: file not found
```

**Solution:**

- Pastikan file ada di `db/migrations/` folder
- Pastikan naming convention benar (NNN_description.up.sql)
- Run `ls -la db/migrations/` untuk verify

### **Problem: PostgreSQL Connection Failed**

```
Error: failed to create migrator
```

**Solution:**

```bash
# Verify .env file
cat .env

# Check PostgreSQL running
# Make sure DB_HOST, DB_USER, DB_PASSWORD correct
```

---

## 📚 Related Documentation

- [GIN_POSTGRESQL_SETUP.md](../GIN_POSTGRESQL_SETUP.md) - Database setup guide
- [DEVELOPMENT_GUIDE.md](../DEVELOPMENT_GUIDE.md) - Full development guide
- [golang-migrate Docs](https://github.com/golang-migrate/migrate) - Official library

---

## ✨ Summary

| Command                           | Purpose                      |
| --------------------------------- | ---------------------------- |
| `go run main.go migrate`          | Apply all pending migrations |
| `go run main.go migrate:rollback` | Undo last migration          |
| `go run main.go migrate:status`   | Check current version        |
| `go run main.go migrate:fresh`    | Reset database completely    |
| `go run main.go`                  | Start server (auto-migrate)  |

**That's it!** Similar to Laravel's `php artisan migrate` 🚀

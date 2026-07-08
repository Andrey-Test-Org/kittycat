# kt-petstore

A Kotlin cat-adoption service backed by **Snowflake**. Manages shelters, breeds, cats, customers, carts, orders, inventory, and an audit trail.

## Layout

```
kt-petstore/
├── seed_data.sql                     # INSERT statements — 15-20 rows per table
└── src/main/kotlin/petstore/
    ├── Main.kt                       # entrypoint — demonstrates all 10 models + services
    ├── model/Models.kt               # 10 data classes (one per table)
    ├── db/Database.kt                # Snowflake JDBC connection + ResultSet helpers
    ├── repository/Repositories.kt    # 10 repository classes (CRUD via JDBC)
    └── service/Services.kt           # 7 service classes with business logic
```

## Seed Data

After running the DDL, populate all tables with sample data:

```sql
-- in a Snowflake worksheet:
USE DATABASE TEST;
USE SCHEMA PETSTORE;
@seed_data.sql
```

Or paste the contents of `seed_data.sql` directly. Inserts 15-20 rows per table with consistent foreign key references.

## 10 Tables

| # | Table | Model class | Purpose |
|---|-------|-------------|---------|
| 1 | `customers` | `Customer` | People who adopt cats |
| 2 | `breeds` | `Breed` | Cat breeds (Maine Coon, Siamese, etc.) |
| 3 | `shelters` | `Shelter` | Rescue organisations / shelters |
| 4 | `cats` | `Cat` | Cats available for adoption |
| 5 | `carts` | `Cart` | A customer's pending selection |
| 6 | `cart_items` | `CartItem` | Line items in a cart |
| 7 | `orders` | `Order` | Adoption orders with lifecycle (pending → paid → shipped / cancelled) |
| 8 | `order_items` | `OrderItem` | Line items in an order |
| 9 | `inventory` | `Inventory` | Per-cat stock levels (available / reserved) |
| 10 | `audit_log` | `AuditLog` | Append-only audit trail |

## Run

Requires the Snowflake JDBC driver on the classpath:

```
export SNOWFLAKE_ACCOUNT=youraccount
export SNOWFLAKE_USER=youruser
export SNOWFLAKE_PASSWORD=yourpassword
export SNOWFLAKE_WAREHOUSE=yourwarehouse
export SNOWFLAKE_ROLE=yourrole

kotlinc src/main/kotlin/petstore -cp snowflake-jdbc-3.x.jar -d out
java -cp out:snowflake-jdbc-3.x.jar petstore.MainKt
```

---

## Snowflake Schema Creation Script

Run this in a Snowflake worksheet. It creates the `PETSTORE` schema inside the **`TEST`** database with all 10 tables.

```sql
USE DATABASE TEST;

CREATE SCHEMA IF NOT EXISTS PETSTORE;
USE SCHEMA PETSTORE;

-- 1. customers
CREATE TABLE IF NOT EXISTS customers (
    id           STRING       NOT NULL,
    full_name    STRING       NOT NULL,
    email        STRING       NOT NULL,
    phone        STRING,
    address      STRING       NOT NULL,
    created_at   TIMESTAMP_NTZ NOT NULL DEFAULT CURRENT_TIMESTAMP(),
    updated_at   TIMESTAMP_NTZ NOT NULL DEFAULT CURRENT_TIMESTAMP(),
    PRIMARY KEY (id)
);

-- 2. breeds
CREATE TABLE IF NOT EXISTS breeds (
    id                STRING       NOT NULL,
    name              STRING       NOT NULL,
    description       STRING,
    origin_country    STRING       NOT NULL,
    avg_lifespan_years INT          NOT NULL,
    created_at        TIMESTAMP_NTZ NOT NULL DEFAULT CURRENT_TIMESTAMP(),
    updated_at        TIMESTAMP_NTZ NOT NULL DEFAULT CURRENT_TIMESTAMP(),
    PRIMARY KEY (id)
);

-- 3. shelters
CREATE TABLE IF NOT EXISTS shelters (
    id          STRING       NOT NULL,
    name        STRING       NOT NULL,
    address     STRING       NOT NULL,
    phone       STRING,
    email       STRING,
    created_at  TIMESTAMP_NTZ NOT NULL DEFAULT CURRENT_TIMESTAMP(),
    updated_at  TIMESTAMP_NTZ NOT NULL DEFAULT CURRENT_TIMESTAMP(),
    PRIMARY KEY (id)
);

-- 4. cats
CREATE TABLE IF NOT EXISTS cats (
    id           STRING       NOT NULL,
    name         STRING       NOT NULL,
    breed_id     STRING       NOT NULL,
    shelter_id   STRING       NOT NULL,
    birth_date   DATE         NOT NULL,
    price_cents  NUMBER(12,2) NOT NULL,
    currency     STRING(3)    NOT NULL,
    status       STRING       NOT NULL DEFAULT 'AVAILABLE',
    description  STRING,
    created_at   TIMESTAMP_NTZ NOT NULL DEFAULT CURRENT_TIMESTAMP(),
    updated_at   TIMESTAMP_NTZ NOT NULL DEFAULT CURRENT_TIMESTAMP(),
    PRIMARY KEY (id),
    FOREIGN KEY (breed_id)   REFERENCES breeds (id),
    FOREIGN KEY (shelter_id) REFERENCES shelters (id)
);

-- 5. carts
CREATE TABLE IF NOT EXISTS carts (
    id          STRING       NOT NULL,
    customer_id STRING       NOT NULL,
    created_at  TIMESTAMP_NTZ NOT NULL DEFAULT CURRENT_TIMESTAMP(),
    updated_at  TIMESTAMP_NTZ NOT NULL DEFAULT CURRENT_TIMESTAMP(),
    PRIMARY KEY (id),
    FOREIGN KEY (customer_id) REFERENCES customers (id)
);

-- 6. cart_items
CREATE TABLE IF NOT EXISTS cart_items (
    id         STRING       NOT NULL,
    cart_id    STRING       NOT NULL,
    cat_id     STRING       NOT NULL,
    quantity   INT          NOT NULL CHECK (quantity > 0),
    created_at TIMESTAMP_NTZ NOT NULL DEFAULT CURRENT_TIMESTAMP(),
    PRIMARY KEY (id),
    FOREIGN KEY (cart_id) REFERENCES carts (id),
    FOREIGN KEY (cat_id)  REFERENCES cats (id)
);

-- 7. orders
CREATE TABLE IF NOT EXISTS orders (
    id           STRING       NOT NULL,
    customer_id  STRING       NOT NULL,
    status       STRING       NOT NULL DEFAULT 'PENDING',
    total_cents  NUMBER(12,2) NOT NULL,
    currency     STRING(3)    NOT NULL,
    ship_address STRING       NOT NULL,
    bill_address STRING       NOT NULL,
    notes        STRING,
    created_at   TIMESTAMP_NTZ NOT NULL DEFAULT CURRENT_TIMESTAMP(),
    updated_at   TIMESTAMP_NTZ NOT NULL DEFAULT CURRENT_TIMESTAMP(),
    placed_at    TIMESTAMP_NTZ,
    shipped_at   TIMESTAMP_NTZ,
    cancelled_at TIMESTAMP_NTZ,
    PRIMARY KEY (id),
    FOREIGN KEY (customer_id) REFERENCES customers (id)
);

-- 8. order_items
CREATE TABLE IF NOT EXISTS order_items (
    id          STRING       NOT NULL,
    order_id    STRING       NOT NULL,
    cat_id      STRING       NOT NULL,
    quantity    INT          NOT NULL CHECK (quantity > 0),
    price_cents NUMBER(12,2) NOT NULL,
    currency    STRING(3)    NOT NULL,
    created_at  TIMESTAMP_NTZ NOT NULL DEFAULT CURRENT_TIMESTAMP(),
    PRIMARY KEY (id),
    FOREIGN KEY (order_id) REFERENCES orders (id),
    FOREIGN KEY (cat_id)   REFERENCES cats (id)
);

-- 9. inventory
CREATE TABLE IF NOT EXISTS inventory (
    cat_id     STRING       NOT NULL,
    available  INT          NOT NULL DEFAULT 0,
    reserved   INT          NOT NULL DEFAULT 0,
    updated_at TIMESTAMP_NTZ NOT NULL DEFAULT CURRENT_TIMESTAMP(),
    PRIMARY KEY (cat_id),
    FOREIGN KEY (cat_id) REFERENCES cats (id)
);

-- 10. audit_log
CREATE TABLE IF NOT EXISTS audit_log (
    id         STRING       NOT NULL,
    actor      STRING       NOT NULL,
    action     STRING       NOT NULL,
    target     STRING       NOT NULL,
    created_at TIMESTAMP_NTZ NOT NULL DEFAULT CURRENT_TIMESTAMP(),
    PRIMARY KEY (id)
);
```

### Drop (clean slate)

```sql
USE DATABASE TEST;
DROP SCHEMA IF EXISTS PETSTORE CASCADE;
```

-- TransferStore
CREATE TABLE IF NOT EXISTS transfer_stores (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(255) NOT NULL,
    user_id CHAR(36),
    status VARCHAR(50) CHECK (status IN ('progress', 'done', 'cancel')),
    store_id CHAR(36)
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_transfer_stores_code ON transfer_stores(code);
CREATE INDEX IF NOT EXISTS idx_transfer_stores_user_store ON transfer_stores(user_id, store_id);

-- Store
CREATE TABLE IF NOT EXISTS stores (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL,
    phone VARCHAR(20),
    address TEXT,
    is_active BOOLEAN DEFAULT true
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_stores_slug_phone ON stores(slug, phone);

-- StockStore
CREATE TABLE IF NOT EXISTS stock_stores (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id CHAR(36),
    store_id CHAR(36),
    transfer_store_id CHAR(36),
    status VARCHAR(50) CHECK (status IN ('stock', 'sale', 'out', 'discrepancy')),
    date_out TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_stock_stores_product_store ON stock_stores(product_id, store_id);

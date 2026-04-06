-- EXTENSION
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- =========================
-- Document
-- =========================
CREATE TABLE IF NOT EXISTS product_documents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(255) NOT NULL,
    file_name VARCHAR(255),
    file_item INT,
    file_price INT,
    status VARCHAR(50) CHECK (status IN ('pending', 'progress', 'done', 'cancel')),
    type VARCHAR(50) CHECK (type IN ('bast', 'bulk', 'manual', 'sku')),
    header_barcode VARCHAR(255),
    header_name VARCHAR(255),
    header_item VARCHAR(255),
    header_price VARCHAR(255),
    user_id UUID,
    supplier VARCHAR(255),
    type_product VARCHAR(50) CHECK (type_product IN ('reguler', 'sticker', 'refurbish')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_product_documents_code 
ON product_documents(code);

-- =========================
-- Product Pendings
-- =========================
CREATE TABLE IF NOT EXISTS product_pendings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    document_id UUID,
    barcode VARCHAR(255),
    name TEXT,
    item INT,
    price DECIMAL(15,2),
    is_sku BOOLEAN DEFAULT false,
    status VARCHAR(50) CHECK (status IN ('good', 'damaged', 'abnormal', 'non')),
    note TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,

    CONSTRAINT fk_product_pendings_document
    FOREIGN KEY (document_id)
    REFERENCES product_documents(id)
    ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_product_pendings_document_id 
ON product_pendings(document_id);

-- =========================
-- Rack Displays
-- =========================
CREATE TABLE IF NOT EXISTS rack_displays (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_rack_displays_code 
ON rack_displays(code);

-- =========================
-- Rack Stagings
-- =========================
CREATE TABLE IF NOT EXISTS rack_stagings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    rack_display_id UUID,
    code VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    is_moved BOOLEAN DEFAULT false,

    CONSTRAINT fk_rack_stagings_display
    FOREIGN KEY (rack_display_id)
    REFERENCES rack_displays(id)
    ON DELETE SET NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_rack_stagings_code 
ON rack_stagings(code);

CREATE INDEX IF NOT EXISTS idx_rack_stagings_display_id 
ON rack_stagings(rack_display_id);

-- =========================
-- Product Masters
-- =========================
CREATE TABLE IF NOT EXISTS product_masters (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_pending_id UUID,
    document_id UUID,
    barcode VARCHAR(255),
    barcode_warehouse VARCHAR(255),
    name TEXT,
    name_warehouse TEXT,
    item INT,
    item_warehouse INT,
    price DECIMAL(15,2),
    price_warehouse DECIMAL(15,2),
    category_id UUID,
    sticker_id UUID,
    is_sku BOOLEAN DEFAULT false,
    location VARCHAR(50) CHECK (
        location IN (
            'staging_reguler', 'staging_sticker', 'display',
            'cargo', 'scrap', 'qcd', 'repair', 'staging_sku'
        )
    ),
    bundle VARCHAR(50) CHECK (bundle IN ('yes', 'no')) DEFAULT 'no',
    bundle_parent_id UUID,
    date_out TIMESTAMP,
    type_out VARCHAR(50) CHECK (
        type_out IN ('regular_sales', 'cargo', 'scrap', 'qcd', 'transfer')
    ),
    rack_staging_id UUID,
    rack_display_id UUID,
    bag_id UUID,
    user_id UUID,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,

    CONSTRAINT fk_product_masters_pending
    FOREIGN KEY (product_pending_id)
    REFERENCES product_pendings(id)
    ON DELETE SET NULL,

    CONSTRAINT fk_product_masters_document
    FOREIGN KEY (document_id)
    REFERENCES product_documents(id)
    ON DELETE CASCADE,

    CONSTRAINT fk_product_masters_rack_staging
    FOREIGN KEY (rack_staging_id)
    REFERENCES rack_stagings(id)
    ON DELETE SET NULL,

    CONSTRAINT fk_product_masters_rack_display
    FOREIGN KEY (rack_display_id)
    REFERENCES rack_displays(id)
    ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_product_masters_barcode 
ON product_masters(barcode);

CREATE INDEX IF NOT EXISTS idx_product_masters_lookup 
ON product_masters(category_id, sticker_id, bag_id);
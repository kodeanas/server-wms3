-- Document
CREATE TABLE IF NOT EXISTS product_documents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(255) NOT NULL,
    file_name VARCHAR(255),
    file_item INT,
    file_price INT,
    status VARCHAR(50) CHECK (status IN ('pending', 'progress', 'done', 'cancel')),
    type VARCHAR(50) CHECK (type IN ('baist', 'bulk', 'manual', 'sku')),
    sku VARCHAR(255),
    header_barcode VARCHAR(255),
    header_name VARCHAR(255),
    header_item VARCHAR(255),
    header_price VARCHAR(255),
    user_id CHAR(36),
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_product_documents_code ON product_documents(code);

-- ProductPending
CREATE TABLE IF NOT EXISTS product_pendings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    document_id CHAR(36),
    barcode VARCHAR(255),
    name TEXT,
    item INT,
    price DECIMAL(15,2),
    is_sku BOOLEAN DEFAULT false,
    status VARCHAR(50) CHECK (status IN ('discrepancy', 'good', 'damaged', 'abnormal', 'non')),
    note TEXT,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_product_pendings_document_id ON product_pendings(document_id);

-- RackDisplay
CREATE TABLE IF NOT EXISTS rack_displays (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_rack_displays_code ON rack_displays(code);

-- RackStaging
CREATE TABLE IF NOT EXISTS rack_stagings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    rack_display_id CHAR(36),
    code VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_moved BOOLEAN DEFAULT false
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_rack_stagings_code ON rack_stagings(code);
CREATE INDEX IF NOT EXISTS idx_rack_stagings_display_id ON rack_stagings(rack_display_id);

-- ProductMaster
CREATE TABLE IF NOT EXISTS product_masters (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    document_id CHAR(36),
    barcode VARCHAR(255),
    barcode_warehouse VARCHAR(255),
    name TEXT,
    name_warehouse TEXT,
    item INT,
    item_warehouse INT,
    price DECIMAL(15,2),
    price_warehouse DECIMAL(15,2),
    category_id CHAR(36),
    sticker_id CHAR(36),
    is_sku BOOLEAN DEFAULT false,
    location VARCHAR(50) CHECK (location IN ('rack', 'card', 'bundle')),
    bundle_parent_id CHAR(36),
    date_out TIMESTAMP,
    type_out VARCHAR(50) CHECK (type_out IN ('regular_sales', 'cargo', 'scrap', 'qcd', 'transfer')),
    rack_staging_id CHAR(36),
    rack_display_id CHAR(36),
    bag_id CHAR(36),
    user_id CHAR(36),
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_product_masters_barcode ON product_masters(barcode);
CREATE INDEX IF NOT EXISTS idx_product_masters_lookup ON product_masters(category_id, sticker_id, bag_id);

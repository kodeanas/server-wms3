-- Create Product Documents Table
CREATE TABLE IF NOT EXISTS product_documents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    code VARCHAR(255) NOT NULL,
    header_barcode VARCHAR(255),
    header_name VARCHAR(255),
    header_quantity VARCHAR(255),
    header_price VARCHAR(255),
    type_in VARCHAR(50) CHECK (type_in IN ('baist', 'bulking', 'manual', 'sku')),
    type_bulking VARCHAR(50) CHECK (type_bulking IN ('regular', 'sticker')),
    total_list INT,
    total_price DECIMAL(15,2),
    user_id CHAR(36),
    supplier_id CHAR(36),
    status VARCHAR(50) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_product_documents_id ON product_documents(id);
CREATE INDEX IF NOT EXISTS idx_product_documents_code ON product_documents(code);
CREATE INDEX IF NOT EXISTS idx_product_documents_user_supplier ON product_documents(user_id, supplier_id);

-- Create Product Masters Table
CREATE TABLE IF NOT EXISTS product_masters (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_document_id CHAR(36),
    barcode VARCHAR(255) NOT NULL,
    barcode_warehouse VARCHAR(255),
    barcode_source VARCHAR(50),
    qty INT,
    qty_source INT,
    price DECIMAL(15,2),
    price_source DECIMAL(15,2),
    price_warehouse DECIMAL(15,2),
    category_id CHAR(36),
    sticker_id CHAR(36),
    reguleher VARCHAR(50) CHECK (reguleher IN ('regular', 'wholesale')),
    out_quantity INT DEFAULT 0,
    is_reidentify BOOLEAN DEFAULT false,
    bundle_id CHAR(36),
    sku_id CHAR(36),
    bag_id CHAR(36),
    user_id CHAR(36),
    status VARCHAR(50) CHECK (status IN ('lulos', 'damaged', 'abnormal', 'non')),
    status_source VARCHAR(50) CHECK (status_source IN ('lulos', 'damaged', 'abnormal', 'non')),
    store_id INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_product_masters_id ON product_masters(id);
CREATE INDEX IF NOT EXISTS idx_product_masters_barcode ON product_masters(barcode);
CREATE INDEX IF NOT EXISTS idx_product_masters_category_sticker_bundle_sku_bag ON product_masters(category_id, sticker_id, bundle_id, sku_id, bag_id);

-- Create Product Logs Table
CREATE TABLE IF NOT EXISTS product_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_master_id CHAR(36),
    data_name VARCHAR(255),
    prev_data TEXT,
    new_data TEXT,
    user_id CHAR(36),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_product_logs_id ON product_logs(id);
CREATE INDEX IF NOT EXISTS idx_product_logs_product_master_user ON product_logs(product_master_id, user_id);

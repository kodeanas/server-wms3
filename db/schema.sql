-- WMS Database Schema
-- Warehouse Management System Database Structure
-- Generated for PostgreSQL

-- ============================================
-- 1. USER MANAGEMENT
-- ============================================

-- Users Table
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    phone VARCHAR(20) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_phone ON users(phone);

-- Taxes Table
CREATE TABLE IF NOT EXISTS taxes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tax INT NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_taxes_id ON taxes(id);

-- ============================================
-- 2. PRODUCT MANAGEMENT
-- ============================================

-- Categories Table
CREATE TABLE IF NOT EXISTS categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL UNIQUE,
    discount INT DEFAULT 0,
    min_price NUMERIC(19, 2),
    max_price NUMERIC(19, 2),
    status VARCHAR(50) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_categories_id ON categories(id);
CREATE INDEX IF NOT EXISTS idx_categories_slug ON categories(slug);

-- Stickers Table
CREATE TABLE IF NOT EXISTS stickers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code_hex VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL UNIQUE,
    type VARCHAR(50),
    fixed_price INT,
    min_price NUMERIC(19, 2),
    max_price NUMERIC(19, 2),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_stickers_id ON stickers(id);
CREATE INDEX IF NOT EXISTS idx_stickers_slug ON stickers(slug);

-- ============================================
-- 3. BUYER & CLASS MANAGEMENT
-- ============================================

-- Classes Table
CREATE TABLE IF NOT EXISTS classes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    min_order INT NOT NULL,
    disc INT,
    min_transaction_value NUMERIC(19, 2),
    week INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_classes_id ON classes(id);
CREATE INDEX IF NOT EXISTS idx_classes_name_min_order ON classes(name, min_order);

-- Buyers Table
CREATE TABLE IF NOT EXISTS buyers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255),
    phone VARCHAR(20) NOT NULL,
    class_id CHAR(36) REFERENCES classes(id) ON DELETE SET NULL,
    address TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_buyers_id ON buyers(id);
CREATE INDEX IF NOT EXISTS idx_buyers_email_phone ON buyers(email, phone);
CREATE INDEX IF NOT EXISTS idx_buyers_class_id ON buyers(class_id);

-- ============================================
-- 4. STORE MANAGEMENT
-- ============================================

-- Stores Table
CREATE TABLE IF NOT EXISTS stores (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    phone VARCHAR(20) NOT NULL,
    email VARCHAR(255),
    address TEXT,
    user_id CHAR(36) UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    is_cashier BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_stores_id ON stores(id);
CREATE INDEX IF NOT EXISTS idx_stores_phone_email ON stores(phone, email);
CREATE INDEX IF NOT EXISTS idx_stores_user_id ON stores(user_id);

-- Store Crews Table
CREATE TABLE IF NOT EXISTS store_crews (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    phone VARCHAR(20) NOT NULL,
    email VARCHAR(255),
    address TEXT,
    password CHAR(36),
    store_id CHAR(36) NOT NULL REFERENCES stores(id) ON DELETE CASCADE,
    is_cashier BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_store_crews_id ON store_crews(id);
CREATE INDEX IF NOT EXISTS idx_store_crews_phone_email ON store_crews(phone, email);
CREATE INDEX IF NOT EXISTS idx_store_crews_store_id ON store_crews(store_id);

-- ============================================
-- 5. WAREHOUSE INFRASTRUCTURE
-- ============================================

-- Racks Table
CREATE TABLE IF NOT EXISTS racks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    type VARCHAR(50),
    name VARCHAR(255) NOT NULL,
    rack_id VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_racks_id ON racks(id);
CREATE INDEX IF NOT EXISTS idx_racks_rack_id ON racks(rack_id);

-- ============================================
-- 6. PRODUCT CATALOG
-- ============================================

-- Product Documents Table
CREATE TABLE IF NOT EXISTS product_documents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    code VARCHAR(255) NOT NULL UNIQUE,
    header_barcode VARCHAR(255),
    header_name VARCHAR(255),
    header_quantity VARCHAR(255),
    header_price VARCHAR(255),
    type_in VARCHAR(50), -- bast, bulking, manual, auto, shu
    type_bulking VARCHAR(50), -- regular, sticker, sticker_ordered, damaged, qcd
    user_id INT,
    supplier_id CHAR(36),
    total_list INT,
    total_price NUMERIC(19, 2),
    status VARCHAR(50) DEFAULT 'done',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_product_documents_id ON product_documents(id);
CREATE INDEX IF NOT EXISTS idx_product_documents_user_supplier ON product_documents(user_id, supplier_id);
CREATE INDEX IF NOT EXISTS idx_product_documents_code ON product_documents(code);

-- Product Masters Table
CREATE TABLE IF NOT EXISTS product_masters (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_document_id CHAR(36) REFERENCES product_documents(id),
    barcode VARCHAR(255) NOT NULL,
    barcode_source VARCHAR(255),
    barcode_warehouse VARCHAR(255),
    qty INT NOT NULL,
    qty_source INT,
    price NUMERIC(19, 2) NOT NULL,
    price_source NUMERIC(19, 2),
    price_warehouse NUMERIC(19, 2),
    category_id CHAR(36) REFERENCES categories(id),
    sticker_id CHAR(36) REFERENCES stickers(id),
    bundle_id CHAR(36),
    sku_id CHAR(36),
    bag_id CHAR(36),
    user_id CHAR(36) NOT NULL REFERENCES users(id),
    out VARCHAR(50), -- regular, wholesale
    is_reidentify BOOLEAN DEFAULT FALSE,
    date_time TIMESTAMP,
    rack_id CHAR(36) REFERENCES racks(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_product_masters_id ON product_masters(id);
CREATE INDEX IF NOT EXISTS idx_product_masters_barcode_warehouse ON product_masters(barcode_warehouse);
CREATE INDEX IF NOT EXISTS idx_product_masters_category_sticker ON product_masters(category_id, sticker_id, bundle_id, sku_id, bag_id);
CREATE INDEX IF NOT EXISTS idx_product_masters_user_id ON product_masters(user_id);

-- Product Logs Table
CREATE TABLE IF NOT EXISTS product_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_master_id CHAR(36) NOT NULL REFERENCES product_masters(id) ON DELETE CASCADE,
    data_name VARCHAR(255),
    prev_data TEXT,
    new_data TEXT,
    date_time TIMESTAMP,
    user_id CHAR(36) NOT NULL REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_product_logs_id ON product_logs(id);
CREATE INDEX IF NOT EXISTS idx_product_logs_product_id ON product_logs(product_master_id);

-- ============================================
-- 7. CARGO MANAGEMENT
-- ============================================

-- Cargos Table
CREATE TABLE IF NOT EXISTS cargos (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(255) NOT NULL UNIQUE,
    date_time TIMESTAMP NOT NULL,
    user_id CHAR(36) NOT NULL REFERENCES users(id),
    total_quantity INT,
    total_price NUMERIC(19, 2),
    status VARCHAR(50) DEFAULT 'on_progress', -- on_progress, done
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_cargos_id ON cargos(id);
CREATE INDEX IF NOT EXISTS idx_cargos_code ON cargos(code);
CREATE INDEX IF NOT EXISTS idx_cargos_user_id ON cargos(user_id);

-- Bags Table
CREATE TABLE IF NOT EXISTS bags (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    cargo_id CHAR(36) REFERENCES cargos(id),
    code VARCHAR(255) NOT NULL UNIQUE,
    type VARCHAR(50), -- regular, sku, sticker, bi, scrap, qcd
    date_time TIMESTAMP NOT NULL,
    user_id CHAR(36) NOT NULL REFERENCES users(id),
    total_quantity INT,
    total_price NUMERIC(19, 2),
    status VARCHAR(50) DEFAULT 'done', -- done, on_progress
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_bags_id ON bags(id);
CREATE INDEX IF NOT EXISTS idx_bags_code ON bags(code);
CREATE INDEX IF NOT EXISTS idx_bags_user_id ON bags(user_id);

-- ============================================
-- 8. ORDER MANAGEMENT
-- ============================================

-- Orders Table
CREATE TABLE IF NOT EXISTS orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(255) NOT NULL UNIQUE,
    user_id CHAR(36) NOT NULL REFERENCES users(id),
    type VARCHAR(50), -- regular, wcargo, wjed, wcargo_w
    buyer_id CHAR(36) NOT NULL REFERENCES buyers(id),
    date_time TIMESTAMP NOT NULL,
    total_item INT,
    total_price_product NUMERIC(19, 2),
    total_ppn NUMERIC(19, 2),
    carton_box_price NUMERIC(19, 2),
    quantity_carton_box INT,
    voucher NUMERIC(19, 2),
    fixed_discount INT,
    class_discount INT,
    class_discount_amount NUMERIC(19, 2),
    grand_total_before NUMERIC(19, 2),
    grand_total_after NUMERIC(19, 2),
    status VARCHAR(50) DEFAULT 'on_progress', -- on_progress, finish
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_orders_id ON orders(id);
CREATE INDEX IF NOT EXISTS idx_orders_user_buyer_class ON orders(user_id, buyer_id);
CREATE INDEX IF NOT EXISTS idx_orders_code ON orders(code);

-- Order Items Table
CREATE TABLE IF NOT EXISTS order_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_master_id CHAR(36) NOT NULL REFERENCES product_masters(id),
    order_id CHAR(36) NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    price NUMERIC(19, 2),
    price_warehouse NUMERIC(19, 2),
    price_cut NUMERIC(19, 2),
    price_final NUMERIC(19, 2),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_order_items_id ON order_items(id);
CREATE INDEX IF NOT EXISTS idx_order_items_product_order ON order_items(product_master_id, order_id);

-- Order Cargos Table
CREATE TABLE IF NOT EXISTS order_cargos (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    cargo_id CHAR(36) NOT NULL REFERENCES cargos(id),
    order_id CHAR(36) NOT NULL REFERENCES orders(id),
    price NUMERIC(19, 2),
    price_warehouse NUMERIC(19, 2),
    price_cut NUMERIC(19, 2),
    price_final NUMERIC(19, 2),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_order_cargos_id ON order_cargos(id);
CREATE INDEX IF NOT EXISTS idx_order_cargos_product_order ON order_cargos(cargo_id, order_id);

-- ============================================
-- 9. STORE TRANSFER
-- ============================================

-- Store Transfers Table
CREATE TABLE IF NOT EXISTS store_transfers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    store_id CHAR(36) NOT NULL REFERENCES stores(id),
    date_time TIMESTAMP NOT NULL,
    total_item INT,
    total_price NUMERIC(19, 2),
    total_price_warehouse NUMERIC(19, 2),
    status VARCHAR(50) DEFAULT 'on_progress', -- on_progress, done
    user_id CHAR(36) REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_store_transfers_id ON store_transfers(id);
CREATE INDEX IF NOT EXISTS idx_store_transfers_store_bag ON store_transfers(store_id);

-- Store Transfer Bags Table
CREATE TABLE IF NOT EXISTS store_transfer_bags (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    store_transfer_id CHAR(36) NOT NULL REFERENCES store_transfers(id) ON DELETE CASCADE,
    bag_id CHAR(36) REFERENCES bags(id),
    quantity INT,
    total_price NUMERIC(19, 2),
    total_cogs NUMERIC(19, 2),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_store_transfer_bags_id ON store_transfer_bags(id);
CREATE INDEX IF NOT EXISTS idx_store_transfer_bags_store_bag ON store_transfer_bags(store_transfer_id, bag_id);

-- ============================================
-- 10. USER CLASS LOG
-- ============================================

-- User Class Logs Table
CREATE TABLE IF NOT EXISTS user_class_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    buyer_id CHAR(36) NOT NULL REFERENCES buyers(id),
    prev_class_id CHAR(36) REFERENCES classes(id),
    new_class_id CHAR(36) NOT NULL REFERENCES classes(id),
    order_id CHAR(36) NOT NULL REFERENCES orders(id),
    change_type VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_user_class_logs_id ON user_class_logs(id);
CREATE INDEX IF NOT EXISTS idx_user_class_logs_composite ON user_class_logs(prev_class_id, new_class_id, order_id, buyer_id);

-- ============================================
-- 11. SLOW MOVING INVENTORY
-- ============================================

-- Slow Movings Table
CREATE TABLE IF NOT EXISTS slow_movings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    date TIMESTAMP NOT NULL,
    total_item INT,
    total_price NUMERIC(19, 2),
    total_price_warehouse NUMERIC(19, 2),
    is_damaged BOOLEAN DEFAULT FALSE,
    user_id CHAR(36) NOT NULL REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_slow_movings_id ON slow_movings(id);
CREATE INDEX IF NOT EXISTS idx_slow_movings_product_slow ON slow_movings(user_id);

-- Slow Moving Items Table
CREATE TABLE IF NOT EXISTS slow_moving_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    slow_moving_id CHAR(36) NOT NULL REFERENCES slow_movings(id) ON DELETE CASCADE,
    product_master_id CHAR(36) NOT NULL REFERENCES product_masters(id),
    is_damaged BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_slow_moving_items_id ON slow_moving_items(id);
CREATE INDEX IF NOT EXISTS idx_slow_moving_items_product_id ON slow_moving_items(product_master_id, slow_moving_id);

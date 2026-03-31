-- Create Orders Table
CREATE TABLE IF NOT EXISTS orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(255) NOT NULL,
    user_id CHAR(36),
    type VARCHAR(50) CHECK (type IN ('regular', 'wcargo', 'wjed', 'wacop')),
    buyer_id CHAR(36),
    datetime TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    total_item INT,
    total_price_product DECIMAL(15,2),
    total_ppn DECIMAL(15,2),
    carton_box_price DECIMAL(15,2),
    quantity_carton_box INT,
    voucher DECIMAL(15,2),
    fixed_discount INT,
    class_discount INT,
    class_discount_amount DECIMAL(15,2),
    class_id CHAR(36),
    grand_total_before DECIMAL(15,2),
    grand_total_after DECIMAL(15,2),
    status VARCHAR(50) CHECK (status IN ('on_progress', 'finish')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_orders_id ON orders(id);
CREATE INDEX IF NOT EXISTS idx_orders_code ON orders(code);
CREATE INDEX IF NOT EXISTS idx_orders_user_buyer_class ON orders(user_id, buyer_id, class_id);

-- Create Order Items Table
CREATE TABLE IF NOT EXISTS order_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_master_id CHAR(36),
    order_id CHAR(36),
    price DECIMAL(15,2),
    price_warehouse DECIMAL(15,2),
    price_cut DECIMAL(15,2),
    prce_final DECIMAL(15,2),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_order_items_id ON order_items(id);
CREATE INDEX IF NOT EXISTS idx_order_items_product_order ON order_items(product_master_id, order_id);

-- Create Order Cargos Table
CREATE TABLE IF NOT EXISTS order_cargos (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    cargo_id CHAR(36),
    order_id CHAR(36),
    price DECIMAL(15,2),
    price_warehouse DECIMAL(15,2),
    price_cut DECIMAL(15,2),
    prce_final DECIMAL(15,2),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_order_cargos_id ON order_cargos(id);
CREATE INDEX IF NOT EXISTS idx_order_cargos_cargo_order ON order_cargos(cargo_id, order_id);

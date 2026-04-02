-- Order
CREATE TABLE IF NOT EXISTS orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(255) NOT NULL,
    type VARCHAR(50) CHECK (type IN ('regular', 'cargo', 'qcd', 'scrap')),
    buyer_id CHAR(36),
    rack_id CHAR(36),
    user_id CHAR(36),
    status VARCHAR(50) CHECK (status IN ('progress', 'pending', 'done', 'cancel')),
    is_tax BOOLEAN DEFAULT false,
    tax_value DECIMAL(15,2),
    tax DECIMAL(15,2),
    total_price DECIMAL(15,2),
    total_box INT,
    price_box DECIMAL(15,2),
    grand_total DECIMAL(15,2),
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_orders_code ON orders(code);
CREATE INDEX IF NOT EXISTS idx_orders_buyer_user ON orders(buyer_id, user_id);

-- ProductOrder
CREATE TABLE IF NOT EXISTS product_orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id CHAR(36),
    product_id CHAR(36),
    name VARCHAR(255),
    price DECIMAL(15,2),
    price_warehouse DECIMAL(15,2),
    discount DECIMAL(15,2)
);

CREATE INDEX IF NOT EXISTS idx_product_orders_order_product ON product_orders(order_id, product_id);

-- CargoOrder
CREATE TABLE IF NOT EXISTS cargo_orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id CHAR(36),
    cargo_id CHAR(36),
    name VARCHAR(255),
    price DECIMAL(15,2),
    price_warehouse DECIMAL(15,2)
);

CREATE INDEX IF NOT EXISTS idx_cargo_orders_order_cargo ON cargo_orders(order_id, cargo_id);

-- BagOrder
CREATE TABLE IF NOT EXISTS bag_orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id CHAR(36),
    cargo_id CHAR(36),
    name VARCHAR(255),
    price DECIMAL(15,2),
    price_warehouse DECIMAL(15,2)
);

CREATE INDEX IF NOT EXISTS idx_bag_orders_order_cargo ON bag_orders(order_id, cargo_id);

-- DiscountOrder
CREATE TABLE IF NOT EXISTS discount_orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id CHAR(36),
    type VARCHAR(50) CHECK (type IN ('voucher', 'rank', 'additional')),
    name VARCHAR(255),
    is_nominal BOOLEAN DEFAULT true,
    value_nominal DECIMAL(15,2),
    value_percentage INT
);

CREATE INDEX IF NOT EXISTS idx_discount_orders_order_id ON discount_orders(order_id);

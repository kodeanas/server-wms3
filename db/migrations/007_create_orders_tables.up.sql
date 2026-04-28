-- =========================
-- Orders
-- =========================
CREATE TABLE IF NOT EXISTS orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(255) NOT NULL,

    type VARCHAR(50) CHECK (type IN ('regular', 'cargo', 'qcd', 'scrap')),

    buyer_id UUID,
    class_id UUID,
    user_id UUID,

    status VARCHAR(50) CHECK (status IN ('progress', 'pending', 'done', 'cancel')),

    is_tax BOOLEAN DEFAULT false,
    tax DECIMAL(15,2),
    tax_value DECIMAL(15,2),

    total_price DECIMAL(15,2),
    total_box INT,
    price_box DECIMAL(15,2),
    grand_total DECIMAL(15,2),

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,

    -- 🔗 RELATIONS
    CONSTRAINT fk_orders_buyer
        FOREIGN KEY (buyer_id)
        REFERENCES buyers(id)
        ON DELETE SET NULL,

    CONSTRAINT fk_orders_class
        FOREIGN KEY (class_id)
        REFERENCES classes(id)
        ON DELETE SET NULL,

    CONSTRAINT fk_orders_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE SET NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_orders_code ON orders(code);
CREATE INDEX IF NOT EXISTS idx_orders_lookup 
ON orders(buyer_id, class_id, user_id);

-- =========================
-- ProductOrder
-- =========================
CREATE TABLE IF NOT EXISTS product_orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    order_id UUID,
    product_id UUID,

    name VARCHAR(255),
    price DECIMAL(15,2),
    price_warehouse DECIMAL(15,2),
    discount DECIMAL(15,2),

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,

    CONSTRAINT fk_product_orders_order
        FOREIGN KEY (order_id)
        REFERENCES orders(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_product_orders_product
        FOREIGN KEY (product_id)
        REFERENCES product_masters(id)
        ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_product_orders_lookup
ON product_orders(order_id, product_id);

-- =========================
-- CargoOrder
-- =========================
CREATE TABLE IF NOT EXISTS cargo_orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    order_id UUID,
    cargo_id UUID,

    name VARCHAR(255),
    price DECIMAL(15,2),
    price_warehouse DECIMAL(15,2),

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,

    CONSTRAINT fk_cargo_orders_order
        FOREIGN KEY (order_id)
        REFERENCES orders(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_cargo_orders_cargo
        FOREIGN KEY (cargo_id)
        REFERENCES cargos(id)
        ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_cargo_orders_lookup
ON cargo_orders(order_id, cargo_id);

-- =========================
-- BagOrder
-- =========================
CREATE TABLE IF NOT EXISTS bag_orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    order_id UUID,
    cargo_id UUID,

    name VARCHAR(255),
    price DECIMAL(15,2),
    price_warehouse DECIMAL(15,2),

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,

    CONSTRAINT fk_bag_orders_order
        FOREIGN KEY (order_id)
        REFERENCES orders(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_bag_orders_cargo
        FOREIGN KEY (cargo_id)
        REFERENCES cargos(id)
        ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_bag_orders_lookup
ON bag_orders(order_id, cargo_id);

-- =========================
-- DiscounOrder
-- =========================
CREATE TABLE IF NOT EXISTS discount_orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    order_id UUID,

    type VARCHAR(50) CHECK (type IN ('voucher', 'class', 'additional')),
    name VARCHAR(255),

    is_nominal BOOLEAN DEFAULT true,
    value DECIMAL(15,2),
    usage INT,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,

    CONSTRAINT fk_discount_orders_order
        FOREIGN KEY (order_id)
        REFERENCES orders(id)
        ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_discount_orders_order_id 
ON discount_orders(order_id);

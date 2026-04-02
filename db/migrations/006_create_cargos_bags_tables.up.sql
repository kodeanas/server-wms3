-- ProductRepair
CREATE TABLE IF NOT EXISTS product_repairs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id CHAR(36),
    status VARCHAR(50) CHECK (status IN ('progress', 'done', 'out')),
    date_out TIMESTAMP,
    user_id CHAR(36),
    price_update DECIMAL(15,2),
    item_update INT,
    category_id CHAR(36),
    sticker_id CHAR(36),
    remark_origin TEXT,
    remark_after TEXT
);

CREATE INDEX IF NOT EXISTS idx_product_repairs_lookup ON product_repairs(product_id, user_id, category_id, sticker_id);

-- Cargo
CREATE TABLE IF NOT EXISTS cargos (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(255) NOT NULL,
    status VARCHAR(50) CHECK (status IN ('open', 'lock', 'cancel')),
    is_sale BOOLEAN DEFAULT false,
    is_online BOOLEAN DEFAULT false,
    user_id CHAR(36)
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_cargos_code ON cargos(code);
CREATE INDEX IF NOT EXISTS idx_cargos_user_id ON cargos(user_id);

-- Bag
CREATE TABLE IF NOT EXISTS bags (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(255) NOT NULL,
    type VARCHAR(50) CHECK (type IN ('sticker', 'regular', 'qcd', 'scrap', 'bkl')),
    user_id CHAR(36),
    lock_status VARCHAR(50) CHECK (lock_status IN ('open', 'lock', 'out')),
    date_out DATE,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    cargo_id CHAR(36),
    transfer_store_id CHAR(36)
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_bags_code ON bags(code);
CREATE INDEX IF NOT EXISTS idx_bags_user_id ON bags(user_id);

-- ItemBag
CREATE TABLE IF NOT EXISTS item_bags (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id CHAR(36),
    bag_id CHAR(36)
);

CREATE INDEX IF NOT EXISTS idx_item_bags_product_bag ON item_bags(product_id, bag_id);

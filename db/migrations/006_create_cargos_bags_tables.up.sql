-- =========================
-- Product Repairs
-- =========================

CREATE TABLE IF NOT EXISTS product_repairs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id UUID,  -- Ubah dari CHAR(36) ke UUID
    document_id UUID,
    status VARCHAR(50) CHECK (status IN ('progress', 'done', 'out')),
    date_out TIMESTAMP,
    user_id UUID,  -- Ubah dari CHAR(36) ke UUID
    price_before DECIMAL(15,2),
    price_update DECIMAL(15,2),
    item_before INT,
    item_update INT,
    category_id UUID,  -- Ubah dari CHAR(36) ke UUID
    sticker_id UUID,  -- Ubah dari CHAR(36) ke UUID
    remark_origin TEXT,
    remark_after TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,

    CONSTRAINT fk_product_repairs_document
    FOREIGN KEY (document_id)
    REFERENCES product_documents(id)
    ON DELETE CASCADE,

    CONSTRAINT fk_product_repairs_user
    FOREIGN KEY (user_id)
    REFERENCES users(id)
    ON DELETE SET NULL,

    CONSTRAINT fk_product_repairs_category
    FOREIGN KEY (category_id)
    REFERENCES categories(id)
    ON DELETE SET NULL,

    CONSTRAINT fk_product_repairs_sticker
    FOREIGN KEY (sticker_id)
    REFERENCES stickers(id)
    ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_product_repairs_lookup ON product_repairs(product_id, user_id, category_id, sticker_id);

-- =========================
-- Cargo
-- =========================
CREATE TABLE IF NOT EXISTS cargos (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(255) NOT NULL,
    status VARCHAR(50) CHECK (status IN ('open', 'lock', 'cancel')),
    is_sale BOOLEAN DEFAULT false,
    is_online BOOLEAN DEFAULT false,
    user_id UUID,  -- Ubah dari CHAR(36) ke UUID
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_cargos_code ON cargos(code);
CREATE INDEX IF NOT EXISTS idx_cargos_user_id ON cargos(user_id);

-- =========================
-- Bag
-- =========================
CREATE TABLE IF NOT EXISTS bags (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(255) NOT NULL,
    type VARCHAR(50) CHECK (type IN ('sticker', 'regular', 'qcd', 'scrap', 'bkl')),
    user_id UUID,  -- Ubah dari CHAR(36) ke UUID
    lock_status VARCHAR(50) CHECK (lock_status IN ('open', 'lock', 'out')),
    date_out DATE,
    cargo_id UUID,  -- Ubah dari CHAR(36) ke UUID
    transfer_store_id UUID,  -- Ubah dari CHAR(36) ke UUID
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_bags_code ON bags(code);
CREATE INDEX IF NOT EXISTS idx_bags_user_id ON bags(user_id);

-- =========================
-- Item Bag
-- =========================
CREATE TABLE IF NOT EXISTS item_bags (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id UUID,  -- Ubah dari CHAR(36) ke UUID
    bag_id UUID,  -- Ubah dari CHAR(36) ke UUID
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_item_bags_product_bag ON item_bags(product_id, bag_id);
-- =========================
-- Product Repairs
-- =========================

CREATE TABLE IF NOT EXISTS product_repairs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id UUID,  -- Ubah dari CHAR(36) ke UUID
    document_id UUID,
    user_id UUID,  -- Ubah dari CHAR(36) ke UUID
    category_id UUID,  -- Ubah dari CHAR(36) ke UUID
    sticker_id UUID,  -- Ubah dari CHAR(36) ke UUID

    status VARCHAR(50) CHECK (status IN ('progress', 'done', 'out')),
    date_out TIMESTAMP,

    price_before DECIMAL(15,2),
    price_update DECIMAL(15,2),
    item_before INT,
    item_update INT,
    
    remark_origin TEXT,
    remark_after TEXT,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,

    CONSTRAINT fk_product_repairs_product
        FOREIGN KEY (product_id)
        REFERENCES product_masters(id)
        ON DELETE RESTRICT,

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
    user_id UUID,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,

    CONSTRAINT fk_cargos_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE SET NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_cargos_code ON cargos(code);
CREATE INDEX IF NOT EXISTS idx_cargos_user_id ON cargos(user_id);

-- =========================
-- Bag
-- =========================
CREATE TABLE IF NOT EXISTS bags (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(255) NOT NULL,
    type VARCHAR(50) CHECK (type IN ('sticker', 'reguler', 'qcd', 'scrap', 'bkl')),
    user_id UUID,
    is_moved BOOLEAN DEFAULT false,
    date_out DATE,
    cargo_id UUID,
    transfer_store_id UUID,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,

    CONSTRAINT fk_bags_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE SET NULL,

    CONSTRAINT fk_bags_cargo
        FOREIGN KEY (cargo_id)
        REFERENCES cargos(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_bags_transfer_store
        FOREIGN KEY (transfer_store_id)
        REFERENCES stores(id)
        ON DELETE SET NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_bags_code ON bags(code);
CREATE INDEX IF NOT EXISTS idx_bags_user_id ON bags(user_id);
CREATE INDEX IF NOT EXISTS idx_bags_cargo_id ON bags(cargo_id);

-- =========================
-- Item Bag
-- =========================
CREATE TABLE IF NOT EXISTS item_bags (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id UUID,
    bag_id UUID,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,

    CONSTRAINT fk_item_bags_product
        FOREIGN KEY (product_id)
        REFERENCES product_masters(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_item_bags_bag
        FOREIGN KEY (bag_id)
        REFERENCES bags(id)
        ON DELETE CASCADE
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_item_bags_product_bag 
ON item_bags(product_id, bag_id);
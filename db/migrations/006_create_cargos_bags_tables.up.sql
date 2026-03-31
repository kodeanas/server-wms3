-- Create Cargos Table
CREATE TABLE IF NOT EXISTS cargos (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(255) NOT NULL,
    datetime TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    user_id CHAR(36),
    total_quantity INT,
    total_price DECIMAL(15,2),
    status VARCHAR(50) CHECK (status IN ('on_progress', 'done')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_cargos_id ON cargos(id);
CREATE INDEX IF NOT EXISTS idx_cargos_code ON cargos(code);
CREATE INDEX IF NOT EXISTS idx_cargos_user_id ON cargos(user_id);

-- Create Bags Table
CREATE TABLE IF NOT EXISTS bags (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    cargo_id CHAR(36),
    code VARCHAR(255) NOT NULL,
    type VARCHAR(50) CHECK (type IN ('regular', 'sku', 'sticker', 'bulk', 'scrap', 'qcd')),
    datetime TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    user_id CHAR(36),
    total_quantity INT,
    total_price DECIMAL(15,2),
    status VARCHAR(50) CHECK (status IN ('done', 'on_progress')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_bags_id ON bags(id);
CREATE INDEX IF NOT EXISTS idx_bags_code ON bags(code);
CREATE INDEX IF NOT EXISTS idx_bags_cargo_user ON bags(cargo_id, user_id);

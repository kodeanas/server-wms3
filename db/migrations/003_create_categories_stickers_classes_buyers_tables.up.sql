-- Category
CREATE TABLE IF NOT EXISTS categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL,
    discount INT DEFAULT 0,
    min_price DECIMAL(15,2),
    max_price DECIMAL(15,2),
    status VARCHAR(50) CHECK (status IN ('active', 'inactive')) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_categories_slug ON categories(slug);

-- Sticker
CREATE TABLE IF NOT EXISTS stickers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code_hex VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL,
    type VARCHAR(50) CHECK (type IN ('big', 'small', 'tiny')),
    fixed_price INT,
    min_price DECIMAL(15,2),
    max_price DECIMAL(15,2),
    status VARCHAR(50) CHECK (status IN ('active', 'inactive')) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_stickers_slug ON stickers(slug);

-- Class
CREATE TABLE IF NOT EXISTS classes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    min_order INT NOT NULL,
    disc INT DEFAULT 0,
    min_transaction_value DECIMAL(15,2),
    week INT,
    iteration INT,
    status VARCHAR(50) CHECK (status IN ('active', 'inactive')) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_classes_name_min_order ON classes(name, min_order);

-- Buyer
CREATE TABLE IF NOT EXISTS buyers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255),
    email VARCHAR(255),
    phone VARCHAR(20),
    class_id CHAR(36),
    address VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_buyers_class_id ON buyers(class_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_buyers_email_phone ON buyers(email, phone);

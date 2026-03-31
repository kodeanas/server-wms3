-- Create Categories Table
CREATE TABLE IF NOT EXISTS categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) UNIQUE NOT NULL,
    discount INT DEFAULT 0,
    min_price DECIMAL(15,2),
    max_price DECIMAL(15,2),
    status VARCHAR(50) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_categories_id ON categories(id);
CREATE INDEX IF NOT EXISTS idx_categories_slug ON categories(slug);

-- Create Stickers Table
CREATE TABLE IF NOT EXISTS stickers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code_hex VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) UNIQUE NOT NULL,
    type VARCHAR(50) CHECK (type IN ('big', 'small')),
    fixed_price INT,
    min_price DECIMAL(15,2),
    max_price DECIMAL(15,2),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_stickers_id ON stickers(id);
CREATE INDEX IF NOT EXISTS idx_stickers_slug ON stickers(slug);

-- Create Classes Table
CREATE TABLE IF NOT EXISTS classes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    min_order INT NOT NULL,
    disc INT DEFAULT 0,
    min_transaction_value DECIMAL(15,2),
    week INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_classes_id ON classes(id);
CREATE INDEX IF NOT EXISTS idx_classes_name_min_order ON classes(name, min_order);

-- Create Buyers Table
CREATE TABLE IF NOT EXISTS buyers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    nama VARCHAR(255),
    email VARCHAR(255),
    phone VARCHAR(20),
    class_id CHAR(36),
    address TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_buyers_id ON buyers(id);
CREATE INDEX IF NOT EXISTS idx_buyers_email_phone ON buyers(email, phone);
CREATE INDEX IF NOT EXISTS idx_buyers_class_id ON buyers(class_id);

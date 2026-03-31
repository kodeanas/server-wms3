-- Create Stores Table
CREATE TABLE IF NOT EXISTS stores (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    phone VARCHAR(20),
    email VARCHAR(255),
    address TEXT,
    user_id CHAR(36) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_stores_id ON stores(id);
CREATE INDEX IF NOT EXISTS idx_stores_phone_email ON stores(phone, email);
CREATE INDEX IF NOT EXISTS idx_stores_user_id ON stores(user_id);

-- Create Store Crews Table
CREATE TABLE IF NOT EXISTS store_crews (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    phone VARCHAR(20),
    email VARCHAR(255),
    address TEXT,
    store_id CHAR(36) NOT NULL,
    is_cashier BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_store_crews_id ON store_crews(id);
CREATE INDEX IF NOT EXISTS idx_store_crews_phone_email ON store_crews(phone, email);
CREATE INDEX IF NOT EXISTS idx_store_crews_store_id ON store_crews(store_id);

-- Create Racks Table
CREATE TABLE IF NOT EXISTS racks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    type VARCHAR(50) CHECK (type IN ('display', 'staging')),
    name VARCHAR(255) NOT NULL,
    rack_id VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_racks_id ON racks(id);
CREATE INDEX IF NOT EXISTS idx_racks_rack_id ON racks(rack_id);

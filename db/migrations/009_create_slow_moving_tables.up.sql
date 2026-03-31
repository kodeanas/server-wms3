-- Create Slow Moving Table
CREATE TABLE IF NOT EXISTS slow_movings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    date TIMESTAMP,
    total_item INT,
    total_price DECIMAL(15,2),
    total_price_warehouse DECIMAL(15,2),
    is_damaged BOOLEAN DEFAULT false,
    user_id CHAR(36),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_slow_movings_id ON slow_movings(id);
CREATE INDEX IF NOT EXISTS idx_slow_movings_user_id ON slow_movings(user_id);

-- Create Slow Moving Items Table
CREATE TABLE IF NOT EXISTS slow_moving_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    slow_moving_id CHAR(36),
    product_master_id CHAR(36),
    is_damaged BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_slow_moving_items_id ON slow_moving_items(id);
CREATE INDEX IF NOT EXISTS idx_slow_moving_items_product_master_slow_moving ON slow_moving_items(product_master_id, slow_moving_id);

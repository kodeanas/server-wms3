-- UserClassLog
CREATE TABLE IF NOT EXISTS user_class_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    buyer_id CHAR(36),
    prev_class_id CHAR(36),
    new_class_id CHAR(36),
    order_id CHAR(36),
    change_type VARCHAR(255),
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_user_class_logs_prev_new_order_buyer ON user_class_logs(prev_class_id, new_class_id, order_id, buyer_id);

-- MovementProduct
CREATE TABLE IF NOT EXISTS movement_products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id CHAR(36),
    is_sku BOOLEAN DEFAULT false,
    movement_type VARCHAR(50) CHECK (movement_type IN ('bundle', 'unbundle', 'abnormal')),
    type_out VARCHAR(50) CHECK (type_out IN ('regular_sales', 'cargo', 'scrap', 'qcd', 'transfer')),
    from_location VARCHAR(255),
    to_location VARCHAR(255),
    qty INT,
    datetime TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_movement_products_product_id ON movement_products(product_id);

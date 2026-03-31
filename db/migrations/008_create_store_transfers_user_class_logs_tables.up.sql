-- Create Store Transfers Table
CREATE TABLE IF NOT EXISTS store_transfers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    store_id CHAR(36),
    datetime TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    total_item INT,
    total_price DECIMAL(15,2),
    total_price_warehouse DECIMAL(15,2),
    status VARCHAR(50) CHECK (status IN ('on_progress', 'done')),
    user_id CHAR(36),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_store_transfers_id ON store_transfers(id);
CREATE INDEX IF NOT EXISTS idx_store_transfers_user_id ON store_transfers(user_id);

-- Create Store Transfer Bags Table
CREATE TABLE IF NOT EXISTS store_transfer_bags (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    store_transfer_id CHAR(36),
    bag_id CHAR(36),
    quantity INT,
    total_price DECIMAL(15,2),
    total_cogs DECIMAL(15,2),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_store_transfer_bags_id ON store_transfer_bags(id);
CREATE INDEX IF NOT EXISTS idx_store_transfer_bags_transfer_bag ON store_transfer_bags(store_transfer_id, bag_id);

-- Create User Class Logs Table
CREATE TABLE IF NOT EXISTS user_class_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    buyer_id CHAR(36),
    prev_class_id CHAR(36),
    new_class_id CHAR(36),
    order_id CHAR(36),
    change_type VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_user_class_logs_id ON user_class_logs(id);
CREATE INDEX IF NOT EXISTS idx_user_class_logs_prev_new_order_buyer ON user_class_logs(prev_class_id, new_class_id, order_id, buyer_id);

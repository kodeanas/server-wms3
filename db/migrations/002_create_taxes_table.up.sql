-- Tax
CREATE TABLE IF NOT EXISTS taxes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tax INT NOT NULL,
    is_active BOOLEAN DEFAULT true
);

CREATE INDEX IF NOT EXISTS idx_taxes_tax ON taxes(tax);

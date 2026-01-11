-- Create mst_currencies table (previously sys_currency)
-- Stores currency reference data
CREATE TABLE IF NOT EXISTS mst_currencies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- Primary key using UUID
    
    -- Main data
    code VARCHAR(5) NOT NULL,       -- ISO currency code (e.g., IDR, USD)
    description VARCHAR(30),        -- Full currency name (e.g., Rupiah)
    sort_order INT NOT NULL DEFAULT 0, -- Display priority
    
    -- Audit fields
    created_by UUID,    -- User who created this record
    updated_by UUID,    -- User who last updated this record
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- Record creation timestamp
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- Record update timestamp
    deleted_at TIMESTAMP WITH TIME ZONE -- Soft deletion timestamp
);

-- Create indexes
CREATE UNIQUE INDEX IF NOT EXISTS idx_mst_currencies_code ON mst_currencies(code) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_mst_currencies_sort_order ON mst_currencies(sort_order) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_mst_currencies_deleted_at ON mst_currencies(deleted_at);

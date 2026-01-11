-- Create mst_banks table (previously sys_banks)
-- Stores bank reference data
CREATE TABLE IF NOT EXISTS mst_banks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- Primary key using UUID
    
    -- Main data
    name VARCHAR(255) NOT NULL, -- Full name of the bank
    code VARCHAR(10) NOT NULL,  -- Unique bank code (Swift/BI code)
    original_code VARCHAR(50),  -- Legacy system bank code
    description VARCHAR(255),   -- General bank description
    
    -- Classification
    category VARCHAR(100),      -- Bank category (e.g., Persero, Swasta Nasional)
    sub_category VARCHAR(100),  -- Sub-category or grouping
    is_sharia BOOLEAN DEFAULT false, -- Whether the bank follows Sharia principles
    
    -- Audit fields
    created_by UUID,    -- User who created this record
    updated_by UUID,    -- User who last updated this record
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- Record creation timestamp
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- Record update timestamp
    deleted_at TIMESTAMP WITH TIME ZONE -- Soft deletion timestamp
);

-- Create indexes
CREATE UNIQUE INDEX IF NOT EXISTS idx_mst_banks_code ON mst_banks(code) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_mst_banks_category ON mst_banks(category) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_mst_banks_is_sharia ON mst_banks(is_sharia) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_mst_banks_deleted_at ON mst_banks(deleted_at);

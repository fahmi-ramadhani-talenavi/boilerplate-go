-- Create mst_months table (previously sys_monts)
-- Stores month reference data
CREATE TABLE IF NOT EXISTS mst_months (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- Primary key using UUID
    
    -- Main data
    code VARCHAR(2) NOT NULL,       -- Two-digit month code (e.g., '01' for January)
    name VARCHAR(20) NOT NULL,      -- Full month name
    sort_order INT NOT NULL DEFAULT 0, -- Display priority
    
    -- Audit fields
    created_by UUID,    -- User who created this record
    updated_by UUID,    -- User who last updated this record
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- Record creation timestamp
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- Record update timestamp
    deleted_at TIMESTAMP WITH TIME ZONE -- Soft deletion timestamp
);

-- Create indexes
CREATE UNIQUE INDEX IF NOT EXISTS idx_mst_months_code ON mst_months(code) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_mst_months_sort_order ON mst_months(sort_order) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_mst_months_deleted_at ON mst_months(deleted_at);

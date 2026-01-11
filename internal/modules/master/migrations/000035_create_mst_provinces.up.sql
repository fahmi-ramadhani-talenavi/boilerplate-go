-- Create mst_provinces table (previously sys_province)
-- Stores province reference data
CREATE TABLE IF NOT EXISTS mst_provinces (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- Primary key using UUID
    
    -- Main data
    code VARCHAR(15) NOT NULL,       -- Unique province (Provinsi) code
    description VARCHAR(100),       -- Province name
    sort_order INT NOT NULL DEFAULT 0, -- Display priority
    
    -- Audit fields
    created_by UUID,    -- User who created this record
    updated_by UUID,    -- User who last updated this record
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- Record creation timestamp
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- Record update timestamp
    deleted_at TIMESTAMP WITH TIME ZONE -- Soft deletion timestamp
);

-- Create indexes
CREATE UNIQUE INDEX IF NOT EXISTS idx_mst_provinces_code ON mst_provinces(code) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_mst_provinces_description ON mst_provinces(description) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_mst_provinces_deleted_at ON mst_provinces(deleted_at);

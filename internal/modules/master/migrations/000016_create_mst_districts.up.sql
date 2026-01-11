-- Create mst_districts table (previously sys_district)
-- Stores district/kecamatan reference data
CREATE TABLE IF NOT EXISTS mst_districts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- Primary key using UUID
    
    -- Main data
    code VARCHAR(15) NOT NULL,      -- District/Kecamatan unique code
    description VARCHAR(100),       -- District name
    sort_order INT NOT NULL DEFAULT 0, -- Display priority
    
    -- Province reference
    province_code VARCHAR(15),      -- Legacy province reference code
    province_id UUID,                -- Link to mst_provinces (Provinsi)
    
    -- Audit fields
    created_by UUID,    -- User who created this record
    updated_by UUID,    -- User who last updated this record
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- Record creation timestamp
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- Record update timestamp
    deleted_at TIMESTAMP WITH TIME ZONE -- Soft deletion timestamp
);

-- Create indexes
CREATE UNIQUE INDEX IF NOT EXISTS idx_mst_districts_code ON mst_districts(code) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_mst_districts_description ON mst_districts(description) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_mst_districts_sort_order ON mst_districts(sort_order) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_mst_districts_province_id ON mst_districts(province_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_mst_districts_deleted_at ON mst_districts(deleted_at);

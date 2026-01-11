-- Create mst_sub_districts table (previously sys_subdistrict)
-- Stores sub-district/kelurahan reference data
CREATE TABLE IF NOT EXISTS mst_sub_districts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- Primary key using UUID
    
    -- Main data
    code VARCHAR(15) NOT NULL,       -- Unique sub-district (Kelurahan/Desa) code
    description VARCHAR(100),       -- Sub-district name
    sort_order INT NOT NULL DEFAULT 0, -- Display priority
    
    -- Geography references
    province_code VARCHAR(15),      -- Legacy province reference code
    district_code VARCHAR(15),      -- Legacy district reference code
    district_id UUID REFERENCES mst_districts(id) ON DELETE CASCADE, -- Link to parent district
    
    -- Audit fields
    created_by UUID,    -- User who created this record
    updated_by UUID,    -- User who last updated this record
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- Record creation timestamp
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- Record update timestamp
    deleted_at TIMESTAMP WITH TIME ZONE -- Soft deletion timestamp
);

-- Create indexes
CREATE UNIQUE INDEX IF NOT EXISTS idx_mst_sub_districts_code ON mst_sub_districts(code) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_mst_sub_districts_description ON mst_sub_districts(description) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_mst_sub_districts_district_id ON mst_sub_districts(district_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_mst_sub_districts_deleted_at ON mst_sub_districts(deleted_at);

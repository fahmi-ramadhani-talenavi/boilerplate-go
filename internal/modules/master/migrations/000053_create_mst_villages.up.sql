-- Create mst_villages table (previously sys_village)
-- Stores village/kelurahan/desa reference data, completing the geographic hierarchy
CREATE TABLE IF NOT EXISTS mst_villages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- Primary key using UUID
    
    -- Main data
    code VARCHAR(15) NOT NULL,       -- Unique village (Kelurahan/Desa) code
    description VARCHAR(100),       -- Village name
    sort_order INT NOT NULL DEFAULT 0, -- Display priority
    
    -- Geography references
    province_code VARCHAR(15),      -- Legacy province reference code
    district_code VARCHAR(15),      -- Legacy district reference code
    sub_district_code VARCHAR(15),  -- Legacy sub-district reference code
    sub_district_id UUID REFERENCES mst_sub_districts(id) ON DELETE CASCADE, -- Link to parent sub-district
    
    -- Audit fields
    created_by UUID,    -- User who created this record
    updated_by UUID,    -- User who last updated this record
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- Record creation timestamp
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- Record update timestamp
    deleted_at TIMESTAMP WITH TIME ZONE -- Soft deletion timestamp
);

-- Create indexes
CREATE UNIQUE INDEX IF NOT EXISTS idx_mst_villages_code ON mst_villages(code) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_mst_villages_description ON mst_villages(description) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_mst_villages_sub_district_id ON mst_villages(sub_district_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_mst_villages_deleted_at ON mst_villages(deleted_at);

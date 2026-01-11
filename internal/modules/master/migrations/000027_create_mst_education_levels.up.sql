-- Create mst_education_levels table (previously sys_last_edu)
-- Stores education level reference data
CREATE TABLE IF NOT EXISTS mst_education_levels (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- Primary key using UUID
    
    -- Main data
    description VARCHAR(255) NOT NULL, -- Education level (e.g., SD, SMP, SMA, S1)
    sort_order INT DEFAULT 0,          -- Display priority
    
    -- Audit fields
    created_by UUID,    -- User who created this record
    updated_by UUID,    -- User who last updated this record
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- Record creation timestamp
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- Record update timestamp
    deleted_at TIMESTAMP WITH TIME ZONE -- Soft deletion timestamp
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_mst_education_levels_description ON mst_education_levels(description) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_mst_education_levels_sort_order ON mst_education_levels(sort_order) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_mst_education_levels_deleted_at ON mst_education_levels(deleted_at);

-- Create mst_main_branches table (previously sys_main_branch)
-- Stores main branch/regional office reference data
CREATE TABLE IF NOT EXISTS mst_main_branches (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- Primary key using UUID
    
    -- Main data
    code VARCHAR(255) NOT NULL,      -- Unique main branch/regional code
    name VARCHAR(255) NOT NULL,      -- Main branch name
    is_active BOOLEAN NOT NULL DEFAULT true, -- Whether the main branch is operational
    sort_order INT NOT NULL DEFAULT 0,      -- Display priority
    
    -- Relations
    area_id UUID NOT NULL, -- Link to mst_areas for regional grouping
    
    -- Audit fields
    created_by UUID,    -- User who created this record
    updated_by UUID,    -- User who last updated this record
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- Record creation timestamp
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- Record update timestamp
    deleted_at TIMESTAMP WITH TIME ZONE -- Soft deletion timestamp
);

-- Create indexes
CREATE UNIQUE INDEX IF NOT EXISTS idx_mst_main_branches_code ON mst_main_branches(code) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_mst_main_branches_area_id ON mst_main_branches(area_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_mst_main_branches_is_active ON mst_main_branches(is_active) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_mst_main_branches_sort_order ON mst_main_branches(sort_order) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_mst_main_branches_deleted_at ON mst_main_branches(deleted_at);

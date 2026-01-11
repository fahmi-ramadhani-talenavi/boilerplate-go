-- Create mst_branches table (previously sys_branchs)
-- Stores branch/office location data
CREATE TABLE IF NOT EXISTS mst_branches (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- Primary key using UUID
    
    -- Branch identification
    code VARCHAR(255) NOT NULL,        -- Unique branch office code
    description VARCHAR(255) NOT NULL, -- Branch office name/description
    
    -- Contact information
    address VARCHAR(255) NOT NULL, -- Full physical address of the office
    phone VARCHAR(255) NOT NULL,   -- Primary contact phone number
    phone2 VARCHAR(255),           -- Secondary contact or fax number
    
    -- Branch classification
    branch_type SMALLINT NOT NULL DEFAULT 2, -- Category: 1 for Head Office (Pusat), 2 for Regional (Cabang)
    is_default BOOLEAN NOT NULL DEFAULT false, -- Whether this is the primary system branch
    is_active BOOLEAN NOT NULL DEFAULT true, -- Whether the branch is currently operational
    flag SMALLINT NOT NULL DEFAULT 0,       -- Internal processing flag
    sort_order INT NOT NULL DEFAULT 0,      -- Display priority in UI lists
    
    -- Relations
    area_id UUID NOT NULL,   -- Link to mst_areas for regional grouping
    main_branch_id UUID,     -- Self-reference link for sub-branch hierarchy
    
    -- Audit fields
    created_by UUID,    -- User who created this record
    updated_by UUID,    -- User who last updated this record
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- Record creation timestamp
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- Record update timestamp
    deleted_at TIMESTAMP WITH TIME ZONE -- Soft deletion timestamp
);

-- Create indexes
CREATE UNIQUE INDEX IF NOT EXISTS idx_mst_branches_code ON mst_branches(code) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_mst_branches_area_id ON mst_branches(area_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_mst_branches_main_branch_id ON mst_branches(main_branch_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_mst_branches_branch_type ON mst_branches(branch_type) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_mst_branches_is_active ON mst_branches(is_active) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_mst_branches_is_default ON mst_branches(is_default) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_mst_branches_sort_order ON mst_branches(sort_order) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_mst_branches_created_by ON mst_branches(created_by);
CREATE INDEX IF NOT EXISTS idx_mst_branches_updated_by ON mst_branches(updated_by);
CREATE INDEX IF NOT EXISTS idx_mst_branches_deleted_at ON mst_branches(deleted_at);

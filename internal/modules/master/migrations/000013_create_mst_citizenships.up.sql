-- Create mst_citizenships table (previously sys_citizenship)
-- Stores citizenship/nationality reference data
CREATE TABLE IF NOT EXISTS mst_citizenships (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- Primary key using UUID
    
    -- Main data
    code VARCHAR(5) NOT NULL,       -- Citizenship code (e.g., WNI, WNA)
    description VARCHAR(30),        -- Full description (e.g., Warga Negara Indonesia)
    sort_order INT NOT NULL DEFAULT 0, -- Display priority
    
    -- Audit fields
    created_by UUID,    -- User who created this record
    updated_by UUID,    -- User who last updated this record
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- Record creation timestamp
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- Record update timestamp
    deleted_at TIMESTAMP WITH TIME ZONE -- Soft deletion timestamp
);

-- Create indexes
CREATE UNIQUE INDEX IF NOT EXISTS idx_mst_citizenships_code ON mst_citizenships(code) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_mst_citizenships_sort_order ON mst_citizenships(sort_order) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_mst_citizenships_deleted_at ON mst_citizenships(deleted_at);

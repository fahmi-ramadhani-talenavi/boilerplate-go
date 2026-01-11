-- Create mst_contribution_types table (previously sys_tipe_iuran)
-- Stores reference data for different types of contributions/iuran
CREATE TABLE IF NOT EXISTS mst_contribution_types (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- Primary key using UUID
    
    -- Main data
    description VARCHAR(30),        -- Type of contribution (e.g., Iuran Peserta, Pemberi Kerja)
    sort_order INT NOT NULL DEFAULT 0, -- Display priority
    
    -- Audit fields
    created_by UUID,    -- User who created this record
    updated_by UUID,    -- User who last updated this record
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- Record creation timestamp
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- Record update timestamp
    deleted_at TIMESTAMP WITH TIME ZONE -- Soft deletion timestamp
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_mst_contribution_types_sort_order ON mst_contribution_types(sort_order) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_mst_contribution_types_deleted_at ON mst_contribution_types(deleted_at);

-- Create mst_giro_types table (previously sys_giro_types)
-- Stores reference data for giro/current account types and their ATS configurations
CREATE TABLE IF NOT EXISTS mst_giro_types (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- Primary key using UUID
    
    -- Main data
    label VARCHAR(255) NOT NULL,    -- Display name for the giro type
    ats_types TEXT,                 -- Automatic Transfer System (ATS) type configurations
    sort_order INT NOT NULL DEFAULT 0, -- Display priority
    
    -- Audit fields
    created_by UUID,    -- User who created this record
    updated_by UUID,    -- User who last updated this record
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- Record creation timestamp
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- Record update timestamp
    deleted_at TIMESTAMP WITH TIME ZONE -- Soft deletion timestamp
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_mst_giro_types_sort_order ON mst_giro_types(sort_order) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_mst_giro_types_deleted_at ON mst_giro_types(deleted_at);

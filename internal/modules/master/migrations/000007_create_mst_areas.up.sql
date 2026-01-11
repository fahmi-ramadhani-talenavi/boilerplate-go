-- Create mst_areas table (previously sys_areas)
-- Stores geographic area/region reference data
CREATE TABLE IF NOT EXISTS mst_areas (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- Primary key using UUID
    
    -- Main data
    code VARCHAR(50) NOT NULL,      -- Area/Region unique code
    name VARCHAR(255) NOT NULL,     -- Area/Region display name
    is_active BOOLEAN DEFAULT true, -- Whether the area is active for selection
    
    -- Audit fields
    created_by UUID,    -- User who created this record
    updated_by UUID,    -- User who last updated this record
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- Record creation timestamp
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- Record update timestamp
    deleted_at TIMESTAMP WITH TIME ZONE -- Soft deletion timestamp
);

-- Create indexes
CREATE UNIQUE INDEX IF NOT EXISTS idx_mst_areas_code ON mst_areas(code) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_mst_areas_is_active ON mst_areas(is_active) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_mst_areas_deleted_at ON mst_areas(deleted_at);

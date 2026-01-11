-- Create mst_genders table (previously sys_gender)
-- Stores gender reference data
CREATE TABLE IF NOT EXISTS mst_genders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- Primary key using UUID
    
    -- Main data
    code VARCHAR(5) NOT NULL,       -- Gender code (e.g., L, P)
    description VARCHAR(30),        -- Full gender description (e.g., Laki-laki, Perempuan)
    sort_order INT NOT NULL DEFAULT 0, -- Display priority
    
    -- Audit fields
    created_by UUID,    -- User who created this record
    updated_by UUID,    -- User who last updated this record
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- Record creation timestamp
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- Record update timestamp
    deleted_at TIMESTAMP WITH TIME ZONE -- Soft deletion timestamp
);

-- Create indexes
CREATE UNIQUE INDEX IF NOT EXISTS idx_mst_genders_code ON mst_genders(code) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_mst_genders_sort_order ON mst_genders(sort_order) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_mst_genders_deleted_at ON mst_genders(deleted_at);

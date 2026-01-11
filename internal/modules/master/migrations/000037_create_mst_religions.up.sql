-- Create mst_religions table (previously sys_religion)
-- Stores religion reference data
CREATE TABLE IF NOT EXISTS mst_religions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- Primary key using UUID
    
    -- Main data
    description VARCHAR(30) NOT NULL, -- Religion name (e.g., Islam, Kristen, Katolik, Hindu, Budha)
    
    -- Audit fields
    created_by UUID,    -- User who created this record
    updated_by UUID,    -- User who last updated this record
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- Record creation timestamp
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- Record update timestamp
    deleted_at TIMESTAMP WITH TIME ZONE -- Soft deletion timestamp
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_mst_religions_description ON mst_religions(description) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_mst_religions_deleted_at ON mst_religions(deleted_at);

-- Create mst_salutations table (previously sys_salutasi)
-- Stores salutation reference data
CREATE TABLE IF NOT EXISTS mst_salutations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- Primary key using UUID
    
    -- Main data
    code VARCHAR(10) NOT NULL,      -- Salutation code (e.g., MR, MRS, MS)
    description VARCHAR(30),        -- Full salutation display (e.g., Bapak, Ibu, Sdr)
    
    -- Audit fields
    created_by UUID,    -- User who created this record
    updated_by UUID,    -- User who last updated this record
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- Record creation timestamp
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- Record update timestamp
    deleted_at TIMESTAMP WITH TIME ZONE -- Soft deletion timestamp
);

-- Create indexes
CREATE UNIQUE INDEX IF NOT EXISTS idx_mst_salutations_code ON mst_salutations(code) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_mst_salutations_description ON mst_salutations(description) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_mst_salutations_deleted_at ON mst_salutations(deleted_at);

-- Create mst_beneficiaries table (previously sys_ahli_waris)
-- Stores beneficiary/heir types for member accounts
CREATE TABLE IF NOT EXISTS mst_beneficiaries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- Primary key using UUID
    
    -- Main data
    description VARCHAR(255) NOT NULL, -- Beneficiary type (e.g., Ahli Waris, Istri, Anak)
    
    -- Audit fields
    created_by UUID,    -- User who created this record
    updated_by UUID,    -- User who last updated this record
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- Record creation timestamp
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- Record update timestamp
    deleted_at TIMESTAMP WITH TIME ZONE -- Soft deletion timestamp
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_mst_beneficiaries_deleted_at ON mst_beneficiaries(deleted_at);

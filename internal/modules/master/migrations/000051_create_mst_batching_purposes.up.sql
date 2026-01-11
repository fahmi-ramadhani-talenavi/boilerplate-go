-- Create mst_batching_purposes table (previously sys_tujuan_batching)
-- Stores reference data for different purposes of batching operations
CREATE TABLE IF NOT EXISTS mst_batching_purposes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- Primary key using UUID
    
    -- Main data
    description VARCHAR(255) NOT NULL, -- Description of the batching purpose (e.g., Membership Update, Financial Processing)
    
    -- Audit fields
    created_by UUID,    -- User who created this record
    updated_by UUID,    -- User who last updated this record
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- Record creation timestamp
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- Record update timestamp
    deleted_at TIMESTAMP WITH TIME ZONE -- Soft deletion timestamp
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_mst_batching_purposes_deleted_at ON mst_batching_purposes(deleted_at);

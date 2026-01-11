-- Create mst_jobs table (previously sys_pekerjaan)
-- Stores job reference data
CREATE TABLE IF NOT EXISTS mst_jobs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- Primary key using UUID
    
    -- Main data
    description VARCHAR(255) NOT NULL, -- Job title or occupation category (e.g., PNS, Swasta, TNI/POLRI)
    
    -- Audit fields
    created_by UUID,    -- User who created this record
    updated_by UUID,    -- User who last updated this record
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- Record creation timestamp
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- Record update timestamp
    deleted_at TIMESTAMP WITH TIME ZONE -- Soft deletion timestamp
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_mst_jobs_description ON mst_jobs(description) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_mst_jobs_deleted_at ON mst_jobs(deleted_at);

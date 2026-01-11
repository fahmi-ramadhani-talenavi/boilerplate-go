-- Create mst_job_statuses table (previously mst_pekerjaan_statuses)
-- Stores job status reference data
CREATE TABLE IF NOT EXISTS mst_job_statuses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- Primary key using UUID
    
    -- Main data
    description VARCHAR(255) NOT NULL, -- Employment status (e.g., Tetap, Kontrak, Pensiun)
    
    -- Audit fields
    created_by UUID,    -- User who created this record
    updated_by UUID,    -- User who last updated this record
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- Record creation timestamp
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- Record update timestamp
    deleted_at TIMESTAMP WITH TIME ZONE -- Soft deletion timestamp
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_mst_job_statuses_description ON mst_job_statuses(description) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_mst_job_statuses_deleted_at ON mst_job_statuses(deleted_at);

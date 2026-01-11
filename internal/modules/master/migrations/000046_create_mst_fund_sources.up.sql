-- Create mst_fund_sources table (previously sy_sumber_danas)
-- Stores fund source reference data
CREATE TABLE IF NOT EXISTS mst_fund_sources (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- Primary key using UUID
    
    -- Main data
    description VARCHAR(255) NOT NULL, -- Source of funds (e.g., Salary, Business, Inheritance)
    is_personal BOOLEAN NOT NULL DEFAULT false, -- Whether the source is for personal or corporate use
    
    -- Audit fields
    created_by UUID,    -- User who created this record
    updated_by UUID,    -- User who last updated this record
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- Record creation timestamp
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- Record update timestamp
    deleted_at TIMESTAMP WITH TIME ZONE -- Soft deletion timestamp
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_mst_fund_sources_description ON mst_fund_sources(description) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_mst_fund_sources_deleted_at ON mst_fund_sources(deleted_at);

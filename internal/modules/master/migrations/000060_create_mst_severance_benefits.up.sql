-- Create mst_severance_benefits table (previously severance_benefits)
-- Stores specific monetary benefits or multipliers for different reasons of leaving
CREATE TABLE IF NOT EXISTS mst_severance_benefits (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- Primary key using UUID
    
    -- Benefit Configuration
    termination_reason VARCHAR(255) NOT NULL, -- The reason for leaving (e.g., Pensiun, PHK)
    severance_pay DECIMAL(20,2) NOT NULL DEFAULT 0.00, -- Amount of severance pay (Uang Pesangon)
    service_pay DECIMAL(20,2) NOT NULL DEFAULT 0.00,   -- Amount of service pay (Uang Penghargaan Masa Kerja / UPMK)
    
    -- Status and Dates
    effective_date DATE NOT NULL,   -- Date when this benefit configuration becomes effective
    
    -- Audit fields
    created_by UUID,    -- User who created this record
    updated_by UUID,    -- User who last updated this record
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- Record creation timestamp
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- Record update timestamp
    deleted_at TIMESTAMP WITH TIME ZONE -- Soft deletion timestamp
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_mst_severance_benefits_termination_reason ON mst_severance_benefits(termination_reason) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_mst_severance_benefits_effective_date ON mst_severance_benefits(effective_date) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_mst_severance_benefits_deleted_at ON mst_severance_benefits(deleted_at);

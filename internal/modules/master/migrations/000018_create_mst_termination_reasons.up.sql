-- Create mst_termination_reasons table (previously sys_dkp_reason_ciptaker)
-- Stores reasons for employment termination and their respective severance multipliers
CREATE TABLE IF NOT EXISTS mst_termination_reasons (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- Primary key using UUID
    
    -- Main data
    reason_name TEXT NOT NULL, -- Employment termination reason (e.g., Pensiun, PHK)
    
    -- Multipliers (used to multiply the base months from service periods)
    severance_multiplier DOUBLE PRECISION NOT NULL DEFAULT 1.0,   -- Multiplier for Uang Pesangon (Ciptaker)
    service_pay_multiplier DOUBLE PRECISION NOT NULL DEFAULT 1.0, -- Multiplier for Uang Penghargaan Masa Kerja (Ciptaker)
    
    -- Audit fields
    created_by UUID,    -- User who created this record
    updated_by UUID,    -- User who last updated this record
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- Record creation timestamp
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- Record update timestamp
    deleted_at TIMESTAMP WITH TIME ZONE -- Soft deletion timestamp
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_mst_termination_reasons_deleted_at ON mst_termination_reasons(deleted_at);

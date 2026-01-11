-- Create sys_pension_ages table (previously sys_usia_pensiun)
-- Stores configuration for normal and early pension ages
CREATE TABLE IF NOT EXISTS sys_pension_ages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- Primary key using UUID
    
    -- Main data
    effective_date DATE,            -- Date when this pension age configuration becomes effective
    normal_pension_age INT NOT NULL DEFAULT 0, -- Age for normal pension benefit eligibility
    early_pension_age INT NOT NULL DEFAULT 0,  -- Age for early pension benefit eligibility
    reference_regulation TEXT,      -- Legal or company regulation governing these ages
    description TEXT,               -- Additional notes or details about this configuration
    
    -- Audit fields
    created_by UUID,    -- User who created this record
    updated_by UUID,    -- User who last updated this record
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- Record creation timestamp
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- Record update timestamp
    deleted_at TIMESTAMP WITH TIME ZONE -- Soft deletion timestamp
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_sys_pension_ages_effective_date ON sys_pension_ages(effective_date) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_sys_pension_ages_deleted_at ON sys_pension_ages(deleted_at);

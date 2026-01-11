-- Create mst_severance_service_periods table (previously sys_dkp_dedication_ciptaker)
-- Stores severance pay and service pay rates based on duration of employment
CREATE TABLE IF NOT EXISTS mst_severance_service_periods (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- Primary key using UUID
    
    -- Service duration in years
    start_year INT NOT NULL, -- Lower bound of employment duration (Years)
    end_year INT,           -- Upper bound of employment duration (Years, NULL for no limit)
    
    -- Rates (expressed in months of salary)
    severance_pay_months DOUBLE PRECISION NOT NULL DEFAULT 0, -- Multiplier for severance pay (Uang Pesangon)
    service_pay_months DOUBLE PRECISION NOT NULL DEFAULT 0,   -- Multiplier for service pay (Uang Penghargaan Masa Kerja)
    
    -- Additional info
    description VARCHAR(255), -- Descriptive name for the period bracket
    
    -- Audit fields
    created_by UUID,    -- User who created this record
    updated_by UUID,    -- User who last updated this record
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- Record creation timestamp
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- Record update timestamp
    deleted_at TIMESTAMP WITH TIME ZONE -- Soft deletion timestamp
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_mst_severance_service_periods_duration ON mst_severance_service_periods(start_year, end_year) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_mst_severance_service_periods_deleted_at ON mst_severance_service_periods(deleted_at);

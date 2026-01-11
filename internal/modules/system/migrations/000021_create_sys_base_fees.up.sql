-- Create sys_base_fees table (previously sys_fee_base)
-- Stores base fee configurations such as registration and administration fees
CREATE TABLE IF NOT EXISTS sys_base_fees (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- Primary key using UUID
    
    -- Status and Flags
    is_active BOOLEAN NOT NULL DEFAULT false, -- Whether this fee configuration is currently active
    is_default BOOLEAN NOT NULL DEFAULT false, -- Whether this is the default fee set for new accounts
    
    -- Validity Period
    effective_start_date TIMESTAMP WITH TIME ZONE, -- Start date for the fee set validity
    effective_end_date TIMESTAMP WITH TIME ZONE,   -- End date for the fee set validity
    
    -- Fee Details
    group_name VARCHAR(255) NOT NULL,    -- Name of the fee grouping or category
    registration_fee BIGINT NOT NULL DEFAULT 0, -- One-time registration fee amount
    administration_fee BIGINT NOT NULL DEFAULT 0, -- Recurring administration fee amount
    
    -- Audit fields
    created_by UUID,    -- User who created this record
    updated_by UUID,    -- User who last updated this record
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- Record creation timestamp
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- Record update timestamp
    deleted_at TIMESTAMP WITH TIME ZONE -- Soft deletion timestamp
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_sys_base_fees_group_name ON sys_base_fees(group_name) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_sys_base_fees_is_active ON sys_base_fees(is_active) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_sys_base_fees_is_default ON sys_base_fees(is_default) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_sys_base_fees_effective_dates ON sys_base_fees(effective_start_date, effective_end_date) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_sys_base_fees_deleted_at ON sys_base_fees(deleted_at);

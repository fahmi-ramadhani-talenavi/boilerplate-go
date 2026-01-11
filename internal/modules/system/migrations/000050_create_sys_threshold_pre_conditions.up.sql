-- Create sys_threshold_pre_conditions table (previously sys_treshold_pre_conditions)
-- Stores various numeric limits and conditions for pension transactions
CREATE TABLE IF NOT EXISTS sys_threshold_pre_conditions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- Primary key using UUID
    
    -- Partial Withdrawal Thresholds (Tarik Sebagian)
    min_participation_partial_withdrawal INT NOT NULL DEFAULT 0, -- Minimum months of participation required
    withdrawal_interval_partial_withdrawal INT NOT NULL DEFAULT 0, -- Minimum months between withdrawals
    min_balance_partial_withdrawal BIGINT NOT NULL DEFAULT 0, -- Minimum account balance to allow withdrawal
    min_amount_partial_withdrawal BIGINT NOT NULL DEFAULT 0, -- Minimum amount allowed for a single withdrawal
    
    -- Benefit Termination Thresholds (Manfaat Pensiun)
    min_withdrawal_benefit_termination BIGINT NOT NULL DEFAULT 0, -- Minimum amount for termination benefit withdrawal
    max_balance_benefit_termination BIGINT NOT NULL DEFAULT 0, -- Maximum balance threshold for termination benefit logic
    
    -- Pension Transfer Thresholds (Pengalihan Dana)
    max_participation_pension_transfer INT NOT NULL DEFAULT 0, -- Maximum participation duration for certain transfer logic
    
    -- Audit fields
    created_by UUID,    -- User who created this record
    updated_by UUID,    -- User who last updated this record
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- Record creation timestamp
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- Record update timestamp
    deleted_at TIMESTAMP WITH TIME ZONE -- Soft deletion timestamp
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_sys_threshold_pre_conditions_deleted_at ON sys_threshold_pre_conditions(deleted_at);

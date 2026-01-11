-- Create sys_transaction_fees table (previously sys_transaction_fees)
-- Stores configuration for various transaction-related fees with effective dating
CREATE TABLE IF NOT EXISTS sys_transaction_fees (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- Primary key using UUID
    grouping_name VARCHAR(255), -- Name of the fee group/category
    is_default BOOLEAN NOT NULL DEFAULT false, -- Whether this is the default fee group
    is_active BOOLEAN NOT NULL DEFAULT true, -- Whether this fee configuration is active
    
    -- Transaction Fee Details (Pensions)
    benefit_pension_quit_work DECIMAL(12,2), -- Fee for pension benefit payout after quitting work
    move_pension DECIMAL(12,2), -- Fee for moving pension funds
    benefit_pension_claim DECIMAL(12,2), -- Fee for claiming pension benefits
    benefit_pension_yearly_admin DECIMAL(12,2), -- Annual administration fee for pension accounts
    
    -- Transaction Fee Details (Investment Packages)
    move_package_invest_gt_2y DECIMAL(12,2), -- Fee for switching investment packages after 2 years
    move_package_invest_lt_2y DECIMAL(12,2), -- Fee for switching investment packages before 2 years
    
    -- Transaction Fee Details (Withdrawals)
    claim_withdrawal_partial_first DECIMAL(12,2), -- Fee for the first partial withdrawal claim
    claim_withdrawal_partial_second DECIMAL(12,2), -- Fee for the second partial withdrawal claim
    
    -- Transaction Fee Details (Fund Transfers Out)
    move_fund_out_lt_3y DECIMAL(12,2), -- Fee for transferring funds out before 3 years
    move_fund_out_bw_3y DECIMAL(12,2), -- Fee for transferring funds out between specific periods
    move_fund_out_gt_3y DECIMAL(12,2), -- Fee for transferring funds out after 3 years
    
    -- Effective Dates
    effective_at_start TIMESTAMP WITH TIME ZONE, -- Start date when these fees apply
    effective_at_end TIMESTAMP WITH TIME ZONE,   -- End date when these fees stop applying
    
    -- Audit fields
    created_by UUID,    -- User who created this record
    updated_by UUID,    -- User who last updated this record
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- Record creation timestamp
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- Record update timestamp
    deleted_at TIMESTAMP WITH TIME ZONE -- Soft deletion timestamp
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_sys_transaction_fees_grouping_name ON sys_transaction_fees(grouping_name) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_sys_transaction_fees_is_default ON sys_transaction_fees(is_default) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_sys_transaction_fees_is_active ON sys_transaction_fees(is_active) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_sys_transaction_fees_effective_dates ON sys_transaction_fees(effective_at_start, effective_at_end) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_sys_transaction_fees_deleted_at ON sys_transaction_fees(deleted_at);

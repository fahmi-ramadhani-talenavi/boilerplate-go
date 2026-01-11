-- Create sys_bank_fees table
-- Stores bank fee configurations per transaction type
CREATE TABLE IF NOT EXISTS sys_bank_fees (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- Primary key using UUID
    
    -- Main data
    transaction_type VARCHAR(100) NOT NULL, -- Type of transaction (e.g., Withdrawal, Topup)
    description TEXT NOT NULL,           -- Description of the fee rule
    bank_fee BIGINT NOT NULL DEFAULT 0,    -- Amount of fee charged by the bank
    pension_fund_income BIGINT NOT NULL DEFAULT 0, -- Portion of fee income for DPLK
    
    -- Audit fields
    created_by UUID,    -- User who created this record
    updated_by UUID,    -- User who last updated this record
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- Record creation timestamp
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- Record update timestamp
    deleted_at TIMESTAMP WITH TIME ZONE -- Soft deletion timestamp
);

-- Create indexes
CREATE UNIQUE INDEX IF NOT EXISTS idx_sys_bank_fees_transaction_type ON sys_bank_fees(transaction_type) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_sys_bank_fees_deleted_at ON sys_bank_fees(deleted_at);

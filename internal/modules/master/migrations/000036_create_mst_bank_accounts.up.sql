-- Create mst_bank_accounts table (previously sys_rekening)
-- Stores bank account reference data
CREATE TABLE IF NOT EXISTS mst_bank_accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- Primary key using UUID
    
    -- Main data
    bank_name VARCHAR(255),           -- Full name of the bank
    branch_name VARCHAR(255),         -- Bank branch office where the account is held
    account_holder_name VARCHAR(255), -- Full name of the bank account owner
    account_number VARCHAR(100),      -- The bank account number
    
    -- Audit fields
    created_by UUID,    -- User who created this record
    updated_by UUID,    -- User who last updated this record
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- Record creation timestamp
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- Record update timestamp
    deleted_at TIMESTAMP WITH TIME ZONE -- Soft deletion timestamp
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_mst_bank_accounts_bank_name ON mst_bank_accounts(bank_name) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_mst_bank_accounts_account_number ON mst_bank_accounts(account_number) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_mst_bank_accounts_deleted_at ON mst_bank_accounts(deleted_at);

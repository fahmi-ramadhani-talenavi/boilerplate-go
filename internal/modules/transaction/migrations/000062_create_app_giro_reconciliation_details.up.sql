-- Create app_giro_reconciliation_details table (previously reconciliation_giro_details)
-- Stores granular transaction details for each giro reconciliation record
CREATE TABLE IF NOT EXISTS app_giro_reconciliation_details (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- Primary key using UUID
    
    -- References
    giro_id UUID REFERENCES app_giro_reconciliations(id) ON DELETE CASCADE, -- Link to the parent giro reconciliation
    batch_id UUID,                  -- Reference to the system batch (previously sb_id)
    legacy_giro_id BIGINT,          -- Legacy reference for giro_id
    legacy_batch_id BIGINT,         -- Legacy reference for batch_id
    
    -- Transaction Data
    transaction_date TIMESTAMP WITH TIME ZONE NOT NULL, -- Date and time of the transaction
    upload_date TIMESTAMP WITH TIME ZONE NOT NULL,      -- Date and time the record was uploaded
    giro_number VARCHAR(255),       -- Associated giro number
    apac_no VARCHAR(255),           -- APAC (Account/Payment) reference number
    reference_number VARCHAR(255),  -- External reference or transaction number
    
    -- Amounts
    amount DECIMAL(20,2) NOT NULL DEFAULT 0.00,          -- Primary transaction amount
    amount_transfer DECIMAL(20,2) NOT NULL DEFAULT 0.00, -- Amount specifically for transfer
    lumpsum_amount DECIMAL(20,2) NOT NULL DEFAULT 0.00,  -- Lumpsum benefit component
    annuity_amount DECIMAL(20,2) NOT NULL DEFAULT 0.00,  -- Annuity benefit component
    bank_fee DECIMAL(20,2) NOT NULL DEFAULT 0.00,        -- Fee charged by the bank
    dplk_income DECIMAL(20,2) NOT NULL DEFAULT 0.00,     -- Income component for the DPLK
    
    -- Bank Information
    bank_name VARCHAR(255),         -- Name of the involved bank
    bank_account_number VARCHAR(255), -- Target bank account number
    bank_account_name VARCHAR(255),   -- Name registered on the bank account
    
    -- Status and Processing
    status INT NOT NULL DEFAULT 0,  -- Processing status (e.g., 0: Pending, 1: Processed)
    fee_burden_type INT,            -- Who bears the fee (1: DPLK, 2: Participant)
    realization_date TIMESTAMP WITH TIME ZONE, -- Date the transaction was finalized/realized
    verified_at TIMESTAMP WITH TIME ZONE,      -- Timestamp when the record was verified
    notes TEXT,                     -- Additional transaction notes or remarks
    
    -- Audit fields
    created_by UUID,    -- User who created this record
    updated_by UUID,    -- User who last updated this record
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- Record creation timestamp
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- Record update timestamp
    deleted_at TIMESTAMP WITH TIME ZONE -- Soft deletion timestamp
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_app_giro_reconciliation_details_giro_id ON app_giro_reconciliation_details(giro_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_app_giro_reconciliation_details_batch_id ON app_giro_reconciliation_details(batch_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_app_giro_reconciliation_details_giro_number ON app_giro_reconciliation_details(giro_number) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_app_giro_reconciliation_details_status ON app_giro_reconciliation_details(status) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_app_giro_reconciliation_details_transaction_date ON app_giro_reconciliation_details(transaction_date) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_app_giro_reconciliation_details_deleted_at ON app_giro_reconciliation_details(deleted_at);

-- Create app_batchings table (previously app_batchings)
-- Stores batch processing operations for membership and financial updates
CREATE TABLE IF NOT EXISTS app_batchings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- Primary key using UUID
    
    -- References
    membership_batching_id UUID REFERENCES mst_membership_batchings(id), -- Link to membership batching category
    batching_purpose_id UUID REFERENCES mst_batching_purposes(id),      -- Link to the purpose of this batch
    batching_detail_id UUID REFERENCES mst_batching_details(id),        -- Link to specific batching details
    
    -- Batch Identification
    batching_code VARCHAR(255) NOT NULL, -- Unique code/identifier for the batch
    description TEXT,                -- General description or notes about the batch
    
    -- Lifecycle and Status
    is_active BOOLEAN NOT NULL DEFAULT true, -- Whether the batch is currently open/active
    investpro_stored_status VARCHAR(255) NOT NULL DEFAULT 'not_stored', -- Status in target system (e.g., stored, failed)
    investpro_stored_message TEXT,   -- Message or error details from the target system
    
    -- Dates
    nab_date DATE,                   -- Net Asset Value calculation date
    cash_date DATE,                  -- Cash processing/realization date
    aum_date DATE,                   -- Assets Under Management reporting date
    transaction_date DATE,           -- Core transaction or accounting date
    
    -- Closing metadata
    closing_at TIMESTAMP WITH TIME ZONE, -- Timestamp when the batch was closed
    closing_by UUID REFERENCES sys_users(id), -- User who closed the batch
    
    -- Audit fields
    created_by UUID REFERENCES sys_users(id), -- User who created this record
    updated_by UUID REFERENCES sys_users(id), -- User who last updated this record
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- Record creation timestamp
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- Record update timestamp
    deleted_at TIMESTAMP WITH TIME ZONE -- Soft deletion timestamp
);

-- Create indexes
CREATE UNIQUE INDEX IF NOT EXISTS idx_app_batchings_code ON app_batchings(batching_code) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_app_batchings_membership_id ON app_batchings(membership_batching_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_app_batchings_purpose_id ON app_batchings(batching_purpose_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_app_batchings_is_active ON app_batchings(is_active) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_app_batchings_transaction_date ON app_batchings(transaction_date) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_app_batchings_deleted_at ON app_batchings(deleted_at);

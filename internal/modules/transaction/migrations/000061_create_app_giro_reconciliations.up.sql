-- Create app_giro_reconciliations table (previously reconciliation_giro)
-- Stores data for reconciling giro transactions with system batches
CREATE TABLE IF NOT EXISTS app_giro_reconciliations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- Primary key using UUID
    
    -- References and Identifiers
    batch_id UUID,                  -- Reference to the system batch (previously sb_id)
    legacy_batch_id BIGINT,         -- Legacy reference for data migration consistency
    giro_number VARCHAR(255) NOT NULL, -- Unique identifier of the giro document
    
    -- Giro Details
    giro_description TEXT,          -- Detailed description of the giro transaction
    giro_date DATE NOT NULL,        -- The date recorded on the giro
    
    -- Classification and Status
    payment_type INT NOT NULL,      -- Type of payment (e.g., 1: Transfer, 2: Cash, as per system mapping)
    status INT NOT NULL DEFAULT 0,  -- Reconciliation status (e.g., 0: Pending, 1: Matched)
    
    -- Audit fields
    created_by UUID,    -- User who created this record
    updated_by UUID,    -- User who last updated this record
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- Record creation timestamp
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- Record update timestamp
    deleted_at TIMESTAMP WITH TIME ZONE -- Soft deletion timestamp
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_app_giro_reconciliations_batch_id ON app_giro_reconciliations(batch_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_app_giro_reconciliations_giro_number ON app_giro_reconciliations(giro_number) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_app_giro_reconciliations_giro_date ON app_giro_reconciliations(giro_date) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_app_giro_reconciliations_status ON app_giro_reconciliations(status) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_app_giro_reconciliations_deleted_at ON app_giro_reconciliations(deleted_at);

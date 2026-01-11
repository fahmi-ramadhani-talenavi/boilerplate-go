-- Create mst_tax_brackets_ministry_regulation_16 table (previously mst_tax_brackets_ministry_regulation_16)
-- PMK 16 = Peraturan Menteri Keuangan No. 16 (Ministry of Finance Regulation No. 16)
-- Stores configurations for tax brackets based on Ministry Regulation 16
CREATE TABLE IF NOT EXISTS mst_tax_brackets_ministry_regulation_16 (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- Primary key using UUID
    
    -- Income and Rate Configuration
    minimum_income DECIMAL(15,2),    -- Minimum income threshold for this bracket
    maximum_income DECIMAL(15,2),    -- Maximum income threshold for this bracket
    tax_rate DECIMAL(5,2) NOT NULL,  -- Tax rate percentage (e.g., 5.00)
    
    -- Tax Rate Details
    effective_tax_rate_category VARCHAR(255) DEFAULT 'TER A', -- Category (A, B, C)
    
    -- Logic and Status
    effective_date DATE,             -- Date when this configuration becomes effective
    logic_operator VARCHAR(255) DEFAULT '=', -- Operator for logic matching
    is_active BOOLEAN NOT NULL DEFAULT true, -- Whether currently active
    
    -- Audit fields
    created_by UUID,    -- User who created this record
    updated_by UUID,    -- User who last updated this record
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE -- Soft deletion timestamp
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_mst_tax_brackets_mr16_effective_date ON mst_tax_brackets_ministry_regulation_16(effective_date) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_mst_tax_brackets_mr16_is_active ON mst_tax_brackets_ministry_regulation_16(is_active) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_mst_tax_brackets_mr16_category ON mst_tax_brackets_ministry_regulation_16(effective_tax_rate_category) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_mst_tax_brackets_mr16_deleted_at ON mst_tax_brackets_ministry_regulation_16(deleted_at);

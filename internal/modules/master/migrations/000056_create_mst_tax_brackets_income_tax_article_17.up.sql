-- Create mst_tax_brackets_income_tax_article_17 table (previously mst_tax_brackets_income_tax_article_17)
-- PPh 17 = Pajak Penghasilan Pasal 17 (Income Tax Article 17)
-- Stores configurations for tax brackets based on Income Tax Law Article 17
CREATE TABLE IF NOT EXISTS mst_tax_brackets_income_tax_article_17 (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- Primary key using UUID
    
    -- Income and Rate Configuration
    minimum_income DECIMAL(15,2),    -- Minimum income threshold for this bracket
    maximum_income DECIMAL(15,2),    -- Maximum income threshold for this bracket
    tax_rate DECIMAL(5,2) NOT NULL,  -- Tax rate percentage (e.g., 5.00)
    
    -- Tax Rate Details
    effective_tax_rate_category VARCHAR(255), -- Category (if applicable)
    
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
CREATE INDEX IF NOT EXISTS idx_mst_tax_brackets_ita17_effective_date ON mst_tax_brackets_income_tax_article_17(effective_date) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_mst_tax_brackets_ita17_is_active ON mst_tax_brackets_income_tax_article_17(is_active) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_mst_tax_brackets_ita17_deleted_at ON mst_tax_brackets_income_tax_article_17(deleted_at);

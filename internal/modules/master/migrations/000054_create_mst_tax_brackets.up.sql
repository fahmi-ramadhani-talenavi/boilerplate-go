-- Create mst_tax_brackets table (previously tax_brackets)
-- Stores configurations for income tax brackets and Effective Tax Rate (TER) categories
CREATE TABLE IF NOT EXISTS mst_tax_brackets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- Primary key using UUID
    
    -- Income and Rate Configuration
    min_income DECIMAL(15,2),       -- Minimum income threshold for this bracket
    max_income DECIMAL(15,2),       -- Maximum income threshold for this bracket
    tax_rate DECIMAL(5,2) NOT NULL, -- Tax rate percentage for this bracket (e.g., 5.00)
    
    -- TER Details
    effective_tax_rate_category VARCHAR(255) DEFAULT 'TER A', -- TER Category (e.g., TER A, TER B, TER C)
    
    -- Logic and Status
    effective_date DATE,            -- Date when this tax bracket configuration becomes effective
    logic_operator VARCHAR(255) DEFAULT '=', -- Operator used for logic matching (e.g., =, <, >)
    is_active BOOLEAN NOT NULL DEFAULT true, -- Whether this tax bracket is currently used
    
    -- Audit fields
    created_by UUID,    -- User who created this record
    updated_by UUID,    -- User who last updated this record
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- Record creation timestamp
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- Record update timestamp
    deleted_at TIMESTAMP WITH TIME ZONE -- Soft deletion timestamp
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_mst_tax_brackets_effective_date ON mst_tax_brackets(effective_date) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_mst_tax_brackets_is_active ON mst_tax_brackets(is_active) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_mst_tax_brackets_ter_category ON mst_tax_brackets(effective_tax_rate_category) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_mst_tax_brackets_deleted_at ON mst_tax_brackets(deleted_at);

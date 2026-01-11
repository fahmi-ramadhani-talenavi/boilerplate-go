-- Create mst_effective_tax_rate_categories table (previously mst_effective_tax_rate_categories / ter_category)
-- Stores reference data for Effective Tax Rate (TER) categories
CREATE TABLE IF NOT EXISTS mst_effective_tax_rate_categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- Primary key using UUID
    
    -- Main data
    description VARCHAR(255) NOT NULL, -- Name or description of the category (e.g., Category A, Category B)
    
    -- Audit fields
    created_by UUID,    -- User who created this record
    updated_by UUID,    -- User who last updated this record
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- Record creation timestamp
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- Record update timestamp
    deleted_at TIMESTAMP WITH TIME ZONE -- Soft deletion timestamp
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_mst_effective_tax_rate_categories_description ON mst_effective_tax_rate_categories(description) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_mst_effective_tax_rate_categories_deleted_at ON mst_effective_tax_rate_categories(deleted_at);

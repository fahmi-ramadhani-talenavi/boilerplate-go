-- Create mst_effective_tax_rate_layers table (previously mst_effective_tax_rate_layers / ter_layer)
-- Stores configurations for tax layers within each effective tax rate category
CREATE TABLE IF NOT EXISTS mst_effective_tax_rate_layers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- Primary key using UUID
    
    -- References
    effective_tax_rate_category_id UUID REFERENCES mst_effective_tax_rate_categories(id) ON DELETE CASCADE, -- Link to the category
    legacy_category_id BIGINT,      -- Legacy reference to the original category ID
    
    -- Layer Configuration
    logic_operator VARCHAR(10),     -- Logic operator for matching (e.g., =, <, >)
    minimum_amount DECIMAL(20,2) NOT NULL DEFAULT 0.00, -- Minimum income for this layer
    maximum_amount DECIMAL(20,2) NOT NULL DEFAULT 0.00, -- Maximum income for this layer
    tax_rate DECIMAL(10,2) NOT NULL DEFAULT 0.00,    -- Tax rate percentage (e.g., 5.00)
    
    -- Status and Dates
    effective_date DATE NOT NULL,   -- Date when this layer configuration becomes effective
    is_active BOOLEAN NOT NULL DEFAULT true, -- Whether this layer is currently active
    
    -- Audit fields
    created_by UUID,    -- User who created this record
    updated_by UUID,    -- User who last updated this record
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- Record creation timestamp
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- Record update timestamp
    deleted_at TIMESTAMP WITH TIME ZONE -- Soft deletion timestamp
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_mst_effective_tax_rate_layers_category_id ON mst_effective_tax_rate_layers(effective_tax_rate_category_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_mst_effective_tax_rate_layers_effective_date ON mst_effective_tax_rate_layers(effective_date) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_mst_effective_tax_rate_layers_is_active ON mst_effective_tax_rate_layers(is_active) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_mst_effective_tax_rate_layers_deleted_at ON mst_effective_tax_rate_layers(deleted_at);

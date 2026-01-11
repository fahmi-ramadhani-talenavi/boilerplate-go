-- Create mst_asset_incomes table (previously sys_penghasilan_asets)
-- Stores asset income reference data
CREATE TABLE IF NOT EXISTS mst_asset_incomes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- Primary key using UUID
    
    -- Main data
    description VARCHAR(255) NOT NULL, -- Source of asset income (e.g., Gaji, Usaha, Investasi)
    is_personal BOOLEAN NOT NULL DEFAULT false, -- Whether the asset is personal or corporate-owned
    
    -- Audit fields
    created_by UUID,    -- User who created this record
    updated_by UUID,    -- User who last updated this record
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- Record creation timestamp
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- Record update timestamp
    deleted_at TIMESTAMP WITH TIME ZONE -- Soft deletion timestamp
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_mst_asset_incomes_description ON mst_asset_incomes(description) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_mst_asset_incomes_deleted_at ON mst_asset_incomes(deleted_at);

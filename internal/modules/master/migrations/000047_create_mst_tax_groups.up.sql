-- Create mst_tax_groups table (previously sys_tax_group)
-- Stores tax group reference data including PTKP and TER configurations
CREATE TABLE IF NOT EXISTS mst_tax_groups (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- Primary key using UUID
    
    -- Main data
    name VARCHAR(255) NOT NULL,      -- Tax group name or category
    tax_exempt_income_code VARCHAR(255), -- PTKP (Penghasilan Tidak Kena Pajak) code as per tax regulation
    effective_tax_rate_code VARCHAR(255),  -- TER (Tarif Efektif Rata-rata) code for income tax calculation
    is_tax_id_combined BOOLEAN NOT NULL DEFAULT false, -- Whether the tax ID (NPWP) is combined for spouses
    effective_date DATE,            -- Date when this tax grouping becomes effective
    is_active BOOLEAN NOT NULL DEFAULT true, -- Whether the tax group is active for selection
    
    -- Audit fields
    created_by UUID,    -- User who created this record
    updated_by UUID,    -- User who last updated this record
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- Record creation timestamp
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- Record update timestamp
    deleted_at TIMESTAMP WITH TIME ZONE -- Soft deletion timestamp
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_mst_tax_groups_name ON mst_tax_groups(name) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_mst_tax_groups_effective_date ON mst_tax_groups(effective_date) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_mst_tax_groups_deleted_at ON mst_tax_groups(deleted_at);

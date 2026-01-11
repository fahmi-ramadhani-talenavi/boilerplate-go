-- Create mst_statuses table (previously sys_statuses)
-- Stores system status reference data with visual configuration for DPLK and Member portals
CREATE TABLE IF NOT EXISTS mst_statuses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- Primary key using UUID
    
    -- Pension Fund Portal Configuration (Portal DPLK)
    pension_fund_description VARCHAR(255) NOT NULL, -- Status description for Pension Fund portal
    pension_fund_text_color VARCHAR(25) NOT NULL DEFAULT '#D7EAF8', -- Text hex color for status tag
    pension_fund_background_color VARCHAR(25) NOT NULL DEFAULT '#3995DB', -- Background hex color for status tag
    
    -- Member Portal Configuration (Portal Peserta)
    member_description VARCHAR(255) NOT NULL, -- Status description for Member portal
    member_text_color VARCHAR(25) NOT NULL DEFAULT '#D7EAF8', -- Text hex color for status tag
    member_background_color VARCHAR(25) NOT NULL DEFAULT '#3995DB', -- Background hex color for status tag
    
    -- Audit fields
    created_by UUID,    -- User who created this record
    updated_by UUID,    -- User who last updated this record
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- Record creation timestamp
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- Record update timestamp
    deleted_at TIMESTAMP WITH TIME ZONE -- Soft deletion timestamp
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_mst_statuses_pension_fund_desc ON mst_statuses(pension_fund_description) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_mst_statuses_member_desc ON mst_statuses(member_description) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_mst_statuses_deleted_at ON mst_statuses(deleted_at);

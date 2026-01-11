-- Create mst_status_progress table (previously sys_statuses_progress)
-- Stores reference data for status progress with visual configuration for DPLK and Member portals
CREATE TABLE IF NOT EXISTS mst_status_progress (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- Primary key using UUID
    
    -- Pension Fund Portal Configuration (Portal DPLK)
    pension_fund_description VARCHAR(255) NOT NULL, -- Progress description for Pension Fund portal
    pension_fund_text_color VARCHAR(25) NOT NULL DEFAULT '#D7EAF8', -- Text hex color for progress tag
    pension_fund_background_color VARCHAR(25) NOT NULL DEFAULT '#3995DB', -- Background hex color for progress tag
    
    -- Member Portal Configuration (Portal Peserta)
    member_description VARCHAR(255) NOT NULL, -- Progress description for Member portal
    member_text_color VARCHAR(25) NOT NULL DEFAULT '#D7EAF8', -- Text hex color for progress tag
    member_background_color VARCHAR(25) NOT NULL DEFAULT '#3995DB', -- Background hex color for progress tag
    
    -- Audit fields
    created_by UUID,    -- User who created this record
    updated_by UUID,    -- User who last updated this record
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- Record creation timestamp
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- Record update timestamp
    deleted_at TIMESTAMP WITH TIME ZONE -- Soft deletion timestamp
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_mst_status_progress_pension_fund_desc ON mst_status_progress(pension_fund_description) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_mst_status_progress_member_desc ON mst_status_progress(member_description) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_mst_status_progress_deleted_at ON mst_status_progress(deleted_at);

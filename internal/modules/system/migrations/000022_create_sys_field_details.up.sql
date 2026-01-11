-- Create sys_field_details table
-- Stores metadata and configuration for dynamic fields in various tables
CREATE TABLE IF NOT EXISTS sys_field_details (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- Primary key using UUID
    
    -- Target information
    table_name VARCHAR(255) NOT NULL, -- Target database table for the dynamic field
    field_name VARCHAR(255),          -- Source column name in the target table
    
    -- UI Configuration
    label VARCHAR(255),               -- Human-readable display label for the UI
    field_type VARCHAR(255),          -- UI component type (e.g., text, select, date)
    
    -- Relation Metadata
    relation_table VARCHAR(255),      -- Foreign table link for select options
    relation_display_column VARCHAR(255), -- Column to display from the related table (formerly sfd_desc_relation)
    
    -- Behavior Flags
    is_multi_select BOOLEAN DEFAULT false, -- Whether multiple options can be selected (formerly sfd_multi)
    flag SMALLINT DEFAULT 0,              -- Internal behavior or feature flag
    
    -- Row-based Configuration
    row_title VARCHAR(255),           -- Sub-label or grouped row title
    row_field_name VARCHAR(255),       -- Reference name for row grouping
    
    -- Audit fields
    created_by UUID,    -- User who created this record
    updated_by UUID,    -- User who last updated this record
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- Record creation timestamp
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- Record update timestamp
    deleted_at TIMESTAMP WITH TIME ZONE -- Soft deletion timestamp
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_sys_field_details_table_name ON sys_field_details(table_name) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_sys_field_details_field_name ON sys_field_details(field_name) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_sys_field_details_deleted_at ON sys_field_details(deleted_at);

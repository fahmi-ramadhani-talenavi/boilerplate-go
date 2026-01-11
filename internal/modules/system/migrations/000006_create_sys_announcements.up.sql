-- Create sys_announcements table
-- Stores system announcements for users
CREATE TABLE IF NOT EXISTS sys_announcements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- Primary key using UUID
    
    -- Main data
    title VARCHAR(255) NOT NULL,    -- Announcement title
    is_active BOOLEAN DEFAULT true, -- Whether the announcement is active/published
    effective_date DATE,            -- Date when the announcement becomes visible
    description TEXT,               -- Full body content of the announcement
    
    -- Image attachment
    image_description TEXT, -- Caption or description for the attached image
    image_path TEXT,        -- Storage path for the announcement image
    
    -- File attachment
    attachment_description TEXT, -- Caption or description for the attached document
    attachment_path TEXT,        -- Storage path for the attached document
    
    -- Audit fields
    created_by UUID,    -- User who created this record
    updated_by UUID,    -- User who last updated this record
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- Record creation timestamp
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- Record update timestamp
    deleted_at TIMESTAMP WITH TIME ZONE -- Soft deletion timestamp
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_sys_announcements_is_active ON sys_announcements(is_active) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_sys_announcements_effective_date ON sys_announcements(effective_date) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_sys_announcements_deleted_at ON sys_announcements(deleted_at);

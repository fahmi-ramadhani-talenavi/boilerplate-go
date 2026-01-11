-- Create sys_roles table (previously sys_roles)
-- Stores system user roles
CREATE TABLE IF NOT EXISTS sys_roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- Primary key using UUID
    
    -- Main data
    description VARCHAR(255),       -- System user role name (e.g., Administrator, Staff)
    is_active BOOLEAN NOT NULL DEFAULT true, -- Whether the role is currently active for assignment
    
    -- Audit fields
    created_by UUID,    -- User who created this record
    updated_by UUID,    -- User who last updated this record
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- Record creation timestamp
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- Record update timestamp
    deleted_at TIMESTAMP WITH TIME ZONE -- Soft deletion timestamp
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_sys_roles_description ON sys_roles(description) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_sys_roles_is_active ON sys_roles(is_active) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_sys_roles_deleted_at ON sys_roles(deleted_at);

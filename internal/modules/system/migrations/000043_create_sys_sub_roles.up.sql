-- Create sys_sub_roles table (previously sys_sub_roles)
-- Stores system user sub-roles linked to sys_roles
CREATE TABLE IF NOT EXISTS sys_sub_roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- Primary key using UUID
    role_id UUID NOT NULL REFERENCES sys_roles(id) ON DELETE CASCADE, -- Link to parent role
    
    -- Main data
    description VARCHAR(100) NOT NULL, -- Sub-role name (e.g., Maker, Checker, Approver)
    is_active BOOLEAN NOT NULL DEFAULT true, -- Whether the sub-role is active for assignment
    
    -- Audit fields
    created_by UUID,    -- User who created this record
    updated_by UUID,    -- User who last updated this record
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- Record creation timestamp
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- Record update timestamp
    deleted_at TIMESTAMP WITH TIME ZONE -- Soft deletion timestamp
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_sys_sub_roles_role_id ON sys_sub_roles(role_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_sys_sub_roles_description ON sys_sub_roles(description) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_sys_sub_roles_deleted_at ON sys_sub_roles(deleted_at);

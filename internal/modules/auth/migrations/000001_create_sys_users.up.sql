-- Create sys_users table (previously users)
-- Main table for system users with organizational mapping
CREATE TABLE IF NOT EXISTS sys_users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- Primary key using UUID
    
    -- Authentication
    email VARCHAR(255) NOT NULL UNIQUE, -- User's unique email address for authentication
    password VARCHAR(255) NOT NULL,    -- Bcrypt hashed password
    
    -- Profile Information
    full_name VARCHAR(255) NOT NULL,   -- User's full name
    employee_number VARCHAR(100),      -- Employee identification number (PNIP)
    
    -- Organizational Mapping (Foreign keys added in later migrations)
    branch_id UUID,        -- Link to mst_branches (Cabang)
    division_id UUID,      -- Link to divisions
    division_name VARCHAR(100), -- Legacy division description
    
    -- Roles and Access
    role_id UUID,          -- Primary role identifier (Link to sys_roles)
    sub_role_id UUID,      -- Secondary role identifier (Link to sys_sub_roles)
    
    -- Extended Profile IDs
    profile_id VARCHAR(30),        -- Legacy profile identifier
    company_profile_id VARCHAR(30), -- Legacy company profile identifier
    
    -- Status
    is_active BOOLEAN NOT NULL DEFAULT true, -- Whether the account is active
    
    -- Audit fields
    created_by UUID,    -- User who created this record
    updated_by UUID,    -- User who last updated this record
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- Timestamp when record was created
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- Timestamp when record was last updated
    deleted_at TIMESTAMP WITH TIME ZONE -- Timestamp for soft deletion
);

-- Create indexes
CREATE UNIQUE INDEX IF NOT EXISTS idx_sys_users_email ON sys_users(email) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_sys_users_employee_number ON sys_users(employee_number) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_sys_users_branch_id ON sys_users(branch_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_sys_users_is_active ON sys_users(is_active) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_sys_users_deleted_at ON sys_users(deleted_at);

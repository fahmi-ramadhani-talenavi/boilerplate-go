-- Create sys_sub_menus table (previously sub_menus)
-- Stores the application's sub-navigation hierarchy and display properties
CREATE TABLE IF NOT EXISTS sys_sub_menus (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- Primary key using UUID
    
    -- Hierarchy and Reference
    menu_id UUID,                   -- Reference to parent menu (sys_menus) - No FK yet as sys_menus doesn't exist
    parent_id UUID REFERENCES sys_sub_menus(id) ON DELETE CASCADE, -- Self-reference for nested sub-menus
    
    -- Main data
    key VARCHAR(255) NOT NULL,      -- Unique programmatic key for identification
    label VARCHAR(255) NOT NULL,    -- Display label for the sub-menu item
    icon VARCHAR(255),              -- Icon identifier/name for display
    router_link VARCHAR(255),       -- Internal application route/link
    path VARCHAR(255),              -- Physical path or additional routing info
    is_visible BOOLEAN NOT NULL DEFAULT true, -- Whether this sub-menu is visible in the UI
    
    -- Audit fields
    created_by UUID,    -- User who created this record
    updated_by UUID,    -- User who last updated this record
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- Record creation timestamp
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- Record update timestamp
    deleted_at TIMESTAMP WITH TIME ZONE -- Soft deletion timestamp
);

-- Create indexes
CREATE UNIQUE INDEX IF NOT EXISTS idx_sys_sub_menus_key ON sys_sub_menus(key) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_sys_sub_menus_menu_id ON sys_sub_menus(menu_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_sys_sub_menus_parent_id ON sys_sub_menus(parent_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_sys_sub_menus_is_visible ON sys_sub_menus(is_visible) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_sys_sub_menus_deleted_at ON sys_sub_menus(deleted_at);

-- Create mst_risk_profiles table (previously sys_risk_profile)
-- Stores risk profile reference data
CREATE TABLE IF NOT EXISTS mst_risk_profiles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- Primary key using UUID
    
    -- Main data
    name VARCHAR(255) NOT NULL,      -- Risk profile category (e.g., Konservatif, Moderat, Agresif)
    investpro_portfolio_id BIGINT, -- Link to core investment system portfolio ID
    
    -- Audit fields
    created_by UUID,    -- User who created this record
    updated_by UUID,    -- User who last updated this record
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- Record creation timestamp
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- Record update timestamp
    deleted_at TIMESTAMP WITH TIME ZONE -- Soft deletion timestamp
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_mst_risk_profiles_name ON mst_risk_profiles(name) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_mst_risk_profiles_deleted_at ON mst_risk_profiles(deleted_at);

-- Create mst_document_types table (previously sys_dokumen_types)
-- Stores reference data for various document types required in business processes
CREATE TABLE IF NOT EXISTS mst_document_types (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- Primary key using UUID
    
    -- Main data
    name VARCHAR(255) NOT NULL,    -- Document type name (e.g., KTP, KK, Form A)
    description VARCHAR(255),      -- Purpose of the document type
    sort_order INT NOT NULL DEFAULT 0, -- Display priority
    
    -- Classification
    -- Category: 1 = Claim, 2 = Withdrawal, 3 = Transfer (Pengalihan)
    category SMALLINT,             -- Business category classification
    claim_type VARCHAR(255),       -- Specific claim sub-type or role restriction (formerly role_claim)
    
    -- Audit fields
    created_by UUID,    -- User who created this record
    updated_by UUID,    -- User who last updated this record
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- Record creation timestamp
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- Record update timestamp
    deleted_at TIMESTAMP WITH TIME ZONE -- Soft deletion timestamp
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_mst_document_types_category ON mst_document_types(category) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_mst_document_types_sort_order ON mst_document_types(sort_order) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_mst_document_types_deleted_at ON mst_document_types(deleted_at);

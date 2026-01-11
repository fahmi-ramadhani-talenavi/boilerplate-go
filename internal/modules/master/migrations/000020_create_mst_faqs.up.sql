-- Create mst_faqs table (previously sys_faq)
-- Stores frequently asked questions and answers
CREATE TABLE IF NOT EXISTS mst_faqs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- Primary key using UUID
    
    -- Main data
    subject VARCHAR(255) NOT NULL, -- FAQ question title or subject
    content TEXT NOT NULL,          -- FAQ detailed answer or content
    sort_order INT NOT NULL DEFAULT 0, -- Display priority
    
    -- Audit fields
    created_by UUID,    -- User who created this record
    updated_by UUID,    -- User who last updated this record
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- Record creation timestamp
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- Record update timestamp
    deleted_at TIMESTAMP WITH TIME ZONE -- Soft deletion timestamp
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_mst_faqs_subject ON mst_faqs(subject) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_mst_faqs_sort_order ON mst_faqs(sort_order) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_mst_faqs_deleted_at ON mst_faqs(deleted_at);

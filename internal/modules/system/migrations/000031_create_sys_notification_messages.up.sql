-- Create sys_notification_messages table (previously sys_notif_messages)
-- Stores notification message templates or cataloged messages
CREATE TABLE IF NOT EXISTS sys_notification_messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- Primary key using UUID
    
    -- Main data
    content TEXT NOT NULL,      -- Full body content of the notification message
    message_type INT,           -- Target audience or category: 1=DPLK Portal, 2=Member Portal (Peserta)
    is_batch BOOLEAN NOT NULL DEFAULT false, -- Whether the message is part of a batch notification
    
    -- Audit fields
    created_by UUID,    -- User who created this record
    updated_by UUID,    -- User who last updated this record
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- Record creation timestamp
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- Record update timestamp
    deleted_at TIMESTAMP WITH TIME ZONE -- Soft deletion timestamp
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_sys_notification_messages_message_type ON sys_notification_messages(message_type) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_sys_notification_messages_is_batch ON sys_notification_messages(is_batch) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_sys_notification_messages_deleted_at ON sys_notification_messages(deleted_at);

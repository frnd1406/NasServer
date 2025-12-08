-- Migration 004: Honeyfile Security System
-- Intrusion detection via trap files that trigger vault lockdown

CREATE TABLE IF NOT EXISTS honeyfiles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    file_path VARCHAR(512) NOT NULL UNIQUE,
    file_type VARCHAR(50) NOT NULL DEFAULT 'general',  -- 'finance', 'it', 'private', 'general'
    trigger_count INT DEFAULT 0,
    last_triggered_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_by UUID REFERENCES users(id),
    
    CONSTRAINT honeyfiles_type_check CHECK (file_type IN ('finance', 'it', 'private', 'general'))
);

-- Index for fast path lookups (used on every download/preview request)
CREATE INDEX idx_honeyfiles_path ON honeyfiles(file_path);

-- Log initialization
DO $$
BEGIN
    RAISE NOTICE 'Migration 004: Honeyfiles table created';
END $$;

-- migrations/0002_create_chats.sql

CREATE TABLE IF NOT EXISTS chats (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    type TEXT NOT NULL CHECK (type IN ('private', 'group')), 
    created_at TIMESTAMPTZ DEFAULT NOW()
);

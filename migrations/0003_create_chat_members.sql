-- migrations/0003_create_chat_members.sql

CREATE TABLE IF NOT EXISTS chat_members (
    chat_id UUID NOT NULL,
    user_id UUID NOT NULL,
    joined_at TIMESTAMPTZ DEFAULT NOW(),

    PRIMARY KEY (chat_id, user_id),
    CONSTRAINT fk_chat FOREIGN KEY (chat_id) REFERENCES chats(id) ON DELETE CASCADE,
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

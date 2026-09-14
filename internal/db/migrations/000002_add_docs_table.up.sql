CREATE TABLE IF NOT EXISTS documents (
    id UUID PRIMARY KEY,
    owner_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    mime VARCHAR(100),
    has_file BOOLEAN NOT NULL DEFAULT FALSE,
    is_public BOOLEAN NOT NULL DEFAULT FALSE,
    file_path VARCHAR(255),
    json_data JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,

    CONSTRAINT fk_documents_owner FOREIGN KEY (owner_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS documents_users (
    user_id UUID NOT NULL,
    doc_id UUID NOT NULL,

    PRIMARY KEY (doc_id, user_id),

    CONSTRAINT fk_documents_users_doc FOREIGN KEY (doc_id) REFERENCES documents(id) ON DELETE CASCADE,
    CONSTRAINT fk_documents_users_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
)
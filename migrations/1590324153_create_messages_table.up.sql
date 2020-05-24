CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE messages
(
    id         UUID DEFAULT uuid_generate_v4 () PRIMARY KEY,
    name       varchar NOT NULL,
    created_at TIMESTAMP DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMP DEFAULT NOW() NOT NULL
);

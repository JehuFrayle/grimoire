-- Enable UUID extension (if not enabled)
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Users table
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    role VARCHAR(20) NOT NULL CHECK (role IN ('admin', 'user', 'guest')),
    active BOOLEAN NOT NULL DEFAULT TRUE
);

-- Profiles table (1:1 relationship with users)
CREATE TABLE IF NOT EXISTS profiles (
    user_id UUID PRIMARY KEY,
    first_name VARCHAR(50),
    last_name VARCHAR(50),
    bio TEXT,
    avatar_url TEXT,
    CONSTRAINT fk_profiles_user FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Links table (multiple links per user)
CREATE TABLE IF NOT EXISTS links (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL,
    url TEXT NOT NULL,
    title VARCHAR(100),
    icon TEXT,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    CONSTRAINT fk_links_user FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Notes table
CREATE TABLE IF NOT EXISTS notes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR(200) NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Tags table
CREATE TABLE IF NOT EXISTS tags (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(50) UNIQUE NOT NULL
);

-- Join table for many-to-many relationship between notes and tags
CREATE TABLE IF NOT EXISTS note_tags (
    note_id UUID NOT NULL,
    tag_id UUID NOT NULL,
    PRIMARY KEY(note_id, tag_id),
    CONSTRAINT fk_notetags_note FOREIGN KEY(note_id) REFERENCES notes(id) ON DELETE CASCADE,
    CONSTRAINT fk_notetags_tag FOREIGN KEY(tag_id) REFERENCES tags(id) ON DELETE CASCADE
);

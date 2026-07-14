CREATE TABLE IF NOT EXISTS social_links (
    id UUID PRIMARY KEY,
    platform VARCHAR(100) NOT NULL,
    url TEXT NOT NULL,
    icon VARCHAR(255),
    profile_id UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    CONSTRAINT fk_profiles_social_links FOREIGN KEY (profile_id) REFERENCES profiles(id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_profile_platform ON social_links(profile_id, platform);

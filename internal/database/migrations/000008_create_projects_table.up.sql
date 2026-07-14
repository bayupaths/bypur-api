CREATE TABLE IF NOT EXISTS projects (
    id UUID PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    content TEXT,
    image TEXT,
    tech_stack JSONB DEFAULT '[]'::jsonb NOT NULL,
    url TEXT,
    github TEXT,
    featured BOOLEAN DEFAULT FALSE NOT NULL,
    "order" INT DEFAULT 0 NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_projects_slug ON projects(slug);

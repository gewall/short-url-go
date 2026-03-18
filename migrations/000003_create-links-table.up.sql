CREATE TABLE IF NOT EXISTS links (
  id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id      UUID REFERENCES users(id) ON DELETE CASCADE,
  original_url TEXT NOT NULL,
  short_code   VARCHAR(10) UNIQUE NOT NULL,
  title        TEXT,
  is_active    BOOLEAN DEFAULT TRUE,
  expires_at   TIMESTAMPTZ,
  created_at   TIMESTAMPTZ DEFAULT NOW()
);

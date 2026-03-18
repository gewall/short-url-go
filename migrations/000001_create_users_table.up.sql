CREATE TABLE IF NOT EXISTS users (
  id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  username      TEXT UNIQUE NOT NULL,
  password   TEXT NOT NULL,         -- bcrypt
  created_at TIMESTAMPTZ DEFAULT NOW()
);

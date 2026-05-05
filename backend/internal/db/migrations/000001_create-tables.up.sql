CREATE TABLE transactions (
  id SERIAL PRIMARY KEY,
  amount VARCHAR(255) NOT NULL,
  reason TEXT,
  created_at TIMESTAMP DEFAULT NOW(),
  type VARCHAR(25),
  user_id    UUID REFERENCES users(id) ON DELETE CASCADE
);

INSERT INTO schema_migrations (version, dirty) VALUES (4, false);
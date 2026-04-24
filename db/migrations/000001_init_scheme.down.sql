DROP TRIGGER IF EXISTS trg_users_updated_at ON users;

DROP FUNCTION IF EXISTS update_updated_at();

DROP TABLE IF EXISTS refresh_tokens;
DROP TABLE IF EXISTS users;

DROP INDEX IF EXISTS idx_users_email_canonical_unique;
DROP TYPE IF EXISTS user_role;

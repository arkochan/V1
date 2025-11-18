-- Drop the trigger first
DROP TRIGGER IF EXISTS set_updated_at_trigger ON auth;

-- Drop the function
DROP FUNCTION IF EXISTS set_updated_at();

-- Drop the index
DROP INDEX IF EXISTS users_email_idx;

-- Drop the table
DROP TABLE IF EXISTS auth;

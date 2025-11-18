CREATE TABLE auth (
    id                 uuid DEFAULT uuidv7() PRIMARY KEY,
    email              text NOT NULL UNIQUE CHECK (email = lower(email)),
    password_hash      text NOT NULL,
    status             text NOT NULL DEFAULT 'active',  -- active | disabled | banned
    created_at         timestamptz NOT NULL DEFAULT now(),
    updated_at         timestamptz NOT NULL DEFAULT now(),
    deleted_at         timestamptz DEFAULT NULL
);

-- Ensure normalized email uniqueness (case-insensitive)
CREATE UNIQUE INDEX users_email_idx
    ON auth (email);

CREATE FUNCTION set_updated_at()
RETURNS trigger AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

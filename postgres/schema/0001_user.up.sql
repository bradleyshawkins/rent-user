BEGIN;

CREATE TABLE IF NOT EXISTS app_user_details(
    id UUID NOT NULL PRIMARY KEY,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    email_address TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    modified_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS app_user_details_email_address ON app_user_details(email_address);

CREATE TABLE IF NOT EXISTS app_user_credentials (
    id UUID NOT NULL PRIMARY KEY,
    username TEXT NOT NULL,
    password TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    modified_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS app_user_credentials_username ON app_user_credentials(username);

CREATE TABLE IF NOT EXISTS app_user (
    id UUID NOT NULL PRIMARY KEY,
    status TEXT NOT NULL,
    app_user_credentials_id UUID NOT NULL REFERENCES app_user_credentials(id),
    app_user_details_id UUID NOT NULL REFERENCES app_user_details(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    modified_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS app_user_status_idx ON app_user(status);
CREATE INDEX IF NOT EXISTS app_user_credentials_id ON app_user(app_user_credentials_id);

CREATE OR REPLACE FUNCTION update_modified_at_column()
    RETURNS TRIGGER AS $$
BEGIN
    NEW.modified_at = NOW();
    RETURN NEW;
END;
$$ language plpgsql;

CREATE TRIGGER update_app_user_modified_at_column BEFORE INSERT OR UPDATE ON app_user FOR EACH ROW EXECUTE FUNCTION update_modified_at_column();
CREATE TRIGGER update_app_user_details_modified_at_column BEFORE INSERT OR UPDATE ON app_user_details FOR EACH ROW EXECUTE FUNCTION update_modified_at_column();
CREATE TRIGGER update_app_user_credentials_modified_at_column BEFORE INSERT OR UPDATE ON app_user_credentials FOR EACH ROW EXECUTE FUNCTION update_modified_at_column();

COMMIT;

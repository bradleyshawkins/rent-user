BEGIN;

CREATE TABLE IF NOT EXISTS account(
    id UUID NOT NULL PRIMARY KEY,
    status TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    modified_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER update_account_modified_at_column BEFORE INSERT OR UPDATE ON account FOR EACH ROW EXECUTE FUNCTION update_modified_at_column();

CREATE TABLE IF NOT EXISTS membership (
    app_user_id UUID NOT NULL REFERENCES app_user(id),
    account_id UUID NOT NULL REFERENCES account(id),
    role TEXT NOT NULL,
    PRIMARY KEY (app_user_id, account_id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    modified_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER update_membership_modified_at_column BEFORE INSERT OR UPDATE ON membership FOR EACH ROW EXECUTE FUNCTION update_modified_at_column();

COMMIT;

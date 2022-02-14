BEGIN;

DROP TRIGGER IF EXISTS update_membership_modified_at_column ON membership;
DROP TRIGGER  IF EXISTS update_account_modified_at_column ON account;

DROP TABLE IF EXISTS membership;
DROP TABLE IF EXISTS account;

COMMIT;
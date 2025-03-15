CREATE TABLE IF NOT EXISTS file_to_account(
    account_id UUID NOT NULL,
    file_uuid VARCHAR(256) NOT NULL,
    FOREIGN KEY (account_id) REFERENCES account,
    FOREIGN KEY (file_uuid) REFERENCES file(file_uuid)
);

ALTER TABLE file_to_account DROP CONSTRAINT file_to_account_account_id_fkey;

-- Step 2: Add a new foreign key constraint with ON DELETE SET NULL
ALTER TABLE file_to_account
ADD CONSTRAINT file_to_account_account_id_fkey
FOREIGN KEY (account_id) REFERENCES account(account_id) ON DELETE SET NULL;

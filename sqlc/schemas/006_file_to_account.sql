CREATE TABLE IF NOT EXISTS file_to_account(
    account_id UUID NOT NULL,
    file_uuid VARCHAR(256) NOT NULL,
    FOREIGN KEY (account_id) REFERENCES account,
    FOREIGN KEY (file_uuid) REFERENCES file(file_uuid)
);
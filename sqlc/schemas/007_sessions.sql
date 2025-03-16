DO $$ BEGIN
    CREATE TYPE TOKEN_TYPE AS ENUM ('refresh', 'authentication');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

CREATE TABLE IF NOT EXISTS account_session(
    session_id VARCHAR(256) PRIMARY KEY,
    account_id UUID NOT NULL,
    start_date TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expire_date TIMESTAMP WITH TIME ZONE NOT NULL,
    token_type TOKEN_TYPE NOT NULL DEFAULT 'authentication',
    FOREIGN KEY (account_id) REFERENCES account,
    CHECK(expire_date > start_date)
);

CREATE OR REPLACE FUNCTION get_userid_by_session(
    IN session_id VARCHAR(256),
    OUT account_id UUID
) AS $$
    SELECT account_id
    FROM account_session
    WHERE session_id = $1
$$  LANGUAGE SQL;

-- Modifications
ALTER TABLE account_session DROP CONSTRAINT account_session_account_id_fkey;
ALTER TABLE account_session
ADD CONSTRAINT account_session_account_id_fkey
FOREIGN KEY (account_id) REFERENCES account(account_id) ON DELETE CASCADE;

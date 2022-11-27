CREATE TABLE IF NOT EXISTS report(
    file_uuid VARCHAR(256) NOT NULL,
    reason VARCHAR(1024) NOT NULL,
    FOREIGN KEY (file_uuid) REFERENCES file(file_uuid)
);

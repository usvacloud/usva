CREATE TABLE IF NOT EXISTS container(
    container_uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(256) NOT NULL,
    password VARCHAR(256) NOT NULL
);

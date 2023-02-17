CREATE TABLE IF NOT EXISTS file_to_container(
    file_to_container_uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    container_uuid UUID NOT NULL,
    file_uuid TEXT NOT NULL, 
    FOREIGN KEY (container_uuid) REFERENCES container,
    FOREIGN KEY (file_uuid) REFERENCES file
)
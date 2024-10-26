ALTER TABLE users
    RENAME COLUMN is_banned TO banned;

ALTER TABLE users
    RENAME COLUMN is_deleted TO deleted;

ALTER TABLE users
    DROP COLUMN is_admin;
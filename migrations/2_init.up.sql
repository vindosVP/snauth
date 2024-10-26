ALTER TABLE users
    RENAME COLUMN banned TO is_banned;

ALTER TABLE users
    RENAME COLUMN deleted TO is_deleted;

ALTER TABLE users
    ADD COLUMN is_admin boolean NOT NULL default false;

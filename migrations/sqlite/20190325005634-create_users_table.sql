
-- +migrate Up

CREATE TABLE users (
    user_id VARCHAR(15) PRIMARY KEY,
    screen_name VARCHAR(20) NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    password VARCHAR(97) NOT NULL,
    delete_flag VARCHAR(1) NOT NULL DEFAULT '0',
    is_admin VARCHAR(1) NOT NULL DEFAULT '0'
);

-- +migrate Down

DROP TABLE users;
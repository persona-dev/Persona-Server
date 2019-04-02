
-- +migrate Up

CREATE TABLE users (
    user_id VARCHAR(15) NOT NULL PRIMARY KEY,
    screen_name VARCHAR(20) NOT NULL,
    email VARCHAR(25) NOT NULL UNIQUE,
    delete_flag INT DEFAULT 0,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    password VARCHAR(150) NOT NULL,
    is_admin INT DEFAULT 0
);

-- +migrate Down

DROP TABLE users;
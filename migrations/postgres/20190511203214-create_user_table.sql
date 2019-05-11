
-- +migrate Up

CREATE TABLE users (
    user_id VARCHAR(15) NOT NULL PRIMARY KEY,
    screen_name VARCHAR(50) NOT NULL,
    email VARCHAR(254) NOT NULL UNIQUE,
    delete_flag BOOLEAN DEFAULT 'False',
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    password VARCHAR(150) NOT NULL,
    is_admin BOOLEAN DEFAULT 'False'
);

-- +migrate Down

DROP TABLE users;

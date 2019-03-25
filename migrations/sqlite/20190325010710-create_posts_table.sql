
-- +migrate Up

CREATE TABLE posts (
    post_id CHAR(26) PRIMARY KEY,
    user_id VARCHAR(15),
    delete_flag VARCHAR(1) DEFAULT '0',
    body VARCHAR(500) NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    FOREIGN KEY(user_id) REFERENCES users(user_id)
);

-- +migrate Down

DROP TABLE posts;
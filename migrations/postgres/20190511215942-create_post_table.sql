
-- +migrate Up

CREATE TABLE posts (
    post_id VARCHAR(26) NOT NULL PRIMARY KEY,
    user_id VARCHAR(15) NOT NULL,
    delete_flag BOOLEAN DEFAULT 'false' NOT NULL,
    body VARCHAR(500) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    FOREIGN KEY(user_id) REFERENCES users(user_id)
);

-- +migrate Down

DROP TABLE posts;

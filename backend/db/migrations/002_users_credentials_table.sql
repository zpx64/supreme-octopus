-- Write your migrate up statements here

CREATE TABLE IF NOT EXISTS users_credentials (
  user_id  INTEGER      NOT NULL REFERENCES users(user_id),
  email    VARCHAR(256) UNIQUE NOT NULL,
  password VARCHAR(256) NOT NULL,
  pow      VARCHAR(32)  NOT NULL
);

---- create above / drop below ----

DROP TABLE IF EXISTS users_credentials;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.

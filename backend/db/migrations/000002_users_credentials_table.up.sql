CREATE TABLE IF NOT EXISTS users_credentials (
  user_id  INTEGER      NOT NULL REFERENCES users(user_id),
  email    VARCHAR(256) UNIQUE NOT NULL,
  password VARCHAR(256) NOT NULL,
  pow      VARCHAR(32)  NOT NULL
);

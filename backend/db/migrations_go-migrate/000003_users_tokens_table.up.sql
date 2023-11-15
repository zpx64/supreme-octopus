CREATE TABLE IF NOT EXISTS users_tokens (
  user_id       INTEGER     NOT NULL REFERENCES users(user_id),
  device_id     VARCHAR(64) UNIQUE NOT NULL,
  refresh_token VARCHAR(36) UNIQUE NOT NULL,
  token_date    BIGINT      NOT NULL
);

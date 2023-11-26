-- Write your migrate up statements here

CREATE TABLE IF NOT EXISTS users_tokens (
  token_id      SERIAL       PRIMARY KEY,
  user_id       INTEGER      NOT NULL REFERENCES users(user_id),
  device_id     VARCHAR(20)  NOT NULL,
  refresh_token VARCHAR(36)  UNIQUE NOT NULL,
  user_agent    VARCHAR(256) NOT NULL,
  token_date    BIGINT       NOT NULL
);

---- create above / drop below ----

DROP TABLE IF EXISTS users_tokens;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.

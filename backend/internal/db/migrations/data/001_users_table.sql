-- Write your migrate up statements here

CREATE TABLE IF NOT EXISTS users (
  user_id       SERIAL       PRIMARY KEY,
  nickname      VARCHAR(256) UNIQUE NOT NULL,
  creation_date TIMESTAMPTZ  NOT NULL,
  name          VARCHAR(256),
  surname       VARCHAR(256)
);

---- create above / drop below ----

DROP TABLE IF EXISTS users;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.

CREATE TABLE IF NOT EXISTS users (
  user_id       SERIAL       PRIMARY KEY,
  nickname      VARCHAR(256) UNIQUE NOT NULL,
  creation_date TIMESTAMPTZ  NOT NULL,
  name          VARCHAR(256),
  surname       VARCHAR(256)
);

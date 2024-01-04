-- Write your migrate up statements here

CREATE TABLE IF NOT EXISTS users_posts (
  post_id                SERIAL      PRIMARY KEY,
  user_id                INTEGER     NOT NULL REFERENCES users(user_id),
  creation_date          TIMESTAMPTZ NOT NULL,
  post_type              INTEGER     NOT NULL,
  body                   TEXT        NOT NULL, -- https://stackoverflow.com/questions/7310558/postgresql-big-text-column-performance
  attachments            TEXT[]      NOT NULL,
  votes_amount           INTEGER     NOT NULL,
  comments_amount        INTEGER     NOT NULL,
  is_comments_disallowed BOOLEAN     NOT NULL
);

---- create above / drop below ----

DROP TABLE IF EXISTS users_posts;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.

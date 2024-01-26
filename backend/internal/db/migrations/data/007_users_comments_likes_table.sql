-- Write your migrate up statements here

CREATE TABLE IF NOT EXISTS users_comments_likes (
  like_id       SERIAL      PRIMARY KEY,
  user_id       INTEGER     NOT NULL REFERENCES users(user_id),
  comment_id    INTEGER     NOT NULL REFERENCES users_comments(comment_id),
  vote_type     INTEGER     NOT NULL,
  creation_date TIMESTAMPTZ NOT NULL,
  UNIQUE(user_id, comment_id)
);

---- create above / drop below ----

DROP TABLE IF EXISTS users_comments_likes;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.

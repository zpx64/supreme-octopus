-- Write your migrate up statements here

CREATE TABLE IF NOT EXISTS users_comments (
  comment_id    SERIAL      PRIMARY KEY,
  user_id       INTEGER     NOT NULL REFERENCES users(user_id),
  post_id       INTEGER     NOT NULL REFERENCES users_posts(post_id),
  body          TEXT        NOT NULL,
  attachments   TEXT[]      NOT NULL,
  creation_date TIMESTAMPTZ NOT NULL, 
  votes_amount  INTEGER     NOT NULL,
  reply_id      INTEGER     REFERENCES users_comments(comment_id) -- naive implementation of response threads
  -- I don't know if links should be used or something else.
  -- but it should be converted into some sort of graph
  -- and then send it to the frontend afterwards
)

---- create above / drop below ----

DROP TABLE IF EXISTS users_comments;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.

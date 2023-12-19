-- Write your migrate up statements here

CREATE TABLE IF NOT EXISTS users_comments (
  comment_id    SERIAL      PRIMARY KEY,
  user_id       INTEGER     NOT NULL REFERENCES users(user_id),
  post_id       INTEGER     NOT NULL REFERENCES users_posts(post_id),
  body          TEXT        NOT NULL,
  attachments   TEXT[]      NOT NULL,
  creation_date TIMESTAMPTZ NOT NULL, 
  votes_amount  INTEGER     NOT NULL,
  reply_id      INTEGER     REFERENCES users_comments(comment_id) -- naive implementation of reply threads
  -- i dont really know need we reference or something other kind of rellation
  -- but it should be converted to some kind of graph
  -- and after it sends to frontend
)

---- create above / drop below ----

DROP TABLE IF EXISTS users_comments;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.

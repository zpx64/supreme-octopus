-- Write your migrate up statements here

CREATE TABLE IF NOT EXISTS users_info ( 
  user_id      INTEGER      NOT NULL REFERENCES users(user_id),
  pronounse    VARCHAR(256) NOT NULL,
  gender       BOOLEAN      NOT NULL,
  about        TEXT         NOT NULL,
  about_line   VARCHAR(256) NOT NULL,
  social_links TEXT[]       NOT NULL
);

---- create above / drop below ----

DROP TABLE IF EXISTS users_info;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.

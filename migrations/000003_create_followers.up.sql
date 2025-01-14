
CREATE TABLE tag (
  id uuid PRIMARY KEY,
  slug varchar UNIQUE NOT NULL,
  level integer NOT NULL,
  created_at timestamp NOT NULL DEFAULT 'now()',
  updated_at timestamp NOT NULL DEFAULT 'now()'
);

CREATE TABLE user_tag (
  id uuid PRIMARY KEY,
  user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  tag_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  created_at timestamp NOT NULL DEFAULT 'now()',
  updated_at timestamp NOT NULL DEFAULT 'now()'
);

CREATE UNIQUE INDEX ON "user_tag" ("tag_id", "user_id");

CREATE TABLE follower (
  id uuid PRIMARY KEY,
  follower_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  following_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  created_at timestamp NOT NULL DEFAULT 'now()',
  updated_at timestamp NOT NULL DEFAULT 'now()'
);


CREATE UNIQUE INDEX ON "follower" ("follower_id", "following_id");

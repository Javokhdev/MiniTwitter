CREATE TYPE attachment_type AS ENUM (
  'photo',
  'video'
);

CREATE TYPE tweet_status AS ENUM (
    'draft',
    'published'
);


CREATE TABLE tweet (
  id uuid PRIMARY KEY,
  owner_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  content text NOT NULL,
  tags json,
  status tweet_status NOT NULL DEFAULT 'draft',
  created_at timestamp NOT NULL DEFAULT 'now()',
  updated_at timestamp NOT NULL DEFAULT 'now()'
);

CREATE TABLE tweet_attachment (
  id uuid PRIMARY KEY,
  tweet_id uuid NOT NULL REFERENCES tweet(id) ON DELETE CASCADE,
  filepath varchar NOT NULL,
  content_type attachment_type NOT NULL,
  created_at timestamp NOT NULL DEFAULT 'now()',
  updated_at timestamp NOT NULL DEFAULT 'now()'
);

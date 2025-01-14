CREATE TYPE user_type AS ENUM (
  'user',
  'admin'
);

CREATE TYPE user_role AS ENUM (
  'user',
  'admin',
  'superadmin'
);

CREATE TYPE gender AS ENUM (
  'male',
  'female'
);

CREATE TYPE user_status AS ENUM (
  'active',
  'blocked',
  'inverify'
);

CREATE TYPE platform AS ENUM (
  'admin',
  'web',
  'mobile'
);

CREATE TABLE users (
  id uuid PRIMARY KEY,
  user_type user_type NOT NULL,
  user_role user_role NOT NULL,
  full_name varchar(50) NOT NULL,
  username varchar(50) UNIQUE NOT NULL,
  password varchar(150) NOT NULL,
  email varchar(50) UNIQUE NOT NULL,
  gender gender NOT NULL DEFAULT 'male',
  avatar_id varchar(50) NOT NULL,
  status user_status NOT NULL DEFAULT 'inverify',
  created_at timestamp NOT NULL DEFAULT 'now()',
  updated_at timestamp NOT NULL DEFAULT 'now()'
);

CREATE TABLE session (
  id uuid PRIMARY KEY,
  user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  user_agent text NOT NULL,
  platform platform NOT NULL,
  ip_address varchar(64) NOT NULL,
  created_at timestamp NOT NULL DEFAULT 'now()',
  updated_at timestamp NOT NULL DEFAULT 'now()'
);

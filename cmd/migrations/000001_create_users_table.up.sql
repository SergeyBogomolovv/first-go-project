CREATE TABLE IF NOT EXISTS users (
  id int GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  username varchar(255) NOT NULL,
  email varchar(255) UNIQUE NOT NULL,
  password bytea NOT NULL,
  created_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
);
CREATE TABLE IF NOT EXISTS users (
  id int GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  content text NOT NULL,
  user_id int REFERENCES users(id) NOT NULL,
  title text NOT NULL,
  created_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
);
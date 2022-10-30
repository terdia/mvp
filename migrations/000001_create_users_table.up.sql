CREATE TABLE IF NOT EXISTS users (
     id bigserial PRIMARY KEY,
     username varchar UNIQUE NOT NULL,
     password_hash bytea NOT NULL,
     deposit numeric NOT NULL default 0,
     role varchar(15) NOT NULL,
     created_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS username_idx ON users (username);
CREATE INDEX IF NOT EXISTS username_idx ON users (role);

ALTER TABLE users ADD CONSTRAINT role_check CHECK (role in ('seller', 'buyer'));
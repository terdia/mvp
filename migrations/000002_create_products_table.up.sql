CREATE TABLE IF NOT EXISTS products (
     id bigserial PRIMARY KEY,
     name varchar NOT NULL,
     cost numeric NOT NULL,
     quantity numeric NOT NULL,
     seller_id bigint NOT NULL REFERENCES users ON DELETE CASCADE,
     created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
     UNIQUE (name, seller_id)
);

ALTER TABLE products ADD CONSTRAINT cost_check CHECK (cost % 5 = 0);

CREATE INDEX IF NOT EXISTS products_title_idx ON products USING GIN (to_tsvector('simple', name));

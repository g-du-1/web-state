CREATE TABLE IF NOT EXISTS pagestates (
    id serial PRIMARY KEY,
    url varchar(500) NOT NULL,
    scroll_pos integer,
    visible_text text,
    updated_at timestamp DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (url)
);
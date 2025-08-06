CREATE TABLE IF NOT EXISTS pagestates (
    id serial,
    url varchar(500),
    scroll_pos integer,
    created_at timestamp DEFAULT CURRENT_TIMESTAMP
);
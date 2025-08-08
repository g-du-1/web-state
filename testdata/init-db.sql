CREATE TABLE IF NOT EXISTS pagestates (
    id serial,
    url varchar(500),
    scroll_pos integer,
    visible_text text,
    created_at timestamp DEFAULT CURRENT_TIMESTAMP
);
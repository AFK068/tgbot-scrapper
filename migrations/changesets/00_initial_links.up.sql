CREATE TABLE tg_users (
    tg_id BIGINT PRIMARY KEY
);

CREATE TABLE links (
    id BIGSERIAL PRIMARY KEY,
    url TEXT NOT NULL UNIQUE
);

CREATE TABLE user_link (
    tg_user_id BIGINT REFERENCES tg_users(tg_id) ON DELETE CASCADE,
    link_id BIGINT REFERENCES links(id),
    last_update TIMESTAMP,
    filters TEXT[],
    tags TEXT[],
    PRIMARY KEY (tg_user_id, link_id)
);

CREATE INDEX user_link_user_id_last_update_idx ON user_link(tg_user_id, last_update);

CREATE INDEX user_link_tags_idx ON user_link(tags);

CREATE INDEX links_last_update_idx ON user_link(last_update);

CREATE INDEX links_url_idx ON links(url);
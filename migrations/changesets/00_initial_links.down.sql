DROP INDEX IF EXISTS user_link_user_id_last_update_idx;
DROP INDEX IF EXISTS user_link_tags_idx;
DROP INDEX IF EXISTS links_last_update_idx;
DROP INDEX IF EXISTS links_url_idx;

DROP TABLE IF EXISTS user_link;
DROP TABLE IF EXISTS links;
DROP TABLE IF EXISTS tg_users;
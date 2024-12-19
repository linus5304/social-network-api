drop extension if exists pg_trgm;

drop index if exists idx_comments_content;

drop index if exists idx_posts_title;

drop index if exists idx_posts_tags;

drop index if exists idx_users_username;

drop index if exists idx_posts_users_id;

drop index if exists idx_comments_post_id;
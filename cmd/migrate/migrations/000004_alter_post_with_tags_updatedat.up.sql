alter table posts
add column tags varchar(100) [];

alter table posts
add column updated_at timestamptz not null default now ();
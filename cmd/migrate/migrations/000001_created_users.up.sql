create extension if not exists citext;

create table
    if not exists users (
        id bigserial primary key,
        username varchar(255) unique not null,
        email citext unique not null,
        password bytea not null,
        created_at timestamptz not null default now (),
        updated_at timestamptz not null default now ()
    );
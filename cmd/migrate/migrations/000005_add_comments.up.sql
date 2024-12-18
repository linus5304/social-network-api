create table
    if not exists comments (
        id bigserial primary key,
        post_id bigserial not null,
        user_id bigserial not null,
        content text not null,
        created_at timestamptz not null default now (),
        updated_at timestamptz not null default now ()
    );
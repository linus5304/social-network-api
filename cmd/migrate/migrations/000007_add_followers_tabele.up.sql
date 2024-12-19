create table
    if not exists followers (
        user_id bigint not null,
        follower_id bigint not null,
        created_at timestamptz not null default now (),
        primary key (user_id, follower_id), -- composite key
        foreign key (user_id) references users (id) on delete cascade,
        foreign key (user_id) references users (id) on delete cascade
    );
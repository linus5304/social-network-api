create table
    if not exists roles (
        id bigserial primary key,
        name varchar(255) not null unique,
        level int not null default 0,
        description text
    );

insert into
    roles (name, description, level)
values
    ('user', 'A user can create posts and comments', 1);

insert into
    roles (name, description, level)
values
    (
        'moderator',
        'A moderator can update other users posts',
        2
    );

insert into
    roles (name, description, level)
values
    (
        'admin',
        'An admin can update and delete other users posts',
        3
    );
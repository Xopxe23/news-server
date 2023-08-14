CREATE TABLE users (
    id serial not null unique,
    name varchar(255) not null,
    email varchar(255) not null unique,
    password_hash varchar(255) not null
);

CREATE TABLE articles (
    id serial not null unique,
    author_id int references users (id) on delete cascade not null,
    title varchar(255) not null,
    content text not null,
    created_at timestamp default now() not null
);
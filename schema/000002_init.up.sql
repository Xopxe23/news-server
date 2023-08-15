CREATE TABLE refresh_tokens (
    id serial not null unique,
    user_id int references users (id) on delete cascade not null,
    token varchar(255) not null unique,
    expires_at timestamp not null
);
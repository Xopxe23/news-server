CREATE TABLE authors (
    id serial not null unique,
    name varchar(255) not null,
    surname varchar(255) not null
);
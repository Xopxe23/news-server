CREATE TABLE users (
    id serial not null unique,
    name varchar(255) not null,
    email varchar(255) not null unique,
    password_hash varchar(255) not null
);

CREATE TABLE authors (
    id serial not null unique,
    name varchar(255) not null,
    surname varchar(255) not null
);

CREATE TABLE articles (
    id serial not null unique,
    author_id int references authors (id) on delete cascade not null,
    title varchar(255) not null,
    content text not null,
    created_at timestamp default now() not null
);

CREATE TABLE refresh_tokens (
    id serial not null unique,
    user_id int references users (id) on delete cascade not null,
    token varchar(255) not null unique,
    expires_at timestamp not null
);

CREATE TABLE bookmarks (
    id serial not null unique,
    user_id int references users (id) on delete cascade not null,
    article_id int references articles (id) on delete cascade not null,
    UNIQUE (user_id, article_id)
);

INSERT INTO users (name, email, password_hash)
VALUES ('George', 'xopxe23@gmail.com', '73616c741c9cba7259cb51451bd282a646782900e56ea62b');

INSERT INTO authors (name, surname)
VALUES ('Alexander', 'Pushkin');

INSERT INTO articles (author_id, title, content) 
VALUES (1, 'Первые шесть внедорожников от органов власти Северной Осетии  готовы к отправке  на фронт', 'После  прохождения техобслуживания автомобили будут переданы бойцам на фронт.Благодаря органам власти, госкорпорациям и предприятиям Северной Осетии в рамках проекта Народный фронт.Всё для победы на фронт будут переданы автомобили повышенной проходимости.'), 
(1, 'Прокуратура Северной Осетии намерена через суд вернуть в общее пользование земли вдоль реки Фиагдон', 'Прокуратура Алагирского района по республике в судебном порядке требует вернуть государству земли, расположенные в береговой полосе реки Фиагдон, незаконно предоставленные в собственность, сообщили «КрыльямTV» в надзорном ведомстве.'), 
(1, 'Героини СВО.', 'Позывной Багира - наша боевая сестра, начальник медицинской службы одного из полков, параллельно решает проблемы по обеспечению дивизии медицинским имуществом и взаимодействию с гуманитарными организациями.');


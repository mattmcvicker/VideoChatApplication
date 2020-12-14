create table if not exists users (
    id int not null auto_increment primary key,
    email varchar(320) not null,
    UNIQUE (email),
    passhash char(60) not null,
    username varchar(255) not null,
    UNIQUE (username),
    firstname varchar(64) not null,
    lastname varchar(128) not null,
    photourl varchar(128) not null
);
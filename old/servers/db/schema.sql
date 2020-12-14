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

create table if not exists logs (
    id int not null auto_increment primary key,
    userid int not null,
    signin_time datetime not null,
    client_IP varchar(50) not null
);



create table if not exists questions {
    questionid int not null auto_increment primary key,
    topicid int not null,
    questionbody varchar(255)
};

create table if not exists answers {
    questionid int not null,
    userid int not null,
    answer boolean not null
};
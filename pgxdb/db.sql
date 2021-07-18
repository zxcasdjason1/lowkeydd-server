CREATE TABLE auth (
    id serial primary key,
    userid varchar(255) not null check(length(userid)>3),
    passwd varchar(255) not null check(length(passwd)>3), 
    is_banned boolean default FALSE,
    create_date timestamp default 'now'
);
GRANT ALL PRIVILEGES ON TABLE auth TO nilson;


CREATE TABLE users (
    id serial primary key,
    userid varchar(255) not null check(length(userid)>3),
    data text,
    is_del boolean default FALSE,
    create_date timestamp default 'now'
);
GRANT ALL PRIVILEGES ON TABLE users TO nilson;


CREATE TABLE visit (
    id serial primary key,
    userid varchar(255) not null check(length(userid)>3),
    data text,
    is_del boolean default FALSE,
    create_date timestamp default 'now'
);
GRANT ALL PRIVILEGES ON TABLE visit TO nilson;


GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO nilson;
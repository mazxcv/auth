CREATE TABLE IF NOT EXISTS users 
(
    id integer primary key,
    email text not null unique,
    pass_hash BLOB not null
);

create index if not EXISTS idx_email on users(email);

create TABLE if not exists app 
(
    id integer primary key,
    name text not null unique,
    secret text not null unique
);
create table users (
    id char(36) primary key,
    name varchar(100) not null,
    email varchar(30) unique not null,
    password varchar(100) not null,
    photo_link varchar(255),
    created_at timestamptz default current_timestamp
    updated_at timestamptz default current_timestamp
)
-- +goose Up
-- +goose StatementBegin

create table if not exists user_role (
    id serial primary key,
    name text
);

insert into user_role (name) values ('unspecified');

insert into user_role (name) values ('user');

insert into user_role (name) values ('admin');

create table if not exists users (
    id serial primary key,
    name text not null,
    email text not null,
    password text not null,
    password_confirmation text not null,
    role integer not null references user_role (id)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists users
drop table if exists user_role
-- +goose StatementEnd
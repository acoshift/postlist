create database if not exists postlist;

create table if not exists posts (
  id serial,
  name string not null,
  content string not null,
  created_at timestamp not null default now(),
  primary key (id),
  index (created_at desc)
);


create table chat
(
  message_id serial       not null
    constraint chat_pkey
      primary key,
  from_user  integer      not null,
  to_user    integer      not null,
  text       varchar(255) not null,
  dt         timestamp default now()
);

alter table chat
  owner to postgres;

create table users
(
  user_id serial not null
    constraint users_pkey
      primary key,
  name    varchar(30)
);

alter table users
  owner to postgres;

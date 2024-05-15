create table usr (
  id int primary key generated always as identity,
  login varchar(255) unique,
  hash varchar(255)
);

-- create type scrt_type as enum('LOGIN_PASSWORD', 'T');
create table scrt (
  id int primary key generated always as identity,
  usr_id int references usr(id),
  type int not null,
  name varchar(255) not null,
  value bytea not null,
  nonce bytea not null,
  encryption_sk bytea not null,
  created_at timestamptz,
  updated_at timestamptz,
  constraint UQ_scrt_usr_id_name unique(usr_id, name)
);

-- drop table scrt_attr;
-- drop table attr;
create table attr (
  id int primary key generated always as identity,
  name varchar(255) unique
);

create table scrt_attr (
  scrt_id int references scrt(id),
  attr_id int references attr(id),
  value varchar(1024),
  constraint PK_scrt_attr primary key (scrt_id, attr_id)
);
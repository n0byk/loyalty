create extension if not exists "uuid-ossp"; 

create table public."user_catalog" ( 
user_id uuid not null default uuid_generate_v4(),
user_login varchar not null UNIQUE, 
user_password varchar not null,
user_salt uuid not null, 
add_time timestamp NOT NULL DEFAULT (now() at time zone 'UTC'),
update_time timestamp NOT NULL DEFAULT (now() at time zone 'UTC'),
delete_time timestamp,
constraint user_pkey primary key (user_id)
);

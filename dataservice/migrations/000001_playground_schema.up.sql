create extension if not exists "uuid-ossp"; 

CREATE TABLE IF NOT EXISTS public."user_catalog" ( 
user_id uuid not null default uuid_generate_v4(),
user_login varchar not null UNIQUE, 
user_password varchar not null,
user_salt uuid not null, 
add_time timestamp NOT NULL DEFAULT (now() at time zone 'UTC'),
update_time timestamp NOT NULL DEFAULT (now() at time zone 'UTC'),
delete_time timestamp,
constraint user_pkey primary key (user_id)
);

CREATE TYPE order_state_type  as enum ('NEW', 'PROCESSING', 'INVALID', 'PROCESSED');

CREATE TABLE IF NOT EXISTS public."order_catalog" ( 
order_id uuid not null default uuid_generate_v4(),
user_id uuid REFERENCES user_catalog (user_id),
order_number varchar unique,
order_state order_state_type default 'NEW',
add_time timestamp NOT NULL DEFAULT (now() at time zone 'UTC'),
update_time timestamp NOT NULL DEFAULT (now() at time zone 'UTC'),
delete_time timestamp,
constraint order_pkey primary key (order_id),
unique (user_id,order_number)
);

CREATE TABLE IF NOT EXISTS public."balance_catalog" ( 
balance_id uuid not null default uuid_generate_v4(), 
order_number varchar not null, 
accrue decimal default 0, 
withdraw decimal default 0, 
add_time timestamp NOT NULL DEFAULT (now() at time zone 'UTC'),
update_time timestamp NOT NULL DEFAULT (now() at time zone 'UTC'),
delete_time timestamp,
constraint balance_pkey primary key (balance_id)
);
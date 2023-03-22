CREATE TABLE orders (
   order_uid varchar not null primary key,
   track_number varchar not null,
   entry varchar not null,
   delivery jsonb not null,
   payment jsonb not null,
   items jsonb not null,
   locale varchar(5) not null,
   internal_signature varchar,
   customer_id varchar unique not null,
   delivery_service varchar not null,
   shardkey varchar not null,
   sm_id integer not null,
   date_created timestamptz not null,
   oof_shard varchar not null
);
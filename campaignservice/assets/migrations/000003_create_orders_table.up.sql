create type order_status as enum (
    'PENDING',
    'PROCESSING',
    'SUCCESS',
    'FAILED'
);

create table orders (
    id varchar(20) primary key not null,
    status order_status not null default 'PENDING',
    total_amount decimal(10, 2) not null,
    customer_id varchar(20) not null,
    customer_name varchar(100) not null,
    created_at timestamp with time zone not null default now(),
    updated_at timestamp with time zone not null default now()
);

select create_updated_at_trigger('orders');

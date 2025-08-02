create type payment_status as enum (
    'COMPLETED',
    'FAILED'
);

create table payments (
    id varchar(20) primary key not null,
    order_id varchar(20) not null,
    status payment_status not null default 'COMPLETED',
    total_amount decimal(10, 2) not null,
    customer_id varchar(20) not null,
    created_at timestamp with time zone not null default now(),
    updated_at timestamp with time zone not null default now()
);

select create_updated_at_trigger('payments');

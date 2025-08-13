
create table winners (
    id varchar(20) primary key not null,
    campaign_id bigserial not null references campaigns(id),
    customer_id varchar(20) not null,
    created_at timestamp with time zone not null default now(),
    updated_at timestamp with time zone not null default now()
);

select create_updated_at_trigger('winners');



create table campaigns (
    id integer primary key not null,
    name varchar(100) not null,
    description text not null,
    start_time timestamp with time zone not null,
    end_time timestamp with time zone not null,
    policy jsonb not null default '{}',
    created_at timestamp with time zone not null default now(),
    updated_at timestamp with time zone not null default now()
);

select create_updated_at_trigger('campaigns');

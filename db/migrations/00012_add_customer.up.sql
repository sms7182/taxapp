
create table customers
(
    id                        bigserial
        constraint customer_pkey
            primary key,
    created_at timestamp    not null,
    finance_id varchar(256) not null,
    token varchar(256) ,
    public_key text,
    private_key text,
    expire_time timestamp

);

alter table customers
    owner to postgres;



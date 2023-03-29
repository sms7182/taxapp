
create table tax_process_history
(
    id                        bigserial
        constraint tax_process_history_pkey
            primary key,
    tax_process_id       bigint       not null,
    status                    varchar(128) not null,
    tax_type                  varchar(128) not null,
    tax_raw_id                bigserial,
    tax_unique_id             varchar(128) not null,
    created_at                timestamp    not null,
    tax_org_reference_id      varchar(256),
    tax_id                    varchar(256),
    standard_invoice          jsonb,
    internal_trn              varchar(256),
    inquiry_uuid              varchar(256),
    confirmation_reference_id varchar(256),
    company_name              varchar(256)
);

alter table tax_process
    owner to pgsql;



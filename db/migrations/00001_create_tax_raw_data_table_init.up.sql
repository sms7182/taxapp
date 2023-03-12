create table if not exists tax_raw_data
(
    id  bigserial not null
    primary key,
    created_at  timestamp not null,
    tax_data  jsonb
    
	
	
);



create table if not exists tax_office_request_response_log
(
    id              bigserial    not null
    constraint tax_office_request_response_log_pkey
    primary key,
    tax_raw_id      bigint  null ,
    tax_process_id  bigint null,
    api_name     text         not null,
    request_unique_id text,
    logged_at       timestamp    not null,
    url             text         not null,
    status_code     int          not null,
    request         text         not null,
    response        text,
    error_message   text
   
);

create table if not exists tax_process(
    id bigserial not null
    constraint tax_process_pkey
    primary key,
    status varchar(128) not null,
    tax_type varchar(128) not null,
    tax_raw_id bigserial,
    tax_unique_id varchar(128) not null,
    created_at timestamp not null,
    updated_at timestamp not null
);
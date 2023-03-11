create table if not exists tax_raw_data
(
    id  bigserial not null
    primary key,
    created_at  timestamp not null,
    updated_at timestamp not null,
    trace_id text,
    tax_data  jsonb,
	tax_type     text,
	process_status text
);

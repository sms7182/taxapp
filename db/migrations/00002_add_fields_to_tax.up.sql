alter table tax_raw_data add column if not exists unique_id varchar(256) not null default '';
alter table tax_raw_data add column if not exists tax_type varchar(256) not null default '';

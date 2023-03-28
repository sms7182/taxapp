alter table tax_process add column if not exists internal_trn varchar(256);
alter table tax_process add column if not exists inquiry_uuid varchar(256);


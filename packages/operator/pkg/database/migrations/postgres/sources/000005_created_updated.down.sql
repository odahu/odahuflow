BEGIN;
alter table odahu_operator_training drop column created;
alter table odahu_operator_training drop column updated;
alter table odahu_operator_packaging drop column created;
alter table odahu_operator_packaging drop column updated;
alter table odahu_operator_deployment drop column created;
alter table odahu_operator_deployment drop column updated;
COMMIT;
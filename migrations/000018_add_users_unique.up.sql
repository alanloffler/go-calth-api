ALTER TABLE users
ADD CONSTRAINT uq_users_business_id_id UNIQUE (business_id, id);

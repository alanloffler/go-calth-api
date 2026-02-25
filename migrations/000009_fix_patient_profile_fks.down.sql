ALTER TABLE patient_profile
  DROP CONSTRAINT fk_patient_profile_business;

ALTER TABLE patient_profile
  DROP CONSTRAINT fk_patient_profile_user;

ALTER TABLE patient_profile
  ADD CONSTRAINT fk_patient_profile_user FOREIGN KEY (user_id) REFERENCES users (id);

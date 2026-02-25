ALTER TABLE patient_profile
  ADD CONSTRAINT fk_patient_profile_business FOREIGN KEY (business_id) REFERENCES businesses (id) ON DELETE CASCADE;

ALTER TABLE patient_profile
  DROP CONSTRAINT fk_patient_profile_user;

ALTER TABLE patient_profile
  ADD CONSTRAINT fk_patient_profile_user FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE;

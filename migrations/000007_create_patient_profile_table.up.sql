CREATE TABLE patient_profile (
  id                      UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
  business_id             UUID        NOT NULL,
  user_id                 UUID        NOT NULL,
  gender                  VARCHAR(20) NOT NULL,
  birth_day               DATE        NOT NULL,
  blood_type              VARCHAR(20) NOT NULL,
  weight                  NUMERIC     NOT NULL,
  height                  NUMERIC     NOT NULL,
  emergency_contact_name  VARCHAR(50) NOT NULL,
  emergency_contact_phone VARCHAR(10) NOT NULL,
  created_at              TIMESTAMPTZ NOT NULL    DEFAULT now(),
  updated_at              TIMESTAMPTZ NOT NULL    DEFAULT now(),
  deleted_at              TIMESTAMPTZ,

  CONSTRAINT fk_patient_profile_user FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE INDEX idx_patient_profile_business ON patient_profile (business_id);
CREATE INDEX idx_patient_profile_business_user ON patient_profile (business_id, user_id);
CREATE UNIQUE INDEX uq_patient_profile_business_user ON patient_profile (business_id, user_id);

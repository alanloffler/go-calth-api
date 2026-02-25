CREATE TABLE professional_profile (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  business_id UUID NOT NULL,
  user_id UUID NOT NULL,
  license_id VARCHAR NOT NULL,
  professional_prefix VARCHAR NOT NULL,
  specialty VARCHAR NOT NULL,
  working_days VARCHAR NOT NULL,
  start_hour VARCHAR NOT NULL DEFAULT '07:00',
  end_hour VARCHAR NOT NULL DEFAULT '20:00',
  slot_duration VARCHAR NOT NULL DEFAULT '60',
  daily_exception_start VARCHAR,
  daily_exception_end VARCHAR,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMPTZ,

  CONSTRAINT uq_prof_profile_business_license UNIQUE (business_id, license_id),
  CONSTRAINT fk_professional_profile_business FOREIGN KEY (business_id) REFERENCES businesses (id) ON DELETE CASCADE,
  CONSTRAINT fk_professional_profile_user FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE INDEX idx_prof_profile_business ON professional_profile (business_id);

CREATE INDEX idx_prof_profile_business_user ON professional_profile (business_id, user_id);

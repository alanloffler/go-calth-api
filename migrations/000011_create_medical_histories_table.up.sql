CREATE TABLE medical_histories (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  business_id UUID NOT NULL,
  user_id UUID NOT NULL,
  professional_id UUID NOT NULL,
  event_id UUID NULL,
  date TIMESTAMPTZ NOT NULL,
  reason VARCHAR NOT NULL,
  recipe BOOLEAN NOT NULL DEFAULT FALSE,
  comments VARCHAR NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  deleted_at TIMESTAMPTZ NULL,
  CONSTRAINT fk_mh_user FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE RESTRICT,
  CONSTRAINT fk_mh_professional FOREIGN KEY (professional_id) REFERENCES users (id) ON DELETE SET NULL,
  CONSTRAINT fk_mh_event FOREIGN KEY (event_id) REFERENCES events (id) ON DELETE SET NULL
);

CREATE INDEX idx_mh_business_user ON medical_histories (business_id, user_id);

CREATE INDEX idx_mh_business_event ON medical_histories (business_id, event_id);

CREATE INDEX idx_mh_business_user_created ON medical_histories (business_id, user_id, created_at);

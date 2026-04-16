CREATE TABLE blocked_days (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  date TIMESTAMPTZ NOT NULL,
  reason VARCHAR(50) NOT NULL,
  business_id UUID NOT NULL,
  professional_id UUID NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  CONSTRAINT fk_blocked_days_professional FOREIGN KEY (business_id, professional_id) REFERENCES users (business_id, id) ON DELETE CASCADE,
  CONSTRAINT uq_blocked_days_unique UNIQUE (business_id, professional_id, date)
);

CREATE INDEX idx_blocked_days_business_professional ON blocked_days (business_id, professional_id, date)

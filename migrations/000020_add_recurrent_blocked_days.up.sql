ALTER TABLE blocked_days
ADD COLUMN recurrent BOOLEAN DEFAULT false;

CREATE INDEX idx_blocked_business_professional_recurrent ON blocked_days (business_id, professional_id)
WHERE
  recurrent IS TRUE;

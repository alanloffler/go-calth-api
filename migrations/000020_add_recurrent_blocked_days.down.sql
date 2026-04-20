ALTER TABLE blocked_days
DROP COLUMN recurrent;

DROP INDEX IF EXISTS idx_blocked_business_professional_recurrent;

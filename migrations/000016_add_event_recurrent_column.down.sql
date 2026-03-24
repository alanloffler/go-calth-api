ALTER TABLE events
DROP COLUMN recurrent_id;

DROP INDEX IF EXISTS idx_events_recurrent_id;

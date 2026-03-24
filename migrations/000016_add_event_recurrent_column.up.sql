ALTER TABLE events
ADD COLUMN recurrent_id UUID;

CREATE INDEX idx_events_recurrent_id ON events (business_id, recurrent_id)
WHERE
  recurrent_id IS NOT NULL;

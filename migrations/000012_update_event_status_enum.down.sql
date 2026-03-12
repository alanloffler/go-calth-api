ALTER TABLE events
ALTER COLUMN status
DROP DEFAULT;

ALTER TYPE event_status
RENAME TO event_status_new;

CREATE TYPE event_status AS ENUM(
  'absent',
  'attended',
  'cancelled',
  'in_progress',
  'pending'
);

ALTER TABLE events
ALTER COLUMN status TYPE event_status USING (
  CASE status::text
    WHEN 'present' THEN 'attended'::event_status
    ELSE status::text::event_status
  END
);

ALTER TABLE events
ALTER COLUMN status
SET DEFAULT 'pending'::event_status;

DROP TYPE event_status_new;

ALTER TABLE events
ALTER COLUMN status
DROP DEFAULT;

ALTER TYPE event_status
RENAME TO event_status_old;

CREATE TYPE event_status AS ENUM(
  'absent',
  'cancelled',
  'in_progress',
  'pending',
  'present'
);

ALTER TABLE events
ALTER COLUMN status TYPE event_status USING (
  CASE status::text
    WHEN 'attended' THEN 'present'::event_status
    ELSE status::text::event_status
  END
);

ALTER TABLE events
ALTER COLUMN status
SET DEFAULT 'pending'::event_status;

DROP TYPE event_status_old;

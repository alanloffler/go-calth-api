CREATE EXTENSION IF NOT EXISTS btree_gist;

ALTER TABLE events
ADD CONSTRAINT events_no_overlap EXCLUDE USING gist (
  business_id
  WITH
    =,
    professional_id
  WITH
    =,
    tstzrange (start_date, end_date, '[)')
  WITH
    &&
)
WHERE
  (deleted_at IS NULL);
